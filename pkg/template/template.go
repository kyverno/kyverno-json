package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/eddycharly/tf-kyverno/pkg/apis/v1alpha1"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
)

var (
	variable = regexp.MustCompile(`{{(.*?)}}`)
	inline   = regexp.MustCompile(`^{{(.*)}}$`)
	parser   = parsing.NewParser()
	caller   = interpreter.NewFunctionCaller(functions.GetDefaultFunctions()...)
)

type Template interface {
	String(string) string
	Interface(string) interface{}
}

type template struct {
	data        interface{}
	interpreter interpreter.Interpreter
}

func New(data interface{}, context ...v1alpha1.ContextEntry) Template {
	bindings := binding.NewBindings()
	for _, entry := range context {
		bindings = bindings.Register("$"+entry.Name, entry.Variable.Value)
	}
	return &template{
		data:        data,
		interpreter: interpreter.NewInterpreter(data, caller, bindings),
	}
}

func (t *template) String(in string) string {
	groups := variable.FindAllStringSubmatch(in, -1)
	for _, group := range groups {
		statement := strings.TrimSpace(group[1])
		result, err := t.jp(statement)
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

func (t *template) Interface(in string) interface{} {
	if inline.MatchString(in) {
		in = strings.TrimPrefix(in, "{{")
		in = strings.TrimSuffix(in, "}}")
		statement := strings.TrimSpace(in)
		result, err := t.jp(statement)
		if err != nil {
			return nil
		}
		return result
	} else {
		return t.String(in)
	}
}

func (t *template) jp(statement string) (interface{}, error) {
	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	return t.interpreter.Execute(compiled, t.data)
}
