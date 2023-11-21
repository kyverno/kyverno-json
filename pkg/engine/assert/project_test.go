package assert

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	tassert "github.com/stretchr/testify/assert"
)

func Test_project(t *testing.T) {
	tests := []struct {
		name     string
		key      interface{}
		value    interface{}
		bindings binding.Bindings
		want     *projection
		wantErr  bool
	}{{
		name: "map index not found",
		key:  "foo",
		value: map[string]interface{}{
			"bar": 42,
		},
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project(context.TODO(), tt.key, tt.value, tt.bindings)
			if tt.wantErr {
				tassert.Error(t, err)
				tassert.Equal(t, tt.want, got)
			} else {
				tassert.NoError(t, err)
			}
		})
	}
}
