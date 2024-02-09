package mutate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Parse(ctx context.Context, mutation any) Mutation {
	switch reflectutils.GetKind(mutation) {
	case reflect.Slice:
		node := sliceNode{}
		valueOf := reflect.ValueOf(mutation)
		for i := 0; i < valueOf.Len(); i++ {
			node = append(node, Parse(ctx, valueOf.Index(i).Interface()))
		}
		return node
	case reflect.Map:
		node := mapNode{}
		iter := reflect.ValueOf(mutation).MapRange()
		for iter.Next() {
			node[iter.Key().Interface()] = Parse(ctx, iter.Value().Interface())
		}
		return node
	default:
		return &scalarNode{rhs: mutation}
	}
}

// mapNode is the mutation type represented by a map.
// it is responsible for projecting the analysed resource and passing the result to the descendant
type mapNode map[any]Mutation

func (n mapNode) mutate(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
	out := map[any]any{}
	for k, v := range n {
		// TODO: very simple implementation
		mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(k))
		if !mapValue.IsValid() {
			continue
		}
		value := mapValue.Interface()
		// TODO: does it make sense to take valueOf.Index(i).Interface() here ?
		if inner, err := v.mutate(ctx, path.Child(fmt.Sprint(k)), value, bindings, opts...); err != nil {
			return nil, err
		} else {
			out[k] = inner
		}
		// projection, err := project(ctx, k, value, bindings, opts...)
		// if err != nil {
		// 	return nil, field.InternalError(path.Child(fmt.Sprint(k)), err)
		// } else {
		// 	if projection.binding != "" {
		// 		bindings = bindings.Register("$"+projection.binding, jpbinding.NewBinding(projection.result))
		// 	}
		// 	if projection.foreach {
		// 		projectedKind := reflectutils.GetKind(projection.result)
		// 		if projectedKind == reflect.Slice {
		// 			valueOf := reflect.ValueOf(projection.result)
		// 			for i := 0; i < valueOf.Len(); i++ {
		// 				bindings := bindings
		// 				if projection.foreachName != "" {
		// 					bindings = bindings.Register("$"+projection.foreachName, jpbinding.NewBinding(i))
		// 				}
		// 				if _errs, err := v.mutate(ctx, path.Child(fmt.Sprint(k)).Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
		// 					return nil, err
		// 				} else {
		// 					errs = append(errs, _errs...)
		// 				}
		// 			}
		// 		} else if projectedKind == reflect.Map {
		// 			iter := reflect.ValueOf(projection.result).MapRange()
		// 			for iter.Next() {
		// 				key := iter.Key().Interface()
		// 				bindings := bindings
		// 				if projection.foreachName != "" {
		// 					bindings = bindings.Register("$"+projection.foreachName, jpbinding.NewBinding(key))
		// 				}
		// 				if _errs, err := v.mutate(ctx, path.Child(fmt.Sprint(k)).Key(fmt.Sprint(key)), iter.Value().Interface(), bindings, opts...); err != nil {
		// 					return nil, err
		// 				} else {
		// 					errs = append(errs, _errs...)
		// 				}
		// 			}
		// 		} else {
		// 			return nil, field.TypeInvalid(path.Child(fmt.Sprint(k)), projection.result, "expected a slice or a map")
		// 		}
		// 	} else {
		// 		if _errs, err := v.mutate(ctx, path.Child(fmt.Sprint(k)), projection.result, bindings, opts...); err != nil {
		// 			return nil, err
		// 		} else {
		// 			errs = append(errs, _errs...)
		// 		}
		// 	}
		// }
	}
	return out, nil
}

// sliceNode is the mutation type represented by a slice.
// it first compares the length of the analysed resource with the length of the descendants.
// if lengths match all descendants are evaluated with their corresponding items.
type sliceNode []Mutation

func (n sliceNode) mutate(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
	if value == nil {
		return nil, nil
	} else if reflectutils.GetKind(value) != reflect.Slice {
		return nil, field.TypeInvalid(path, value, "expected a slice")
	} else {
		var out []any
		valueOf := reflect.ValueOf(value)
		for i := range n {
			// TODO: does it make sense to take valueOf.Index(i).Interface() here ?
			if inner, err := n[i].mutate(ctx, path.Index(i), valueOf.Index(i).Interface(), bindings, opts...); err != nil {
				return nil, err
			} else {
				out = append(out, inner)
			}
		}
		return out, nil
	}
}

// scalarNode is a terminal type of mutation.
// it receives a value and compares it with an expected value.
// the expected value can be the result of an expression.
type scalarNode struct {
	rhs any
}

func (n *scalarNode) mutate(ctx context.Context, path *field.Path, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
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
	return rhs, nil
}
