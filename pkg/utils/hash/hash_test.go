package hash

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{{
		name: "nil",
		in:   nil,
		want: "",
	}, {
		name: "unmarshalable",
		in:   make(chan int),
		want: "",
	}, {
		name: "string",
		in:   "foo",
		want: Hash("foo"),
	}, {
		name: "map",
		in:   map[string]any{"a": 1},
		want: Hash(map[string]any{"a": 1}),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.in); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash_consistent(t *testing.T) {
	in := map[string]any{"a": 1, "b": "two"}
	if Hash(in) != Hash(in) {
		t.Error("Hash() is not consistent for the same input")
	}
}

func TestHash_different_inputs_different_hash(t *testing.T) {
	if Hash("foo") == Hash("bar") {
		t.Error("Hash() returned the same hash for different inputs")
	}
}
