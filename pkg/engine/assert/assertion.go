package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Assertion interface {
	assert(context.Context, *field.Path, interface{}, binding.Bindings) (field.ErrorList, error)
}
