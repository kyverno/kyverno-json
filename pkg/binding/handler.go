package binding

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type Bindings struct {
	jpbindings  binding.Bindings
	data        map[string]any
	celEnvNames []cel.EnvOption
}

func New() Bindings {
	return Bindings{
		celEnvNames: make([]cel.EnvOption, 0),
		data:        make(map[string]any),
		jpbindings:  binding.NewBindings(),
	}
}

func (b Bindings) Register(key string, value any) error {
	b.jpbindings = b.jpbindings.Register(key, binding.NewBinding(value))
	if _, found := b.data[key]; found {
		return fmt.Errorf("binding already exists: %s", key)
	} else {
		b.data[key] = value
		b.celEnvNames = append(b.celEnvNames, cel.Variable(key, cel.AnyType))
	}
	return nil
}

func (b Bindings) JmespathBinding() binding.Bindings {
	return b.jpbindings
}

func (b Bindings) CELEnv() (*cel.Env, map[string]any, error) {
	env, err := cel.NewEnv(b.celEnvNames...)
	if err != nil {
		return nil, nil, err
	}

	return env, b.data, nil
}
