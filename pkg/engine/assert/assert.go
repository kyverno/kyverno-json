package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Assert(ctx context.Context, assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assert(ctx, nil, assertion, value, bindings)
}

func assert(ctx context.Context, path *field.Path, assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assertion.assert(ctx, path, value, bindings)
}
