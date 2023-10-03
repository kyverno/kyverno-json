package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny_DeepCopyInto(t *testing.T) {
	tests := []struct {
		name string
		in   *Any
		out  *Any
	}{{
		name: "nil",
		in:   &Any{nil},
		out:  &Any{nil},
	}, {
		name: "int",
		in:   &Any{42},
		out:  &Any{nil},
	}, {
		name: "string",
		in:   &Any{"foo"},
		out:  &Any{nil},
	}, {
		name: "slice",
		in:   &Any{[]interface{}{42, "string"}},
		out:  &Any{nil},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.DeepCopyInto(tt.out)
			assert.Equal(t, tt.in, tt.out)
		})
	}
}
