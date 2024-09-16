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
		in:   &Any{[]any{42, "string"}},
		out:  &Any{nil},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.DeepCopyInto(tt.out)
			assert.Equal(t, tt.in, tt.out)
		})
	}
	{
		inner := map[string]any{
			"foo": 42,
		}
		in := Any{map[string]any{"inner": inner}}
		out := in.DeepCopy()
		inPtr := in.Value.(map[string]any)["inner"].(map[string]any)
		inPtr["foo"] = 55
		outPtr := out.Value.(map[string]any)["inner"].(map[string]any)
		assert.NotEqual(t, inPtr, outPtr)
	}
}

func TestAny_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    []byte
		wantErr bool
	}{{
		name:    "nil",
		value:   nil,
		want:    []byte("null"),
		wantErr: false,
	}, {
		name:    "int",
		value:   42,
		want:    []byte("42"),
		wantErr: false,
	}, {
		name:    "string",
		value:   "foo",
		want:    []byte(`"foo"`),
		wantErr: false,
	}, {
		name:    "map",
		value:   map[string]any{"foo": 42},
		want:    []byte(`{"foo":42}`),
		wantErr: false,
	}, {
		name:    "error",
		value:   func() {},
		want:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Any{
				Value: tt.value,
			}
			got, err := a.MarshalJSON()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAny_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *Any
		wantErr bool
	}{{
		name:    "nil",
		data:    []byte("null"),
		want:    &Any{},
		wantErr: false,
	}, {
		name:    "int",
		data:    []byte("42"),
		want:    &Any{Value: 42.0},
		wantErr: false,
	}, {
		name:    "string",
		data:    []byte(`"foo"`),
		want:    &Any{Value: "foo"},
		wantErr: false,
	}, {
		name:    "map",
		data:    []byte(`{"foo":42}`),
		want:    &Any{Value: map[string]any{"foo": 42.0}},
		wantErr: false,
	}, {
		name:    "error",
		data:    []byte(`{"foo":`),
		want:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Any{}
			err := a.UnmarshalJSON(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, a)
			}
		})
	}
}
