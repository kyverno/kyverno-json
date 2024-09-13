package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBinding(path *field.Path, bindings binding.Bindings, value any, entry any) binding.Binding {
	return template.NewLazyBinding(
		func() (any, error) {
			expression := ParseExpression(context.TODO(), entry)
			if expression != nil && expression.engine != "" {
				if expression.Foreach {
					return nil, field.Invalid(path.Child("variable"), entry, "foreach is not supported in context")
				}
				if expression.binding != "" {
					return nil, field.Invalid(path.Child("variable"), entry, "binding is not supported in context")
				}
				projected, err := template.Execute(context.Background(), expression.Statement, value, bindings)
				if err != nil {
					return nil, field.InternalError(path.Child("variable"), err)
				}
				return projected, nil
			}
			return entry, nil
		},
	)
}
