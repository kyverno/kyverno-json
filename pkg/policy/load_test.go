package policy

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoad(t *testing.T) {
	basePath := "../../test/policy"
	tests := []struct {
		name    string
		path    string
		want    []*v1alpha1.ValidatingPolicy
		wantErr bool
	}{{
		name:    "confimap",
		path:    filepath.Join(basePath, "configmap.yaml"),
		wantErr: true,
	}, {
		name:    "not found",
		path:    filepath.Join(basePath, "not-found.yaml"),
		wantErr: true,
	}, {
		name:    "empty",
		path:    filepath.Join(basePath, "empty.yaml"),
		wantErr: false,
	}, {
		name:    "no spec",
		path:    filepath.Join(basePath, "no-spec.yaml"),
		wantErr: true,
	}, {
		name:    "no rules",
		path:    filepath.Join(basePath, "no-rules.yaml"),
		wantErr: true,
	}, {
		name:    "invalid rule",
		path:    filepath.Join(basePath, "bad-rule.yaml"),
		wantErr: true,
	}, {
		name:    "rule name missing",
		path:    filepath.Join(basePath, "rule-name-missing.yaml"),
		wantErr: true,
	}, {
		name: "ok",
		path: filepath.Join(basePath, "ok.yaml"),
		want: []*v1alpha1.ValidatingPolicy{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "json.kyverno.io/v1alpha1",
				Kind:       "ValidatingPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ValidatingPolicySpec{
				Rules: []v1alpha1.ValidatingRule{{
					Name: "pod-no-latest",
					Match: &v1alpha1.Match{
						Any: []v1alpha1.Any{{
							Value: map[string]any{
								"apiVersion": "v1",
								"kind":       "Pod",
							},
						}},
					},
					Assert: &v1alpha1.Assert{
						All: []v1alpha1.Assertion{{
							Check: v1alpha1.Any{
								Value: map[string]any{
									"spec": map[string]any{
										"~foo.containers->foos": map[string]any{
											"(at($foos, $foo).image)->foo": map[string]any{
												"(contains($foo, ':'))":        true,
												"(ends_with($foo, ':latest'))": false,
											},
										},
									},
								},
							},
						}},
					},
				}},
			},
		}},
	}, {
		name: "multiple",
		path: filepath.Join(basePath, "multiple.yaml"),
		want: []*v1alpha1.ValidatingPolicy{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "json.kyverno.io/v1alpha1",
				Kind:       "ValidatingPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Spec: v1alpha1.ValidatingPolicySpec{
				Rules: []v1alpha1.ValidatingRule{},
			},
		}, {
			TypeMeta: metav1.TypeMeta{
				APIVersion: "json.kyverno.io/v1alpha1",
				Kind:       "ValidatingPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Spec: v1alpha1.ValidatingPolicySpec{
				Rules: []v1alpha1.ValidatingRule{},
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
