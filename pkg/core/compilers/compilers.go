package compilers

import (
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/cel"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var DefaultCompiler = Compilers{
	Jp:  jp.NewCompiler(),
	Cel: cel.NewCompiler(),
}

type Compilers struct {
	Jp  jp.Compiler
	Cel cel.Compiler
}

func (c Compilers) NewBinding(path *field.Path, value any, bindings binding.Bindings, template any) binding.Binding {
	return binding.NewDelegate(
		sync.OnceValues(
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
						projected, err := Execute(expr.Statement, value, bindings, c.Jp)
						if err != nil {
							return nil, field.InternalError(path.Child("variable"), err)
						}
						return projected, nil
					case expression.EngineCEL:
						projected, err := Execute(expr.Statement, value, bindings, c.Cel)
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
		),
	)
}
