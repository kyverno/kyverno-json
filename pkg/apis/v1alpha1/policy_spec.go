package v1alpha1

type PolicySpec struct {
	// Rules is a list of Rule instances. A Policy contains multiple rules and each rule can validate, mutate, or generate resources.
	Rules []Rule `json:"rules,omitempty"`
}
