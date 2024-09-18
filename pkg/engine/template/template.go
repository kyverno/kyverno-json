package template

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
)

var variable = regexp.MustCompile(`{{(.*?)}}`)

func String(ctx context.Context, in string, value any, bindings binding.Bindings, opts ...Option) string {
	groups := variable.FindAllStringSubmatch(in, -1)
	for _, group := range groups {
		statement := strings.TrimSpace(group[1])
		result, err := ExecuteJP(ctx, statement, value, bindings, opts...)
		if err != nil {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s - %s)", statement, err))
		} else if result == nil {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s not found)", statement))
		} else if result, ok := result.(string); !ok {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s not a string)", statement))
		} else {
			in = strings.ReplaceAll(in, group[0], result)
		}
	}
	return in
}

func ExecuteCEL(ctx context.Context, statement string, value any, bindings binding.Bindings) (any, error) {
	env, err := cel.NewEnv(cel.Variable("object", cel.AnyType))
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
