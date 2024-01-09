package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Assert(ctx context.Context, path *field.Path, assertion Assertion, value any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	return assertion.assert(ctx, path, value, bindings, opts...)
}
