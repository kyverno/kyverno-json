package v1alpha1

// ContextEntry adds variables and data sources to a rule Context.
type ContextEntry struct {
	// Name is the variable name.
	Name string `json:"name"`

	// Variable defines an arbitrary JMESPath context variable that can be defined inline.
	Variable *Variable `json:"variable,omitempty"`

	// ImageRegistry defines requests to an OCI/Docker V2 registry to fetch image details.
	ImageRegistry *ImageRegistry `json:"imageRegistry,omitempty"`
}
