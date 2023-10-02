package assert

import (
	"fmt"
	"reflect"

	"github.com/eddycharly/json-kyverno/pkg/engine/template"
	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

func project(projection interface{}, value interface{}, bindings binding.Bindings) (interface{}, string, string, error) {
	expression := parseExpression(projection)
	if expression != nil {
		if expression.engine != "" {
			projected, err := template.Execute(expression.statement, value, bindings)
			if err != nil {
				return nil, "", "", err
			}
			return projected, expression.foreach, expression.binding, nil
		} else {
			if reflectutils.GetKind(value) == reflect.Map {
				projected := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.statement))
				if !projected.IsValid() {
					return nil, "", "", fmt.Errorf("failed to find the map index `%s`", expression.statement)
				}
				return projected.Interface(), expression.foreach, expression.binding, nil
			}
		}
	}
	if reflectutils.GetKind(value) == reflect.Map {
		projected := reflect.ValueOf(value).MapIndex(reflect.ValueOf(projection))
		if !projected.IsValid() {
			return nil, "", "", fmt.Errorf("failed to find the map index `%v`", projection)
		}
		return projected.Interface(), "", fmt.Sprint(projection), nil
	}
	// TODO is this an error ?
	return value, "", "", nil
}
