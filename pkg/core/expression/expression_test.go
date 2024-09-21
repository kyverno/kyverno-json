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
			Compiler:    CompilerJP,
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
			Compiler:    CompilerJP,
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
			Compiler:    CompilerJP,
		},
	}, {
		name: "Foreach nested field",
		in:   "~.(test.test)",
		want: Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test.test",
			Compiler:    CompilerJP,
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
			Compiler:    CompilerJP,
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
			Compiler:    CompilerJP,
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(CompilerJP, tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
