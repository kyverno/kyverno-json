package assert

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Validate(assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assertion.assert(nil, value, bindings)
}
