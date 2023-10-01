package assert

import (
	"reflect"
	"strings"

	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
)

const (
	foreachMarker          = '~'
	foreachPrefix          = "~"
	expressionPrefix       = "("
	expressionSuffix       = ")"
	legacyExpressionPrefix = "{{"
	legacyExpressionSuffix = "}}"
)

type expression struct {
	foreach   bool
	statement string
	binding   string
	engine    string
}

func parseExpression(value interface{}) *expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	statement := reflect.ValueOf(value).String()
	foreach := false
	binding := ""
	engine := ""
	if statement[0] == foreachMarker {
		foreach = true
	}
	statement = strings.TrimPrefix(statement, foreachPrefix)
	if strings.HasPrefix(statement, legacyExpressionPrefix) {
		statement = strings.TrimPrefix(statement, legacyExpressionPrefix)
		statement = strings.TrimSuffix(statement, legacyExpressionSuffix)
		engine = "jp"
	} else if strings.HasPrefix(statement, expressionPrefix) {
		statement = strings.TrimPrefix(statement, expressionPrefix)
		statement = strings.TrimSuffix(statement, expressionSuffix)
		engine = "jp"
	} else {
		binding = strings.TrimSpace(statement)
	}
	return &expression{
		foreach:   foreach,
		statement: strings.TrimSpace(statement),
		binding:   binding,
		engine:    engine,
	}
}
