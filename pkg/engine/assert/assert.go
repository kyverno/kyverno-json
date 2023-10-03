package assert

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Validate(assertion interface{}, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return validate(nil, assertion, value, bindings)
}

func Assert(assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assert(nil, assertion, value, bindings)
}

func validate(path *field.Path, assertion interface{}, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assert(path, Parse(assertion), value, bindings)
}

func assert(path *field.Path, assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assertion.assert(path, value, bindings)
}
