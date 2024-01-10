package reflect

import (
	"reflect"
	"testing"
)

func TestToNumber(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		want   float64
		wantOk bool
	}{{
		name:   "nil",
		value:  nil,
		want:   0,
		wantOk: false,
	}, {
		name:   "float32",
		value:  float32(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "float64",
		value:  float64(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "int32",
		value:  int32(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "int64",
		value:  int64(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "uint32",
		value:  uint32(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "uint64",
		value:  uint64(42),
		want:   42,
		wantOk: true,
	}, {
		name:   "string",
		value:  "foo",
		want:   0,
		wantOk: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ToNumber(reflect.ValueOf(tt.value))
			if got != tt.want {
				t.Errorf("ToNumber() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("ToNumber() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
