package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/jp"
	"github.com/kyverno/kyverno-json/pkg/syntax/expression"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBinding(path *field.Path, bindings binding.Bindings, value any, entry any, opts ...template.Option) binding.Binding {
	return jp.NewLazyBinding(
		func() (any, error) {
			switch typed := entry.(type) {
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
					projected, err := template.ExecuteJP(context.TODO(), expr.Statement, value, bindings, opts...)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				case expression.EngineCEL:
					return nil, field.Invalid(path.Child("variable"), expr.Engine, "engine not supported")
				default:
					return expr.Statement, nil
				}
			default:
				return typed, nil
			}
		},
	)
}
