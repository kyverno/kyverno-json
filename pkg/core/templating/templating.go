package templating

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

func Execute(statement string, value any, bindings binding.Bindings, compiler compilers.Compiler) (any, error) {
	program, err := compiler.Compile(statement)
	if err != nil {
		return nil, err
	}
	return program(value, bindings)
}
