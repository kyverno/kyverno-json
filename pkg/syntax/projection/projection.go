package projection

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/syntax/expression"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

type Handler = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error)

type Info struct {
	Foreach     bool
	ForeachName string
	Binding     string
}

type Projection struct {
	Info
	Handler
}

func Parse(in any) (projection Projection) {
	switch typed := in.(type) {
	case string:
		// 1. if we have a string, parse the expression
		expr := expression.Parse(typed)
		// 2. record projection infos
		projection.Foreach = expr.Foreach
		projection.ForeachName = expr.ForeachName
		projection.Binding = expr.Binding
		// 3. compute the projection func
		switch expr.Engine {
		case expression.EngineJP:
			parse := sync.OnceValues(func() (parsing.ASTNode, error) {
				parser := parsing.NewParser()
				return parser.Parse(expr.Statement)
			})
			projection.Handler = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				ast, err := parse()
				if err != nil {
					return nil, err
				}
				return template.ExecuteAST(ctx, ast, value, bindings, opts...)
			}
		case expression.EngineCEL:
			panic("engine not supported")
		default:
			projection.Handler = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
				if value == nil {
					return nil, nil
				}
				if reflectutils.GetKind(value) == reflect.Map {
					value := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expr.Statement))
					if !value.IsValid() {
						return nil, nil
					}
					return value.Interface(), nil
				}
				return nil, errors.New("projection not recognized")
			}
		}
	default:
		// 1. compute the projection func
		projection.Handler = func(ctx context.Context, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
			if value == nil {
				return nil, nil
			}
			if reflectutils.GetKind(value) == reflect.Map {
				mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(typed))
				if !mapValue.IsValid() {
					return nil, nil
				}
				return mapValue.Interface(), nil
			}
			return nil, errors.New("projection not recognized")
		}
	}
	return
}
