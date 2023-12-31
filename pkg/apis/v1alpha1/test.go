package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope="Cluster"

// Test declares a test
type Tests struct {
	metav1.TypeMeta   `json:",inline,omitempty" yaml:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// Tests are the scenarios to be checked in the test
	Tests []Test `json:"tests,omitempty" yaml:"tests,omitempty"`
}

// Test declares a test result  fields
type Test struct {
	// Policies mentions the path to kyverno-json policies.
	Policies []string `json:"policies" yaml:"policies"`

	// Policy mentions the name of the policy.
	Payload string `json:"payload" yaml:"payload"`

	// Policy mentions the name of the policy.
	Preprocessors []string `json:"preprocess" yaml:"preprocess"`

	// Selectors mentions the labels selectors for policies.
	Selectors []string `json:"labels" yaml:"labels"`

	// Result mentions the result that the user is expecting.
	// Possible values are pass, fail and skip.
	Result PolicyResult `json:"result" yaml:"result"`
}

// +kubebuilder:validation:Enum=pass;fail;warn;error;skip
// PolicyResult specifies state of a policy result
type PolicyResult string

const (
	StatusPass  PolicyResult = "pass"
	StatusFail  PolicyResult = "fail"
	StatusWarn  PolicyResult = "warn"
	StatusError PolicyResult = "error"
	StatusSkip  PolicyResult = "skip"
)
