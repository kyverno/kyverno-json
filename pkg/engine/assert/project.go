package assert

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

type projection struct {
	foreach     bool
	foreachName string
	binding     string
	result      interface{}
}

func project(ctx context.Context, key interface{}, value interface{}, bindings binding.Bindings) (*projection, error) {
	expression := parseExpression(ctx, key)
	if expression != nil {
		if expression.engine != "" {
			projected, err := template.Execute(ctx, expression.statement, value, bindings)
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
			if reflectutils.GetKind(value) == reflect.Map {
				projected := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.statement))
				if !projected.IsValid() {
					return nil, fmt.Errorf("failed to find the map index `%s`", expression.statement)
				}
				return &projection{
					foreach:     expression.foreach,
					foreachName: expression.foreachName,
					binding:     expression.binding,
					result:      projected.Interface(),
				}, nil
			}
		}
	}
	if reflectutils.GetKind(value) == reflect.Map {
		projected := reflect.ValueOf(value).MapIndex(reflect.ValueOf(key))
		if !projected.IsValid() {
			return nil, fmt.Errorf("failed to find the map index `%v`", key)
		}
		return &projection{
			result: projected.Interface(),
		}, nil
	}
	// TODO is this an error ?
	return &projection{
		result: value,
	}, nil
}
