package v1alpha1

// ValidatingRule defines a validating rule.
type ValidatingRule struct {
	// Name is a label to identify the rule, It must be unique within the policy.
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	// Context defines variables and data sources that can be used during rule execution.
	// +optional
	Context []ContextEntry `json:"context,omitempty"`

	// Match defines when this policy rule should be applied.
	// +optional
	Match *Match `json:"match,omitempty"`

	// Exclude defines when this policy rule should not be applied.
	// +optional
	Exclude *Match `json:"exclude,omitempty"`

	// Assert is used to validate matching resources.
	Assert *Assert `json:"assert"`
}
