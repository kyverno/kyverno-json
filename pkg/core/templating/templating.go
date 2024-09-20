package templating

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

func ExecuteJP(statement string, value any, bindings binding.Bindings, compiler Compiler) (any, error) {
	program, err := compiler.CompileJP(statement)
	if err != nil {
		return nil, err
	}
	return program(value, bindings)
}

func ExecuteCEL(statement string, value any, bindings binding.Bindings, compiler Compiler) (any, error) {
	program, err := compiler.CompileCEL(statement)
	if err != nil {
		return nil, err
	}
	return program(value, bindings)
}
