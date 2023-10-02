package v1alpha1

// ImageRegistry defines requests to an OCI/Docker V2 registry to fetch image details.
type ImageRegistry struct {
	// Reference is image reference to a container image in the registry.
	// Example: ghcr.io/kyverno/kyverno:latest
	Reference string `json:"reference"`

	// ImageRegistryCredentials provides credentials that will be used for authentication with registry.
	ImageRegistryCredentials *ImageRegistryCredentials `json:"imageRegistryCredentials,omitempty"`
}
