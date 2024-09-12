package assert

import (
	"context"
	"reflect"
	"regexp"

	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

var (
	foreachRegex = regexp.MustCompile(`^~(\w+)?\.(.*)`)
	bindingRegex = regexp.MustCompile(`(.*)\s*->\s*(\w+)$`)
	escapeRegex  = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex  = regexp.MustCompile(`^\((?:(\w+):)?(.+)\)$`)
)

type Expression struct {
	Foreach     bool
	ForeachName string
	Statement   string
	binding     string
	engine      string
}

func parseExpressionRegex(_ context.Context, in string) *Expression {
	expression := &Expression{}
	// 1. match foreach
	if match := foreachRegex.FindStringSubmatch(in); match != nil {
		expression.Foreach = true
		expression.ForeachName = match[1]
		in = match[2]
	}
	// 2. match binding
	if match := bindingRegex.FindStringSubmatch(in); match != nil {
		expression.binding = match[2]
		in = match[1]
	}
	// 3. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			expression.engine = match[1]
			// account for default engine
			if expression.engine == "" {
				expression.engine = "jp"
			}
			in = match[2]
		}
	}
	// parse statement
	expression.Statement = in
	if expression.Statement == "" {
		return nil
	}
	return expression
}

func ParseExpression(ctx context.Context, value any) *Expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	return parseExpressionRegex(ctx, reflect.ValueOf(value).String())
}
