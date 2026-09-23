package expression

import (
	"regexp"
	"strings"
)

const (
	CompilerCEL     = "cel"
	CompilerJP      = "jp"
	CompilerDefault = "default"
)

var (
	foreachRegex   = regexp.MustCompile(`^~(\w+)?\.(.*)`)
	bindingRegex   = regexp.MustCompile(`(.*)\s*->\s*(\w+)$`)
	escapeRegex    = regexp.MustCompile(`^\\(.+)\\$`)
	compilerRegex  = regexp.MustCompile(`^\((?:(\w+);)?(.+)\)$`)
	celBindingName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*`)
)

type Expression struct {
	Foreach     bool
	ForeachName string
	Statement   string
	Binding     string
	Compiler    string
}

func Parse(in string) (expression Expression) {
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
		if match := compilerRegex.FindStringSubmatch(in); match != nil {
			expression.Compiler = match[1]
			// account for default engine
			if expression.Compiler == "" {
				expression.Compiler = CompilerDefault
			}
			in = match[2]
		}
	}
	// 4. for CEL statements, rewrite the `$name` shorthand into an explicit
	// `bindings.resolve('name')` call so users don't have to spell that out
	// themselves
	if expression.Compiler == CompilerCEL {
		in = celBindingShorthand(in)
	}
	// 5. assign statement
	expression.Statement = in
	// 6. done
	return
}

// celBindingShorthand rewrites `$name` occurrences in a CEL statement into
// `bindings.resolve('name')`, skipping anything inside a single- or
// double-quoted string literal so a literal `$` in a string (e.g. "cost:
// $5") is left untouched.
func celBindingShorthand(in string) string {
	var out strings.Builder
	var quote rune
	runes := []rune(in)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if quote != 0 {
			out.WriteRune(c)
			if c == '\\' && i+1 < len(runes) {
				// preserve the escaped character verbatim so it can't
				// prematurely close the string literal
				i++
				out.WriteRune(runes[i])
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			out.WriteRune(c)
			continue
		}
		if c == '$' {
			if name := celBindingName.FindString(string(runes[i+1:])); name != "" {
				out.WriteString("bindings.resolve('")
				out.WriteString(name)
				out.WriteString("')")
				i += len(name)
				continue
			}
		}
		out.WriteRune(c)
	}
	return out.String()
}
