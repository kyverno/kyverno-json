package assert

import (
	"context"
	"errors"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

type projection struct {
	foreach     bool
	foreachName string
	binding     string
	result      any
}

// TODO: remove need for key
func project(ctx context.Context, expression *expression, key any, value any, bindings binding.Bindings, opts ...template.Option) (*projection, error) {
	if expression != nil {
		if expression.engine != "" {
			ast, err := expression.ast()
			if err != nil {
				return nil, err
			}
			projected, err := template.ExecuteAST(ctx, ast, value, bindings, opts...)
			if err != nil {
				return nil, err
			}
			return &projection{
				foreach:     expression.foreach,
				foreachName: expression.foreachName,
				binding:     expression.binding,
				result:      projected,
			}, nil
		} else {
			if value == nil {
				return nil, nil
			} else if reflectutils.GetKind(value) == reflect.Map {
				mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.statement))
				if !mapValue.IsValid() {
					return nil, nil
				}
				return &projection{
					foreach:     expression.foreach,
					foreachName: expression.foreachName,
					binding:     expression.binding,
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
