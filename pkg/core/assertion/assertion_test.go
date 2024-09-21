package assertion

import (
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	tassert "github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestAssert(t *testing.T) {
	tests := []struct {
		name      string
		assertion any
		value     any
		bindings  binding.Bindings
		want      field.ErrorList
		wantErr   bool
	}{{
		name: "nil vs empty object",
		assertion: map[string]any{
			"foo": map[string]any{},
		},
		value: map[string]any{
			"foo": nil,
		},
		want: field.ErrorList{
			&field.Error{
				Type:   field.ErrorTypeInvalid,
				Field:  "foo",
				Detail: "invalid value, must not be null",
			},
		},
		wantErr: false,
	}, {
		name: "not nil vs empty object",
		assertion: map[string]any{
			"foo": map[string]any{},
		},
		value: map[string]any{
			"foo": map[string]any{
				"bar": 42,
			},
		},
		want:    nil,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := compilers.DefaultCompiler
			parsed, err := Parse(tt.assertion, compiler, expression.CompilerJP)
			tassert.NoError(t, err)
			got, err := parsed.Assert(nil, tt.value, tt.bindings)
			if tt.wantErr {
				tassert.Error(t, err)
			} else {
				tassert.NoError(t, err)
			}
			tassert.Equal(t, tt.want, got)
		})
	}
}
