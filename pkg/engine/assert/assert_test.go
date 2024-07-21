package assert

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	tassert "github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestAssert(t *testing.T) {
	type args struct {
		assertion Assertion
		value     any
		bindings  binding.Bindings
	}
	tests := []struct {
		name    string
		args    args
		want    field.ErrorList
		wantErr bool
	}{{
		name: "nil vs empty object",
		args: args{
			assertion: Parse(context.TODO(), map[string]any{
				"foo": map[string]any{},
			}),
			value: map[string]any{
				"foo": nil,
			},
		},
		want: field.ErrorList{
			&field.Error{
				Type:   field.ErrorTypeInvalid,
				Field:  "foo",
				Detail: "invalid value, must not be null",
			},
		},
		wantErr: false,
	}, {
		name: "not nil vs empty object",
		args: args{
			assertion: Parse(context.TODO(), map[string]any{
				"foo": map[string]any{},
			}),
			value: map[string]any{
				"foo": map[string]any{
					"bar": 42,
				},
			},
		},
		want:    nil,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Assert(context.TODO(), nil, tt.args.assertion, tt.args.value, tt.args.bindings)
			if tt.wantErr {
				tassert.Error(t, err)
			} else {
				tassert.NoError(t, err)
			}
			tassert.Equal(t, tt.want, got)
		})
	}
}
