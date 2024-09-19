package template

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
)

func Execute(ctx context.Context, statement string, value any, bindings binding.Bindings, opts ...Option) (any, error) {
	parser := parsing.NewParser()
	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	return ExecuteAST(ctx, compiled, value, bindings, opts...)
}

func ExecuteAST(ctx context.Context, ast parsing.ASTNode, value any, bindings binding.Bindings, opts ...Option) (any, error) {
	o := buildOptions(opts...)
	vm := interpreter.NewInterpreter(nil, bindings)
	return vm.Execute(ast, value, interpreter.WithFunctionCaller(o.functionCaller))
}
