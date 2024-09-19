package expression

import (
	"regexp"
)

var (
	foreachRegex = regexp.MustCompile(`^~(\w+)?\.(.*)`)
	bindingRegex = regexp.MustCompile(`(.*)\s*->\s*(\w+)$`)
	escapeRegex  = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex  = regexp.MustCompile(`^\((?:(\w+);)?(.+)\)$`)
)

const (
	EngineJP      = "jp"
	EngineCEL     = "cel"
	EngineDefault = EngineJP
)

type Expression struct {
	Foreach     bool
	ForeachName string
	Statement   string
	Binding     string
	Engine      string
}

func Parse(in string) Expression {
	expression := Expression{}
	// 1. match foreach
	if match := foreachRegex.FindStringSubmatch(in); match != nil {
		expression.Foreach = true
		expression.ForeachName = match[1]
		in = match[2]
	}
	// 2. match binding
	if match := bindingRegex.FindStringSubmatch(in); match != nil {
		expression.Binding = match[2]
		in = match[1]
	}
	// 3. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			expression.Engine = match[1]
			// account for default engine
			if expression.Engine == "" {
				expression.Engine = EngineDefault
			}
			in = match[2]
		}
	}
	// 4. assign statement
	expression.Statement = in
	// 5. done
	return expression
}
