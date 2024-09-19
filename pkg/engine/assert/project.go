package assert

import (
	"context"
	"errors"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/syntax/expression"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

type projection struct {
	foreach     bool
	foreachName string
	binding     string
	result      any
}

// TODO: remove need for key
func project(ctx context.Context, expression *expression.Expression, key any, value any, bindings binding.Bindings, opts ...template.Option) (*projection, error) {
	if expression != nil {
		if expression.Engine != "" {
			projected, err := template.ExecuteJP(ctx, expression.Statement, value, bindings, opts...)
			if err != nil {
				return nil, err
			}
			return &projection{
				foreach:     expression.Foreach,
				foreachName: expression.ForeachName,
				binding:     expression.Binding,
				result:      projected,
			}, nil
		} else {
			if value == nil {
				return nil, nil
			} else if reflectutils.GetKind(value) == reflect.Map {
				mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.Statement))
				if !mapValue.IsValid() {
					return nil, nil
				}
				return &projection{
					foreach:     expression.Foreach,
					foreachName: expression.ForeachName,
					binding:     expression.Binding,
					result:      mapValue.Interface(),
				}, nil
			}
		}
	}
	if reflectutils.GetKind(value) == reflect.Map {
		mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(key))
		if !mapValue.IsValid() {
			return nil, nil
		}
		return &projection{
			result: mapValue.Interface(),
		}, nil
	}
	return nil, errors.New("projection not recognized")
}
