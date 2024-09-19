package message

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var variable = regexp.MustCompile(`{{(.*?)}}`)

type Message interface {
	Original() string
	Format(any, binding.Bindings, ...template.Option) string
}

type substitution = func(string, any, binding.Bindings, ...template.Option) string

type message struct {
	original      string
	substitutions []substitution
}

func (m *message) Original() string {
	return m.original
}

func (m *message) Format(value any, bindings binding.Bindings, opts ...template.Option) string {
	out := m.original
	for _, substitution := range m.substitutions {
		out = substitution(out, value, bindings, opts...)
	}
	return out
}

func Parse(in string) *message {
	groups := variable.FindAllStringSubmatch(in, -1)
	var substitutions []func(string, any, binding.Bindings, ...template.Option) string
	for _, group := range groups {
		statement := strings.TrimSpace(group[1])
		parse := sync.OnceValues(func() (parsing.ASTNode, error) {
			parser := parsing.NewParser()
			return parser.Parse(statement)
		})
		evaluate := func(value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
			ast, err := parse()
			if err != nil {
				return nil, err
			}
			return template.ExecuteAST(context.TODO(), ast, value, bindings, opts...)
		}
		placeholder := group[0]
		substitutions = append(substitutions, func(out string, value any, bindings binding.Bindings, opts ...template.Option) string {
			result, err := evaluate(value, bindings, opts...)
			if err != nil {
				out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("ERR (%s - %s)", statement, err))
			} else if result == nil {
				out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("ERR (%s not found)", statement))
			} else if result, ok := result.(string); !ok {
				out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("ERR (%s not a string)", statement))
			} else {
				out = strings.ReplaceAll(out, placeholder, result)
			}
			return out
		})
	}
	return &message{
		original:      in,
		substitutions: substitutions,
	}
}
