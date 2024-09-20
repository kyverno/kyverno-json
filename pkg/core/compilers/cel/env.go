package cel

import (
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

var (
	BindingsType = cel.OpaqueType("bindings")
	DefaultEnv   = sync.OnceValues(func() (*cel.Env, error) {
		return cel.NewEnv(
			cel.Variable("object", cel.DynType),
			cel.Variable("bindings", BindingsType),
			cel.Function("resolve",
				cel.MemberOverload("bindings_resolve_string",
					[]*cel.Type{BindingsType, cel.StringType},
					cel.AnyType,
					cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
						bindings, ok := lhs.(Val[binding.Bindings])
						if !ok {
							return types.ValOrErr(bindings, "invalid bindings type")
						}
						name, ok := rhs.(types.String)
						if !ok {
							return types.ValOrErr(name, "invalid name type")
						}
						value, err := binding.Resolve("$"+string(name), bindings.Unwrap())
						if err != nil {
							return types.WrapErr(err)
						}
						return types.DefaultTypeAdapter.NativeToValue(value)
					}),
				),
			),
		)
	})
)
