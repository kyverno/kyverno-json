package assert

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Assert(assertion interface{}, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return Validate(Parse(assertion), value, bindings)
}
