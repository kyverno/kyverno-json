package assert

import (
	"reflect"
	"regexp"
	"strings"

	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

const (
	defaultForeachVariable = "index"
	expressionPrefix       = "("
	expressionSuffix       = ")"
	legacyExpressionPrefix = "{{"
	legacyExpressionSuffix = "}}"
)

var (
	foreachRegex = regexp.MustCompile(`^~(?:(\w*)\.)?`)
	bindingRegex = regexp.MustCompile(`@(\w*)$`)
)

type expression struct {
	foreach   string
	statement string
	binding   string
	engine    string
}

func parseExpression(value interface{}) *expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	statement := reflect.ValueOf(value).String()
	foreach := ""
	binding := ""
	engine := ""
	if match := foreachRegex.FindStringSubmatch(statement); match != nil {
		foreach = match[1]
		if foreach == "" {
			foreach = defaultForeachVariable
		}
		statement = strings.TrimPrefix(statement, match[0])
	}
	if match := bindingRegex.FindStringSubmatch(statement); match != nil {
		binding = match[1]
		statement = strings.TrimSuffix(statement, match[0])
	}
	if strings.HasPrefix(statement, legacyExpressionPrefix) {
		statement = strings.TrimPrefix(statement, legacyExpressionPrefix)
		statement = strings.TrimSuffix(statement, legacyExpressionSuffix)
		engine = "jp"
	} else if strings.HasPrefix(statement, expressionPrefix) {
		statement = strings.TrimPrefix(statement, expressionPrefix)
		statement = strings.TrimSuffix(statement, expressionSuffix)
		engine = "jp"
	} else if binding == "" {
		binding = strings.TrimSpace(statement)
	}
	return &expression{
		foreach:   foreach,
		statement: strings.TrimSpace(statement),
		binding:   binding,
		engine:    engine,
	}
}
