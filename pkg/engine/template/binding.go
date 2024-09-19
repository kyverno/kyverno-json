package template

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/kyverno/kyverno-json/pkg/jp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBinding(path *field.Path, bindings binding.Bindings, value any, template any, opts ...Option) binding.Binding {
	return jp.NewLazyBinding(
		func() (any, error) {
			switch typed := template.(type) {
			case string:
				expr := expression.Parse(typed)
				if expr.Foreach {
					return nil, field.Invalid(path.Child("variable"), typed, "foreach is not supported in context")
				}
				if expr.Binding != "" {
					return nil, field.Invalid(path.Child("variable"), typed, "binding is not supported in context")
				}
				switch expr.Engine {
				case expression.EngineJP:
					projected, err := ExecuteJP(context.TODO(), expr.Statement, value, bindings, opts...)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				case expression.EngineCEL:
					projected, err := template.ExecuteCEL(context.Background(), expr.Statement, value, bindings)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				default:
					return expr.Statement, nil
				}
			default:
				return typed, nil
			}
		},
	)
}
