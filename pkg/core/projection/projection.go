package projection

import (
	"errors"
	"reflect"
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

type Handler = func(value any, bindings binding.Bindings) (any, bool, error)

type Info struct {
	Foreach     bool
	ForeachName string
	Binding     string
}

type Projection struct {
	Info
	Handler
}

func Parse(in any, compiler compilers.Compilers) (projection Projection) {
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
			parse := sync.OnceValues(func() (compilers.Program, error) {
				return compiler.Jp.Compile(expr.Statement)
			})
			projection.Handler = func(value any, bindings binding.Bindings) (any, bool, error) {
				program, err := parse()
				if err != nil {
					return nil, false, err
				}
				projected, err := program(value, bindings)
				if err != nil {
					return nil, false, err
				}
				return projected, true, err
			}
		case expression.EngineCEL:
			projection.Handler = func(value any, bindings binding.Bindings) (any, bool, error) {
				program, err := compiler.Cel.Compile(expr.Statement)
				if err != nil {
					return nil, false, err
				}
				projected, err := program(value, bindings)
				if err != nil {
					return nil, false, err
				}
				return projected, true, nil
			}
		default:
			projection.Handler = func(value any, bindings binding.Bindings) (any, bool, error) {
				if value == nil {
					return nil, false, nil
				}
				if reflectutils.GetKind(value) == reflect.Map {
					value := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expr.Statement))
					if !value.IsValid() {
						return nil, false, nil
					}
					return value.Interface(), true, nil
				}
				return nil, false, errors.New("projection not recognized")
			}
		}
	default:
		// 1. compute the projection func
		projection.Handler = func(value any, bindings binding.Bindings) (any, bool, error) {
			if value == nil {
				return nil, false, nil
			}
			if reflectutils.GetKind(value) == reflect.Map {
				mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(typed))
				if !mapValue.IsValid() {
					return nil, false, nil
				}
				return mapValue.Interface(), true, nil
			}
			return nil, false, errors.New("projection not recognized")
		}
	}
	return
}
