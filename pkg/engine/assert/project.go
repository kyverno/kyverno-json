package assert

import (
	"fmt"
	"reflect"

	"github.com/eddycharly/json-kyverno/pkg/engine/template"
	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

func project(projection interface{}, value interface{}, bindings binding.Bindings) (interface{}, bool, string, error) {
	expression := parseExpression(projection)
	if expression != nil {
		if expression.engine != "" {
			projected, err := template.Execute(expression.statement, value, bindings)
			if err != nil {
				return nil, false, "", err
			}
			return projected, expression.foreach, expression.binding, nil
		} else {
			if reflectutils.GetKind(value) == reflect.Map {
				return reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.statement)).Interface(), expression.foreach, expression.binding, nil
			}
		}
	}
	if reflectutils.GetKind(value) == reflect.Map {
		return reflect.ValueOf(value).MapIndex(reflect.ValueOf(projection)).Interface(), false, fmt.Sprint(projection), nil
	}
	return value, false, "", nil
}
