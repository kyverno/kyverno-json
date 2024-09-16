package assert

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/engine/match"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Parse(ctx context.Context, assertion any) (Assertion, error) {
	switch reflectutils.GetKind(assertion) {
	case reflect.Slice:
		node := sliceNode{}
		valueOf := reflect.ValueOf(assertion)
		for i := 0; i < valueOf.Len(); i++ {
			sub, err := Parse(ctx, valueOf.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			node = append(node, sub)
		}
		return node, nil
	case reflect.Map:
		node := mapNode{}
		iter := reflect.ValueOf(assertion).MapRange()
		for iter.Next() {
			sub, err := Parse(ctx, iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			node[iter.Key().Interface()] = sub
		}
		return node, nil
	default:
		return newScalarNode(ctx, nil, assertion)
	}
}

// mapNode is the assertion type represented by a map.
// it is responsible for projecting the analysed resource and passing the result to the descendant
type mapNode map[any]Assertion

func (n mapNode) assert(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	// if we assert against an empty object, value is expected to be not nil
	if len(n) == 0 {
		if value == nil {
			errs = append(errs, field.Invalid(path, value, "invalid value, must not be null"))
		}
		return errs, nil
	}
	for k, v := range n {
		projection, err := project(ctx, k, value, bindings, opts...)
		if err != nil {
			return nil, field.InternalError(path.Child(fmt.Sprint(k)), err)
		} else if projection == nil {
			errs = append(errs, field.Required(path.Child(fmt.Sprint(k)), "field not found in the input object"))
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

func (n sliceNode) assert(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
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
	project func(value any, bindings binding.Bindings, opts ...template.Option) (any, error)
}

func newScalarNode(ctx context.Context, path *field.Path, rhs any) (Assertion, error) {
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
		parser := parsing.NewParser()
		compiled, err := parser.Parse(expression.statement)
		if err != nil {
			return nil, field.InternalError(path, err)
		}
		return &scalarNode{
			project: func(value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				o := template.BuildOptions(opts...)
				vm := interpreter.NewInterpreter(nil, bindings)
				return vm.Execute(compiled, value, interpreter.WithFunctionCaller(o.FunctionCaller))
			},
		}, nil
	} else {
		return &scalarNode{
			project: func(value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				return rhs, nil
			},
		}, nil
	}
}

func (n *scalarNode) assert(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	if rhs, err := n.project(value, bindings, opts...); err != nil {
		return nil, field.InternalError(path, err)
	} else if match, err := match.Match(ctx, rhs, value); err != nil {
		return nil, field.InternalError(path, err)
	} else if !match {
		errs = append(errs, field.Invalid(path, value, expectValueMessage(rhs)))
	}
	return errs, nil
}

func expectValueMessage(value any) string {
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
