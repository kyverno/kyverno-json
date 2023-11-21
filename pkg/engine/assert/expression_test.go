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
		want *expression
	}{{
		name: "empty",
		in:   "",
		want: nil,
	}, {
		name: "simple field",
		in:   "test",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test",
		},
	}, {
		name: "simple field",
		in:   "(test)",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test",
			engine:      "jp",
		},
	}, {
		name: "nested field",
		in:   "test.test",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test.test",
		},
	}, {
		name: "nested field",
		in:   "(test.test)",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test.test",
			engine:      "jp",
		},
	}, {
		name: "foreach simple field",
		in:   "~.test",
		want: &expression{
			foreach:     true,
			foreachName: "",
			statement:   "test",
		},
	}, {
		name: "foreach simple field",
		in:   "~.(test)",
		want: &expression{
			foreach:     true,
			foreachName: "",
			statement:   "test",
			engine:      "jp",
		},
	}, {
		name: "foreach nested field",
		in:   "~.(test.test)",
		want: &expression{
			foreach:     true,
			foreachName: "",
			statement:   "test.test",
			engine:      "jp",
		},
	}, {
		name: "binding",
		in:   "test->foo",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test",
			binding:     "foo",
		},
	}, {
		name: "binding",
		in:   "(test)->foo",
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test",
			binding:     "foo",
			engine:      "jp",
		},
	}, {
		name: "foreach and binding",
		in:   "~.test->foo",
		want: &expression{
			foreach:     true,
			foreachName: "",
			statement:   "test",
			binding:     "foo",
		},
	}, {
		name: "foreach and binding",
		in:   "~.(test)->foo",
		want: &expression{
			foreach:     true,
			foreachName: "",
			statement:   "test",
			binding:     "foo",
			engine:      "jp",
		},
	}, {
		name: "escape",
		in:   `\~(test)->foo\`,
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "~(test)->foo",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\test\`,
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "test",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\(test)\`,
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "(test)",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `\/test/\`,
		want: &expression{
			foreach:     false,
			foreachName: "",
			statement:   "/test/",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\`,
		want: &expression{
			foreach:     true,
			foreachName: "index",
			statement:   "(test)",
			binding:     "",
		},
	}, {
		name: "escape",
		in:   `~index.\(test)\->name`,
		want: &expression{
			foreach:     true,
			foreachName: "index",
			statement:   "(test)",
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
