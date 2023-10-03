package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariable_DeepCopy(t *testing.T) {
	tests := []struct {
		name string
		in   *Variable
	}{{
		name: "nil",
		in:   &Variable{nil},
	}, {
		name: "nil",
		in:   nil,
	}, {
		name: "int",
		in:   &Variable{42},
	}, {
		name: "string",
		in:   &Variable{"foo"},
	}, {
		name: "slice",
		in:   &Variable{[]interface{}{42, "string"}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.DeepCopy(); !reflect.DeepEqual(got, tt.in) {
				t.Errorf("Variable.DeepCopy() = %v, want %v", got, tt.in)
			} else if tt.in != nil {
				assert.False(t, got == tt.in)
			}
		})
	}
}
