package functions

import (
	"testing"
)

func Test_jpfAt(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{{
		name: "valid index",
		args: []any{[]any{"a", "b"}, float64(1)},
		want: "b",
	}, {
		name: "first element",
		args: []any{[]any{"a", "b"}, float64(0)},
		want: "a",
	}, {
		name:    "out of range",
		args:    []any{[]any{"a", "b"}, float64(999)},
		wantErr: true,
	}, {
		name:    "negative index",
		args:    []any{[]any{"a", "b"}, float64(-1)},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpfAt(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("jpfAt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != nil {
					t.Fatalf("jpfAt() returned a value %v on error, want nil", got)
				}
				return
			}
			if got != tt.want {
				t.Fatalf("jpfAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
