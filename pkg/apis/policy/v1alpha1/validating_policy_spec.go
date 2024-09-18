package v1alpha1

// ValidatingPolicySpec contains the policy spec.
type ValidatingPolicySpec struct {
	// Engine defines the default engine to use when evaluating expressions.
	Engine *Engine `json:"engine,omitempty"`

	// Rules is a list of ValidatingRule instances.
	Rules []ValidatingRule `json:"rules"`
}
