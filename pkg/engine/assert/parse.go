package assert

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/kyverno/kyverno-json/pkg/core/projection"
	"github.com/kyverno/kyverno-json/pkg/engine/match"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Parse(ctx context.Context, assertion any) (Assertion, error) {
	switch reflectutils.GetKind(assertion) {
	case reflect.Slice:
		return parseSlice(ctx, assertion)
	case reflect.Map:
		return parseMap(ctx, assertion)
	default:
		return parseScalar(ctx, assertion)
	}
}

// node implements the Assertion interface using a delegate func
type node func(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error)

func (n node) assert(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	return n(ctx, path, value, bindings, opts...)
}

// parseSlice is the assertion represented by a slice.
// it first compares the length of the analysed resource with the length of the descendants.
// if lengths match all descendants are evaluated with their corresponding items.
func parseSlice(ctx context.Context, assertion any) (node, error) {
	var assertions []Assertion
	valueOf := reflect.ValueOf(assertion)
	for i := 0; i < valueOf.Len(); i++ {
		sub, err := Parse(ctx, valueOf.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		assertions = append(assertions, sub)
	}
	return func(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
		var errs field.ErrorList
		if value == nil {
			errs = append(errs, field.Invalid(path, value, "value is null"))
		} else if reflectutils.GetKind(value) != reflect.Slice {
			return nil, field.TypeInvalid(path, value, "expected a slice")
		} else {
			valueOf := reflect.ValueOf(value)
			if valueOf.Len() != len(assertions) {
				errs = append(errs, field.Invalid(path, value, "lengths of slices don't match"))
			} else {
				for i := range assertions {
					if _errs, err := assertions[i].assert(ctx, path.Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
						return nil, err
					} else {
						errs = append(errs, _errs...)
					}
				}
			}
		}
		return errs, nil
	}, nil
}

// parseMap is the assertion represented by a map.
// it is responsible for projecting the analysed resource and passing the result to the descendant
func parseMap(ctx context.Context, assertion any) (node, error) {
	assertions := map[any]struct {
		projection.Projection
		Assertion
	}{}
	iter := reflect.ValueOf(assertion).MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		value := iter.Value().Interface()
		assertion, err := Parse(ctx, value)
		if err != nil {
			return nil, err
		}
		entry := assertions[key]
		entry.Assertion = assertion
		entry.Projection = projection.Parse(key)
		assertions[key] = entry
	}
	return func(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
		var errs field.ErrorList
		// if we assert against an empty object, value is expected to be not nil
		if len(assertions) == 0 {
			if value == nil {
				errs = append(errs, field.Invalid(path, value, "invalid value, must not be null"))
			}
			return errs, nil
		}
		for k, v := range assertions {
			projected, found, err := v.Projection.Handler(ctx, value, bindings, opts...)
			if err != nil {
				return nil, field.InternalError(path.Child(fmt.Sprint(k)), err)
			} else if !found {
				errs = append(errs, field.Required(path.Child(fmt.Sprint(k)), "field not found in the input object"))
			} else {
				if v.Projection.Binding != "" {
					bindings = bindings.Register("$"+v.Projection.Binding, binding.NewBinding(projected))
				}
				if v.Projection.Foreach {
					projectedKind := reflectutils.GetKind(projected)
					if projectedKind == reflect.Slice {
						valueOf := reflect.ValueOf(projected)
						for i := 0; i < valueOf.Len(); i++ {
							bindings := bindings
							if v.Projection.ForeachName != "" {
								bindings = bindings.Register("$"+v.Projection.ForeachName, binding.NewBinding(i))
							}
							if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)).Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
								return nil, err
							} else {
								errs = append(errs, _errs...)
							}
						}
					} else if projectedKind == reflect.Map {
						iter := reflect.ValueOf(projected).MapRange()
						for iter.Next() {
							key := iter.Key().Interface()
							bindings := bindings
							if v.Projection.ForeachName != "" {
								bindings = bindings.Register("$"+v.Projection.ForeachName, binding.NewBinding(key))
							}
							if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)).Key(fmt.Sprint(key)), iter.Value().Interface(), bindings, opts...); err != nil {
								return nil, err
							} else {
								errs = append(errs, _errs...)
							}
						}
					} else {
						return nil, field.TypeInvalid(path.Child(fmt.Sprint(k)), projected, "expected a slice or a map")
					}
				} else {
					if _errs, err := v.assert(ctx, path.Child(fmt.Sprint(k)), projected, bindings, opts...); err != nil {
						return nil, err
					} else {
						errs = append(errs, _errs...)
					}
				}
			}
		}
		return errs, nil
	}, nil
}

// parseScalar is the assertion represented by a leaf.
// it receives a value and compares it with an expected value.
// the expected value can be the result of an expression.
func parseScalar(_ context.Context, assertion any) (node, error) {
	var project func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error)
	switch typed := assertion.(type) {
	case string:
		expr := expression.Parse(typed)
		if expr.Foreach {
			return nil, errors.New("foreach is not supported on the RHS")
		}
		if expr.Binding != "" {
			return nil, errors.New("binding is not supported on the RHS")
		}
		switch expr.Engine {
		case expression.EngineJP:
			parse := sync.OnceValues(func() (parsing.ASTNode, error) {
				parser := parsing.NewParser()
				return parser.Parse(expr.Statement)
			})
			project = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				ast, err := parse()
				if err != nil {
					return nil, err
				}
				return template.ExecuteAST(ctx, ast, value, bindings, opts...)
			}
		case expression.EngineCEL:
			project = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				return template.ExecuteCEL(ctx, expr.Statement, value, bindings)
			}
		default:
			assertion = expr.Statement
		}
	}
	return func(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
		expected := assertion
		if project != nil {
			projected, err := project(ctx, value, bindings, opts...)
			if err != nil {
				return nil, field.InternalError(path, err)
			}
			expected = projected
		}
		var errs field.ErrorList
		if match, err := match.Match(ctx, expected, value); err != nil {
			return nil, field.InternalError(path, err)
		} else if !match {
			errs = append(errs, field.Invalid(path, value, expectValueMessage(expected)))
		}
		return errs, nil
	}, nil
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
