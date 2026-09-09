package functions

import (
	"testing"
)

func Test_jpfAt(t *testing.T) {
	testCases := []struct {
		name      string
		arguments []any
		expected  any
		wantErr   bool
	}{{
		name:      "valid index",
		arguments: []any{[]any{"a", "b", "c"}, 1.0},
		expected:  "b",
	}, {
		name:      "negative index",
		arguments: []any{[]any{"a", "b", "c"}, -1.0},
		wantErr:   true,
	}, {
		name:      "index equal to length",
		arguments: []any{[]any{"a", "b", "c"}, 3.0},
		wantErr:   true,
	}, {
		name:      "index beyond length",
		arguments: []any{[]any{"a", "b", "c"}, 10.0},
		wantErr:   true,
	}, {
		name:      "empty array",
		arguments: []any{[]any{}, 0.0},
		wantErr:   true,
	}, {
		name:      "first argument not an array",
		arguments: []any{"not an array", 0.0},
		wantErr:   true,
	}, {
		name:      "second argument not a number",
		arguments: []any{[]any{"a"}, "not a number"},
		wantErr:   true,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := jpfAt(tc.arguments)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got result %v", result)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Fatalf("jpfAt(%v) = %v, want %v", tc.arguments, result, tc.expected)
			}
		})
	}
}
