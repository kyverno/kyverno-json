package binding

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NewContextBindings(bindings binding.Bindings, value any, entries ...v1alpha1.ContextEntry) binding.Bindings {
	var path *field.Path
	path = path.Child("context")
	for i, entry := range entries {
		bindings = bindings.Register("$"+entry.Name, assert.NewContextBinding(path.Index(i), bindings, value, entry.Variable.Value()))
	}
	return bindings
}
