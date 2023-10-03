package v1alpha1

type Rule struct {
	// Name is a label to identify the rule, It must be unique within the policy.
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	// Context defines variables and data sources that can be used during rule execution.
	Context []ContextEntry `json:"context,omitempty"`

	// Match defines when this policy rule should be applied. The match
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the user name or role.
	// At least one kind is required.
	Match *Match `json:"match,omitempty"`

	// Exclude defines when this policy rule should not be applied. The exclude
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the name or role.
	Exclude *Match `json:"exclude,omitempty"`

	// Validation is used to validate matching resources.
	Validation *Validation `json:"validate,omitempty"`
}
