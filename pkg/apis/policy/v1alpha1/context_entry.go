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

	// APICall defines an HTTP or Kubernetes API call.
	// +optional
	APICall *APICall `json:"apiCall,omitempty"`
}

// APICall defines an HTTP or Kubernetes API call.
type APICall struct {
	// URLPath is the URL path to be used in the HTTP GET or POST request to the
	// Kubernetes API server (e.g. "/api/v1/namespaces" or  "/apis/apps/v1/deployments").
	// +optional
	URLPath Any `json:"urlPath,omitempty"`

	// Method is the HTTP request type (GET or POST). Defaults to GET.
	// +kubebuilder:default=GET
	// +optional
	Method string `json:"method,omitempty"`

	// The data object specifies the POST data sent to the server.
	// Only applicable when the method field is set to POST.
	// +optional
	Data []RequestData `json:"data,omitempty"`

	// Service is an API call to a JSON web service.
	// This is used for non-Kubernetes API server calls.
	// It's mutually exclusive with the URLPath field.
	// +optional
	Service *ServiceCall `json:"service,omitempty"`

	// JMESPath is an optional JSON Match Expression that can be used to
	// transform the JSON response returned from the server.
	// +optional
	JMESPath string `json:"jmesPath,omitempty"`

	// Default is an optional arbitrary JSON object that the context
	// value is set to, if the apiCall returns an error.
	// +optional
	Default *Any `json:"default,omitempty"`
}

// ServiceCall specifies a JSON web service.
type ServiceCall struct {
	// URL is the JSON web service URL.
	URL Any `json:"url"`
}

// RequestData specifies POST data sent to the server.
type RequestData struct {
	// Key is the form key.
	Key string `json:"key"`

	// Value is the form value.
	Value Any `json:"value"`
}
