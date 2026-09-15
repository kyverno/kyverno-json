package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Expression
	}{{
		name: "empty",
		in:   "",
		want: Expression{},
	}, {
		name: "simple field",
		in:   "test",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
		},
	}, {
		name: "simple field",
		in:   "(test)",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "nested field",
		in:   "test.test",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test.test",
		},
	}, {
		name: "nested field",
		in:   "(test.test)",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test.test",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "Foreach simple field",
		in:   "~.test",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
		},
	}, {
		name: "Foreach simple field",
		in:   "~.(test)",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "Foreach nested field",
		in:   "~.(test.test)",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test.test",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "binding",
		in:   "test->foo",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			Binding:     "foo",
		},
	}, {
		name: "binding",
		in:   "(test)->foo",
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			Binding:     "foo",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "Foreach and binding",
		in:   "~.test->foo",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			Binding:     "foo",
		},
	}, {
		name: "Foreach and binding",
		in:   "~.(test)->foo",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			Binding:     "foo",
			Compiler:    CompilerDefault,
		},
	}, {
		name: "escape",
		in:   `\~(test)->foo\`,
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "~(test)->foo",
			Binding:     "",
		},
	}, {
		name: "escape",
		in:   `\test\`,
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			Binding:     "",
		},
	}, {
		name: "escape",
		in:   `\(test)\`,
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "(test)",
			Binding:     "",
		},
	}, {
		name: "escape",
		in:   `\/test/\`,
		want: Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "/test/",
			Binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\`,
		want: Expression{
			Foreach:     true,
			ForeachName: "index",
			Statement:   "(test)",
			Binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\->name`,
		want: Expression{
			Foreach:     true,
			ForeachName: "index",
			Statement:   "(test)",
			Binding:     "name",
		},
	}, {
		name: "cel, single binding shorthand",
		in:   "(cel;$foo == 'bar')",
		want: Expression{
			Statement: "bindings.resolve('foo') == 'bar'",
			Compiler:  CompilerCEL,
		},
	}, {
		name: "cel, multiple binding shorthands",
		in:   "(cel;$outputs_result.Namespaces.exists(ns, ns.ID==$namespace_id))",
		want: Expression{
			Statement: "bindings.resolve('outputs_result').Namespaces.exists(ns, ns.ID==bindings.resolve('namespace_id'))",
			Compiler:  CompilerCEL,
		},
	}, {
		name: "cel, dollar inside single-quoted string is untouched",
		in:   "(cel;'cost: $5' == $price)",
		want: Expression{
			Statement: "'cost: $5' == bindings.resolve('price')",
			Compiler:  CompilerCEL,
		},
	}, {
		name: "cel, dollar inside double-quoted string is untouched",
		in:   `(cel;"cost: $5" == $price)`,
		want: Expression{
			Statement: `"cost: $5" == bindings.resolve('price')`,
			Compiler:  CompilerCEL,
		},
	}, {
		name: "cel, already-explicit resolve call is left as-is",
		in:   "(cel;bindings.resolve('foo') == 'bar')",
		want: Expression{
			Statement: "bindings.resolve('foo') == 'bar'",
			Compiler:  CompilerCEL,
		},
	}, {
		name: "jp, no shorthand rewrite outside cel",
		in:   "(jp;$foo)",
		want: Expression{
			Statement: "$foo",
			Compiler:  "jp",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
