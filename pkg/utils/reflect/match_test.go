package reflect

import (
	"testing"
)

func TestMatchScalar(t *testing.T) {
	tests := []struct {
		name     string
		expected any
		actual   any
		want     bool
		err      bool
	}{
		{
			name:     "string-string same",
			expected: "str",
			actual:   "str",
			want:     true,
			err:      false,
		},
		{
			name:     "string-string different",
			expected: "str",
			actual:   "string",
			want:     false,
			err:      false,
		},
		{
			name:     "int64-float64 same",
			expected: int64(12),
			actual:   float64(12.0),
			want:     true,
			err:      false,
		},
		{
			name:     "int64-float64 different",
			expected: int64(12),
			actual:   float64(13),
			want:     false,
			err:      false,
		},
		{
			name:     "int64-float64 with decimals",
			expected: int64(12),
			actual:   float64(12.2),
			want:     false,
			err:      false,
		},
		{
			name:     "map-slice",
			expected: make(map[string]interface{}),
			actual:   []string{},
			want:     false,
			err:      true,
		},
		{
			name:     "bool-int",
			expected: true,
			actual:   1,
			want:     false,
			err:      true,
		},
		{
			name:     "string-int",
			expected: "13",
			actual:   13,
			want:     false,
			err:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MatchScalar(tt.expected, tt.actual); got != tt.want {
				t.Errorf("MatchScalar() = %v, want %v", got, tt.want)
			} else if err != nil && !tt.err {
				t.Errorf("MatchScalar() err = %v", err)
			}
		})
	}
}
