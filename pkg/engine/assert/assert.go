package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Validate(ctx context.Context, assertion interface{}, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return validate(ctx, nil, assertion, value, bindings)
}

func Assert(ctx context.Context, assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assert(ctx, nil, assertion, value, bindings)
}

func validate(ctx context.Context, path *field.Path, assertion interface{}, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assert(ctx, path, Parse(ctx, assertion), value, bindings)
}

func assert(ctx context.Context, path *field.Path, assertion Assertion, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	return assertion.assert(ctx, path, value, bindings)
}
