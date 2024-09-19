package projection

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	tassert "github.com/stretchr/testify/assert"
)

func TestProjection(t *testing.T) {
	tests := []struct {
		name     string
		key      any
		value    any
		bindings binding.Bindings
		want     any
		wantErr  bool
	}{{
		name: "map index not found",
		key:  "foo",
		value: map[string]any{
			"bar": 42,
		},
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}, {
		name: "map index found",
		key:  "bar",
		value: map[string]any{
			"bar": 42,
		},
		bindings: nil,
		want:     42,
		wantErr:  false,
	}, {
		name: "map index found (and nil)",
		key:  "bar",
		value: map[string]any{
			"bar": nil,
		},
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}, {
		name: "non string key (not found)",
		key:  3,
		value: map[int]any{
			2: "foo",
		},
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}, {
		name: "non string key (found)",
		key:  2,
		value: map[int]any{
			2: "foo",
		},
		bindings: nil,
		want:     "foo",
		wantErr:  false,
	}, {
		name: "non string key (found and nil)",
		key:  2,
		value: map[int]any{
			2: nil,
		},
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}, {
		name:     "nil value",
		key:      "foo",
		value:    nil,
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj := Parse(tt.key)
			got, err := proj.Handler(context.TODO(), tt.value, tt.bindings)
			if tt.wantErr {
				tassert.Error(t, err)
			} else {
				tassert.NoError(t, err)
			}
			tassert.Equal(t, tt.want, got)
		})
	}
}
