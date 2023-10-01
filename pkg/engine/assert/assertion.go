package assert

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Assertion interface {
	assert(*field.Path, interface{}, binding.Bindings) (field.ErrorList, error)
}
