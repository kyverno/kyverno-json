package cel

import (
	"github.com/google/cel-go/cel"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

func Execute(program cel.Program, value any, bindings binding.Bindings) (any, error) {
	data := map[string]interface{}{
		"object":   value,
		"bindings": NewVal(bindings, BindingsType),
	}
	out, _, err := program.Eval(data)
	if err != nil {
		return nil, err
	}
	return out.Value(), nil
}
