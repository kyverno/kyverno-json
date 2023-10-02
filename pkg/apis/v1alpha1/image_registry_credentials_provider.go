package v1alpha1

// ImageRegistryCredentialsProvider provides the list of credential providers required.
// +kubebuilder:validation:Enum=default;amazon;azure;google;github
type ImageRegistryCredentialsProvider string

const (
	DEFAULT ImageRegistryCredentialsProvider = "default"
	AWS     ImageRegistryCredentialsProvider = "amazon"
	ACR     ImageRegistryCredentialsProvider = "azure"
	GCP     ImageRegistryCredentialsProvider = "google"
	GHCR    ImageRegistryCredentialsProvider = "github"
)
