package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/jp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBinding(path *field.Path, bindings binding.Bindings, value any, entry any, opts ...template.Option) binding.Binding {
	return jp.NewLazyBinding(
		func() (any, error) {
			expression := parseExpression(context.TODO(), entry)
			if expression != nil && expression.engine != "" {
				if expression.foreach {
					return nil, field.Invalid(path.Child("variable"), entry, "foreach is not supported in context")
				}
				if expression.binding != "" {
					return nil, field.Invalid(path.Child("variable"), entry, "binding is not supported in context")
				}
				projected, err := template.Execute(context.TODO(), expression.statement, value, bindings, opts...)
				if err != nil {
					return nil, field.InternalError(path.Child("variable"), err)
				}
				return projected, nil
			}
			return entry, nil
		},
	)
}
