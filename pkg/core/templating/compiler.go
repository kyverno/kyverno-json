package templating

import (
	"reflect"
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/kyverno/kyverno-json/pkg/core/templating/jp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type CelOptions struct {
	FunctionCaller interpreter.FunctionCaller
}

type CompilerOptions struct {
	Cel CelOptions
	Jp  []jp.Option
}

type Compiler struct {
	options CompilerOptions
}

func NewCompiler(options CompilerOptions) Compiler {
	return Compiler{
		options: options,
	}
}

type Program func(any, binding.Bindings) (any, error)

func (c Compiler) Options() CompilerOptions {
	return c.options
}

var bindingsType = cel.OpaqueType("bindings")

type b struct {
	binding.Bindings
}

func (b b) Value() interface{} {
	return b
}
func (b b) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	panic("not required")
}

func (b b) ConvertToType(typeVal ref.Type) ref.Val {
	panic("not required")
}
func (x b) Equal(other ref.Val) ref.Val {
	o, ok := other.Value().(b)
	if !ok {
		return types.ValOrErr(other, "no such overload xxx")
	}
	return types.Bool(o == x)
}

func (b b) Type() ref.Type {
	return bindingsType
}

var newEnv = sync.OnceValues(func() (*cel.Env, error) {
	return cel.NewEnv(
		cel.Variable("object", cel.DynType),
		cel.Variable("bindings", bindingsType),
		cel.Function("resolve",
			cel.MemberOverload("bindings_resolve_string",
				[]*cel.Type{bindingsType, cel.StringType},
				cel.AnyType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					bindings := lhs.(b)
					name := rhs.(types.String)
					value, err := binding.Resolve("$"+string(name), bindings)
					if err != nil {
						return types.WrapErr(err)
					}
					return types.DefaultTypeAdapter.NativeToValue(value)
				}),
			),
		),
	)
})

func (c Compiler) CompileCEL(statement string) (Program, error) {
	env, err := newEnv()
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
	return func(value any, bindings binding.Bindings) (any, error) {
		out, _, err := prg.Eval(
			map[string]interface{}{
				"object":   value,
				"bindings": b{bindings},
			},
		)
		if err != nil {
			return nil, err
		}
		return out.Value(), nil
	}, nil
}

func (c Compiler) CompileJP(statement string) (Program, error) {
	parser := parsing.NewParser()
	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	return func(value any, bindings binding.Bindings) (any, error) {
		return jp.Execute(compiled, value, bindings, c.options.Jp...)
	}, nil
}

func (c Compiler) NewBinding(path *field.Path, value any, bindings binding.Bindings, template any) binding.Binding {
	return jp.NewLazyBinding(
		func() (any, error) {
			switch typed := template.(type) {
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
					projected, err := ExecuteJP(expr.Statement, value, bindings, c)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				case expression.EngineCEL:
					projected, err := ExecuteCEL(expr.Statement, value, bindings, c)
					if err != nil {
						return nil, field.InternalError(path.Child("variable"), err)
					}
					return projected, nil
				default:
					return expr.Statement, nil
				}
			default:
				return typed, nil
			}
		},
	)
}
