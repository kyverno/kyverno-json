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
			expression := parseExpression(context.TODO(), entry)
			if expression != nil {
				if expression.foreach {
					return nil, field.Invalid(path.Child("variable"), entry, "foreach is not supported in context")
				}
				if expression.binding != "" {
					return nil, field.Invalid(path.Child("variable"), entry, "binding is not supported in context")
				}
				if expression.engine == "jp" {
					projected, err := template.ExecuteJP(context.Background(), expression.statement, value, bindings)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				}
				if expression.engine == "cel" {
					projected, err := template.ExecuteCEL(context.Background(), expression.statement, value, bindings)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				}
			}
			return entry, nil
		},
	)
}
