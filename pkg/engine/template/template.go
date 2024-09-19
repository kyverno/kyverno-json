package template

import (
	"context"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
)

func ExecuteCEL(ctx context.Context, statement string, value any, bindings binding.Bindings) (any, error) {
	env, err := cel.NewEnv(
		cel.Variable("object", cel.AnyType),
		cel.Function(
			"binding",
			cel.Overload(
				"binding",
				[]*types.Type{types.StringType},
				types.AnyType,
				cel.UnaryBinding(
					func(arg ref.Val) ref.Val {
						name, ok := arg.(types.String)
						if !ok {
							return types.MaybeNoSuchOverloadErr(arg)
						}
						value, err := binding.Resolve("$"+string(name), bindings)
						if err != nil {
							return types.WrapErr(err)
						}
						return types.DefaultTypeAdapter.NativeToValue(value)
					},
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}
	ast, iss := env.Compile(statement)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := prg.Eval(
		map[string]interface{}{
			"object": value,
		},
	)
	if err != nil {
		return nil, err
	}
	return out.Value(), nil
}

func ExecuteJP(ctx context.Context, statement string, value any, bindings binding.Bindings, opts ...Option) (any, error) {
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
