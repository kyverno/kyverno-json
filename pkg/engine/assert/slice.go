package assert

import (
	"reflect"

	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// sliceNode is the assertion type represented by a slice.
// it first compares the length of the analysed resource with the length of the descendants.
// if lengths match all descendants are evaluated with their corresponding items.
type sliceNode []Assertion

func (n sliceNode) assert(path *field.Path, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	var errs field.ErrorList
	if reflectutils.GetKind(value) != reflect.Slice {
		return nil, field.TypeInvalid(path, value, "expected a slice")
	} else {
		valueOf := reflect.ValueOf(value)
		if valueOf.Len() != len(n) {
			errs = append(errs, field.Invalid(path, value, "lengths of slices don't match"))
		} else {
			for i := range n {
				if _errs, err := n[i].assert(path.Index(i), valueOf.Index(i).Interface(), bindings); err != nil {
					return nil, err
				} else {
					errs = append(errs, _errs...)
				}
			}
		}
	}
	return errs, nil
}
