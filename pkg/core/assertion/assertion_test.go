package assertion

import (
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/stretchr/testify/assert"
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
			compiler := compilers.DefaultCompilers
			parsed, err := Parse(tt.assertion, compiler)
			assert.NoError(t, err)
			got, err := parsed.Assert(nil, tt.value, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		assertion any
		want      field.ErrorList
		wantErr   bool
	}{{
		name: "bad scalar",
		assertion: map[string]any{
			"foo": map[string]any{
				"bar": "~.(`42`)",
			},
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := compilers.DefaultCompilers
			parsed, err := Parse(tt.assertion, compiler)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parsed)
			}
		})
	}
}
