package v1alpha1

type ImageRegistryCredentials struct {
	// AllowInsecureRegistry allows insecure access to a registry.
	AllowInsecureRegistry bool `json:"allowInsecureRegistry,omitempty"`

	// Providers specifies a list of OCI Registry names, whose authentication providers are provided.
	// It can be of one of these values: AWS, ACR, GCP, GHCR.
	Providers []ImageRegistryCredentialsProvider `json:"providers,omitempty"`

	// Secrets specifies a list of secrets that are provided for credentials.
	// Secrets must live in the Kyverno namespace.
	Secrets []string `json:"secrets,omitempty"`
}
