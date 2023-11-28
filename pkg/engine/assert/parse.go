package assert

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/match"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Parse(ctx context.Context, assertion interface{}) Assertion {
	switch reflectutils.GetKind(assertion) {
	case reflect.Slice:
		node := sliceNode{}
		valueOf := reflect.ValueOf(assertion)
		for i := 0; i < valueOf.Len(); i++ {
			node = append(node, Parse(ctx, valueOf.Index(i).Interface()))
		}
		return node
	case reflect.Map:
		node := mapNode{}
		iter := reflect.ValueOf(assertion).MapRange()
		for iter.Next() {
			node[iter.Key().Interface()] = Parse(ctx, iter.Value().Interface())
		}
		return node
	default:
		return &scalarNode{rhs: assertion}
	}
}

// mapNode is the assertion type represented by a map.
// it is responsible for projecting the analysed resource and passing the result to the descendant
type mapNode map[interface{}]Assertion

func (n mapNode) assert(ctx context.Context, path *field.Path, value interface{}, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	for k, v := range n {
		projection, err := project(ctx, k, value, bindings, opts...)
		if err != nil {
			return nil, field.InternalError(path.Child(fmt.Sprint(k)), err)
		} else {
			if projection.binding != "" {
				bindings = bindings.Register("$"+projection.binding, jpbinding.NewBinding(projection.result))
			}
			if projection.foreach {
				projectedKind := reflectutils.GetKind(projection.result)
				if projectedKind == reflect.Slice {
					valueOf := reflect.ValueOf(projection.result)
					for i := 0; i < valueOf.Len(); i++ {
						bindings := bindings
						if projection.foreachName != "" {
							bindings = bindings.Register("$"+projection.foreachName, jpbinding.NewBinding(i))
						}
						if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)).Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
							return nil, err
						} else {
							errs = append(errs, _errs...)
						}
					}
				} else if projectedKind == reflect.Map {
					iter := reflect.ValueOf(projection.result).MapRange()
					for iter.Next() {
						key := iter.Key().Interface()
						bindings := bindings
						if projection.foreachName != "" {
							bindings = bindings.Register("$"+projection.foreachName, jpbinding.NewBinding(key))
						}
						if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)).Key(fmt.Sprint(key)), iter.Value().Interface(), bindings, opts...); err != nil {
							return nil, err
						} else {
							errs = append(errs, _errs...)
						}
					}
				} else {
					return nil, field.TypeInvalid(path.Child(fmt.Sprint(k)), projection.result, "expected a slice or a map")
				}
			} else {
				if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)), projection.result, bindings, opts...); err != nil {
					return nil, err
				} else {
					errs = append(errs, _errs...)
				}
			}
		}
	}
	return errs, nil
}

// sliceNode is the assertion type represented by a slice.
// it first compares the length of the analysed resource with the length of the descendants.
// if lengths match all descendants are evaluated with their corresponding items.
type sliceNode []Assertion

func (n sliceNode) assert(ctx context.Context, path *field.Path, value interface{}, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	if value == nil {
		errs = append(errs, field.Invalid(path, value, "value is null"))
	} else if reflectutils.GetKind(value) != reflect.Slice {
		return nil, field.TypeInvalid(path, value, "expected a slice")
	} else {
		valueOf := reflect.ValueOf(value)
		if valueOf.Len() != len(n) {
			errs = append(errs, field.Invalid(path, value, "lengths of slices don't match"))
		} else {
			for i := range n {
				if _errs, err := n[i].assert(ctx, path.Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
					return nil, err
				} else {
					errs = append(errs, _errs...)
				}
			}
		}
	}
	return errs, nil
}

// scalarNode is a terminal type of assertion.
// it receives a value and compares it with an expected value.
// the expected value can be the result of an expression.
type scalarNode struct {
	rhs interface{}
}

func (n *scalarNode) assert(ctx context.Context, path *field.Path, value interface{}, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	rhs := n.rhs
	expression := parseExpression(ctx, rhs)
	// we only project if the expression uses the engine syntax
	// this is to avoid the case where the value is a map and the RHS is a string
	if expression != nil && expression.engine != "" {
		if expression.foreachName != "" {
			return nil, field.Invalid(path, rhs, "foreach is not supported on the RHS")
		}
		if expression.binding != "" {
			return nil, field.Invalid(path, rhs, "binding is not supported on the RHS")
		}
		projected, err := template.Execute(ctx, expression.statement, value, bindings, opts...)
		if err != nil {
			return nil, field.InternalError(path, err)
		}
		rhs = projected
	}
	var errs field.ErrorList
	if match, err := match.Match(ctx, rhs, value); err != nil {
		return nil, field.InternalError(path, err)
	} else if !match {
		errs = append(errs, field.Invalid(path, value, expectValueMessage(rhs)))
	}
	return errs, nil
}

func expectValueMessage(value interface{}) string {
	switch t := value.(type) {
	case int64, int32, float64, float32, bool:
		// use simple printer for simple types
		return fmt.Sprintf("Expected value: %v", value)
	case string:
		return fmt.Sprintf("Expected value: %q", t)
	case fmt.Stringer:
		// anything that defines String() is better than raw struct
		return fmt.Sprintf("Expected value: %s", t.String())
	default:
		// fallback to raw struct
		// TODO: internal types have panic guards against json.Marshalling to prevent
		// accidental use of internal types in external serialized form.  For now, use
		// %#v, although it would be better to show a more expressive output in the future
		return fmt.Sprintf("Expected value: %#v", value)
	}
}
