package projection

import (
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/stretchr/testify/assert"
)

func TestParseMap(t *testing.T) {
	tests := []struct {
		name      string
		key       any
		value     any
		bindings  binding.Bindings
		want      any
		wantFound bool
		wantErr   bool
	}{{
		name: "map index not found",
		key:  "foo",
		value: map[string]any{
			"bar": 42,
		},
		bindings:  nil,
		want:      nil,
		wantFound: false,
		wantErr:   false,
	}, {
		name: "map index found",
		key:  "bar",
		value: map[string]any{
			"bar": 42,
		},
		bindings:  nil,
		want:      42,
		wantFound: true,
		wantErr:   false,
	}, {
		name: "map index found (and nil)",
		key:  "bar",
		value: map[string]any{
			"bar": nil,
		},
		bindings:  nil,
		want:      nil,
		wantFound: true,
		wantErr:   false,
	}, {
		name: "non string key (not found)",
		key:  3,
		value: map[int]any{
			2: "foo",
		},
		bindings:  nil,
		want:      nil,
		wantFound: false,
		wantErr:   false,
	}, {
		name: "non string key (found)",
		key:  2,
		value: map[int]any{
			2: "foo",
		},
		bindings:  nil,
		want:      "foo",
		wantFound: true,
		wantErr:   false,
	}, {
		name: "non string key (found and nil)",
		key:  2,
		value: map[int]any{
			2: nil,
		},
		bindings:  nil,
		want:      nil,
		wantFound: true,
		wantErr:   false,
	}, {
		name:      "nil value",
		key:       "foo",
		value:     nil,
		bindings:  nil,
		want:      nil,
		wantFound: false,
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := compilers.DefaultCompilers
			proj, err := ParseMapKey(nil, tt.key, compiler)
			assert.Nil(t, err)
			{
				got, found, err := proj.Handler(tt.value, tt.bindings)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tt.wantFound, found)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
