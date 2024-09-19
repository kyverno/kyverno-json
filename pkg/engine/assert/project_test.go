package assert

// import (
// 	"context"
// 	"testing"

// 	"github.com/jmespath-community/go-jmespath/pkg/binding"
// 	"github.com/kyverno/kyverno-json/pkg/syntax/expression"
// 	tassert "github.com/stretchr/testify/assert"
// 	"k8s.io/utils/ptr"
// )

// func Test_project(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		key      any
// 		value    any
// 		bindings binding.Bindings
// 		want     *projection
// 		wantErr  bool
// 	}{{
// 		name: "map index not found",
// 		key:  "foo",
// 		value: map[string]any{
// 			"bar": 42,
// 		},
// 		bindings: nil,
// 		want:     nil,
// 		wantErr:  false,
// 	}, {
// 		name: "map index found",
// 		key:  "bar",
// 		value: map[string]any{
// 			"bar": 42,
// 		},
// 		bindings: nil,
// 		want: &projection{
// 			result: 42,
// 		},
// 		wantErr: false,
// 	}, {
// 		name: "map index found (and nil)",
// 		key:  "bar",
// 		value: map[string]any{
// 			"bar": nil,
// 		},
// 		bindings: nil,
// 		want: &projection{
// 			result: nil,
// 		},
// 		wantErr: false,
// 	}, {
// 		name: "non string key (not found)",
// 		key:  3,
// 		value: map[int]any{
// 			2: "foo",
// 		},
// 		bindings: nil,
// 		want:     nil,
// 		wantErr:  false,
// 	}, {
// 		name: "non string key (found)",
// 		key:  2,
// 		value: map[int]any{
// 			2: "foo",
// 		},
// 		bindings: nil,
// 		want: &projection{
// 			result: "foo",
// 		},
// 		wantErr: false,
// 	}, {
// 		name: "non string key (found and nil)",
// 		key:  2,
// 		value: map[int]any{
// 			2: nil,
// 		},
// 		bindings: nil,
// 		want: &projection{
// 			result: nil,
// 		},
// 		wantErr: false,
// 	}, {
// 		name:     "nil value",
// 		key:      "foo",
// 		value:    nil,
// 		bindings: nil,
// 		want:     nil,
// 		wantErr:  false,
// 	}}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var expr *expression.Expression
// 			switch typed := tt.key.(type) {
// 			case string:
// 				expr = ptr.To(expression.Parse(typed))
// 			}
// 			got, err := project(context.TODO(), expr, tt.key, tt.value, tt.bindings)
// 			if tt.wantErr {
// 				tassert.Error(t, err)
// 			} else {
// 				tassert.NoError(t, err)
// 			}
// 			tassert.Equal(t, tt.want, got)
// 		})
// 	}
// }
