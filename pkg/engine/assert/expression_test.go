package assert

import (
	"context"
	"testing"

	tassert "github.com/stretchr/testify/assert"
)

func Test_parseExpressionRegex(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want *Expression
	}{{
		name: "empty",
		in:   "",
		want: nil,
	}, {
		name: "simple field",
		in:   "test",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
		},
	}, {
		name: "simple field",
		in:   "(test)",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			engine:      "jp",
		},
	}, {
		name: "nested field",
		in:   "test.test",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test.test",
		},
	}, {
		name: "nested field",
		in:   "(test.test)",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test.test",
			engine:      "jp",
		},
	}, {
		name: "foreach simple field",
		in:   "~.test",
		want: &Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
		},
	}, {
		name: "foreach simple field",
		in:   "~.(test)",
		want: &Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			engine:      "jp",
		},
	}, {
		name: "foreach nested field",
		in:   "~.(test.test)",
		want: &Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test.test",
			engine:      "jp",
		},
	}, {
		name: "binding",
		in:   "test->foo",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			binding:     "foo",
		},
	}, {
		name: "binding",
		in:   "(test)->foo",
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			binding:     "foo",
			engine:      "jp",
		},
	}, {
		name: "foreach and binding",
		in:   "~.test->foo",
		want: &Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			binding:     "foo",
		},
	}, {
		name: "foreach and binding",
		in:   "~.(test)->foo",
		want: &Expression{
			Foreach:     true,
			ForeachName: "",
			Statement:   "test",
			binding:     "foo",
			engine:      "jp",
		},
	}, {
		name: "escape",
		in:   `\~(test)->foo\`,
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "~(test)->foo",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\test\`,
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "test",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\(test)\`,
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "(test)",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\/test/\`,
		want: &Expression{
			Foreach:     false,
			ForeachName: "",
			Statement:   "/test/",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\`,
		want: &Expression{
			Foreach:     true,
			ForeachName: "index",
			Statement:   "(test)",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\->name`,
		want: &Expression{
			Foreach:     true,
			ForeachName: "index",
			Statement:   "(test)",
			binding:     "name",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExpressionRegex(context.Background(), tt.in)
			tassert.Equal(t, tt.want, got)
		})
	}
}
