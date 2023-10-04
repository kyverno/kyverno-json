package assert

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBindings(bindings binding.Bindings, value interface{}, entries ...v1alpha1.ContextEntry) binding.Bindings {
	var path *field.Path
	path = path.Child("context")
	for i, entry := range entries {
		bindings = bindings.Register("$"+entry.Name, NewContextBinding(path.Index(i), bindings, value, entry))
	}
	return bindings
}

func NewContextBinding(path *field.Path, bindings binding.Bindings, value interface{}, entry v1alpha1.ContextEntry) binding.Binding {
	return template.NewLazyBinding(
		func() (interface{}, error) {
			expression := parseExpression(entry.Variable.Value)
			if expression != nil && expression.engine != "" {
				if expression.foreachName != "" {
					return nil, field.Invalid(path.Child("variable", "value"), entry.Variable.Value, "foreach is not supported in context")
				}
				if expression.binding != "" {
					return nil, field.Invalid(path.Child("variable", "value"), entry.Variable.Value, "binding is not supported in context")
				}
				projected, err := template.Execute(expression.statement, value, bindings)
				if err != nil {
					return nil, field.InternalError(path.Child("variable", "value"), err)
				}
				return projected, nil
			}
			return entry.Variable.Value, nil
		},
	)
}
