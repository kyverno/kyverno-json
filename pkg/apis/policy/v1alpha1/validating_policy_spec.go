package v1alpha1

// ValidatingPolicySpec contains the policy spec.
type ValidatingPolicySpec struct {
	// Rules is a list of ValidatingRule instances.
	Rules []ValidatingRule `json:"rules"`
}
