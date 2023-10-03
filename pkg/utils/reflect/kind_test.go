package reflect

import (
	"reflect"
	"testing"
)

func TestGetKind(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  reflect.Kind
	}{{
		name:  "nil",
		value: nil,
		want:  reflect.Invalid,
	}, {
		name:  "int",
		value: int(42),
		want:  reflect.Int,
	}, {
		name:  "int32",
		value: int32(42),
		want:  reflect.Int32,
	}, {
		name:  "int64",
		value: int64(42),
		want:  reflect.Int64,
	}, {
		name:  "string",
		value: "foo",
		want:  reflect.String,
	}, {
		name:  "map",
		value: map[interface{}]interface{}{},
		want:  reflect.Map,
	}, {
		name:  "slice",
		value: []interface{}{},
		want:  reflect.Slice,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKind(tt.value); got != tt.want {
				t.Errorf("GetKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
