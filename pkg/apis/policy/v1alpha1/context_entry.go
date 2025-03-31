package v1alpha1

// ContextEntry adds variables and data sources to a rule context.
type ContextEntry struct {
	// Compiler defines the default compiler to use when evaluating expressions.
	// +optional
	Compiler *Compiler `json:"compiler,omitempty"`

	// Name is the entry name.
	Name string `json:"name"`

	// Variable defines an arbitrary variable.
	// +optional
	Variable Any `json:"variable,omitempty"`

	// ConfigMap allows referencing a ConfigMap as data source.
	// +optional
	ConfigMap *ConfigMapReference `json:"configMap,omitempty"`
}

// ConfigMapReference references a ConfigMap resource.
type ConfigMapReference struct {
	// Name is the ConfigMap name.
	Name string `json:"name"`

	// Namespace is the ConfigMap namespace.
	Namespace string `json:"namespace"`
}
