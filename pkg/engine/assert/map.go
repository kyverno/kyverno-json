package assert

import (
	"fmt"
	"reflect"

	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// mapNode is the assertion type represented by a map.
// it is reponsible for projecting the analysed resource and passing the result to the descendant
type mapNode map[interface{}]Assertion

func (n mapNode) assert(path *field.Path, value interface{}, bindings jpbinding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	for k, v := range n {
		projected, foreach, binding, err := project(k, value, bindings)
		if err != nil {
			return nil, field.InternalError(path.Child(fmt.Sprint(k)), err)
		} else {
			if binding != "" {
				bindings = bindings.Register("$"+binding, jpbinding.NewBinding(projected))
			}
			if foreach != "" {
				projectedKind := reflectutils.GetKind(projected)
				if projectedKind == reflect.Slice {
					valueOf := reflect.ValueOf(projected)
					for i := 0; i < valueOf.Len(); i++ {
						bindings := bindings.Register("$"+foreach, jpbinding.NewBinding(i))
						if _errs, err := v.assert(path.Child(fmt.Sprint(k)).Index(i), valueOf.Index(i).Interface(), bindings); err != nil {
							return nil, err
						} else {
							errs = append(errs, _errs...)
						}
					}
				} else if projectedKind == reflect.Map {
					iter := reflect.ValueOf(projected).MapRange()
					for iter.Next() {
						key := iter.Key().Interface()
						bindings := bindings.Register("$"+foreach, jpbinding.NewBinding(key))
						if _errs, err := v.assert(path.Child(fmt.Sprint(k)).Key(fmt.Sprint(key)), iter.Value().Interface(), bindings); err != nil {
							return nil, err
						} else {
							errs = append(errs, _errs...)
						}
					}
				} else {
					return nil, field.TypeInvalid(path.Child(fmt.Sprint(k)), projected, "expected a slice or a map")
				}
			} else {
				if _errs, err := v.assert(path.Child(fmt.Sprint(k)), projected, bindings); err != nil {
					return nil, err
				} else {
					errs = append(errs, _errs...)
				}
			}
		}
	}
	return errs, nil
}
