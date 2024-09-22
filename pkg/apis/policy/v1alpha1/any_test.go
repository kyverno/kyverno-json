package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny_DeepCopyInto(t *testing.T) {
	tests := []struct {
		name string
		in   Any
	}{{
		name: "nil",
		in:   NewAny(nil),
	}, {
		name: "int",
		in:   NewAny(42),
	}, {
		name: "string",
		in:   NewAny("foo"),
	}, {
		name: "slice",
		in:   NewAny([]any{42, "string"}),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out Any
			tt.in.DeepCopyInto(&out)
			assert.Equal(t, tt.in, out)
		})
	}
	{
		inner := map[string]any{
			"foo": 42,
		}
		in := NewAny(map[string]any{"inner": inner})
		out := in.DeepCopy()
		inPtr := in._value.(map[string]any)["inner"].(map[string]any)
		inPtr["foo"] = 55
		outPtr := out._value.(map[string]any)["inner"].(map[string]any)
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
			a := NewAny(tt.value)
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
		want    Any
		wantErr bool
	}{{
		name:    "nil",
		data:    []byte("null"),
		want:    NewAny(nil),
		wantErr: false,
	}, {
		name:    "int",
		data:    []byte("42"),
		want:    NewAny(int64(42)),
		wantErr: false,
	}, {
		name:    "string",
		data:    []byte(`"foo"`),
		want:    NewAny("foo"),
		wantErr: false,
	}, {
		name:    "map",
		data:    []byte(`{"foo":42}`),
		want:    NewAny(map[string]any{"foo": int64(42)}),
		wantErr: false,
	}, {
		name:    "error",
		data:    []byte(`{"foo":`),
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Any
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
