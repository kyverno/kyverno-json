package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
)

var (
	variable = regexp.MustCompile(`{{(.*?)}}`)
	parser   = parsing.NewParser()
	caller   = interpreter.NewFunctionCaller(GetFunctions()...)
)

func String(in string, value interface{}, bindings binding.Bindings) string {
	groups := variable.FindAllStringSubmatch(in, -1)
	for _, group := range groups {
		statement := strings.TrimSpace(group[1])
		result, err := Execute(statement, value, bindings)
		if err != nil {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s - %s)", statement, err))
		} else if result == nil {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s not found)", statement))
		} else if result, ok := result.(string); !ok {
			in = strings.ReplaceAll(in, group[0], fmt.Sprintf("ERR (%s not a string)", statement))
		} else {
			in = strings.ReplaceAll(in, group[0], result)
		}
	}
	return in
}

func Execute(statement string, value interface{}, bindings binding.Bindings) (interface{}, error) {
	interpreter := interpreter.NewInterpreter(nil, caller, bindings)
	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	return interpreter.Execute(compiled, value)
}
