package scan

// import (
// 	"testing/fstest"

// 	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
// 	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
// 	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/exception"
// 	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
// 	"k8s.io/api/admissionregistration/v1alpha1"
// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/client-go/openapi"
// 	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
// 	"sigs.k8s.io/yaml"

// 	"github.com/kyverno/playground/backend/data"
// 	"github.com/kyverno/playground/backend/pkg/cluster"
// 	"github.com/kyverno/playground/backend/pkg/engine/models"
// 	"github.com/kyverno/playground/backend/pkg/policy"
// 	"github.com/kyverno/playground/backend/pkg/resource"
// )

type Request struct {
	Payload string `json:"payload"`
}

// func (r *EngineRequest) LoadParameters() (*models.Parameters, error) {
// 	var params models.Parameters
// 	if err := yaml.Unmarshal([]byte(r.Context), &params); err != nil {
// 		return nil, err
// 	}
// 	return &params, nil
// }

// func (r *EngineRequest) LoadPolicies(policyLoader loader.Loader) ([]kyvernov1.PolicyInterface, []v1alpha1.ValidatingAdmissionPolicy, error) {
// 	return policy.Load(policyLoader, []byte(r.Policies))
// }

// func (r *EngineRequest) LoadResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
// 	return resource.LoadResources(resourceLoader, []byte(r.Resources))
// }

// func (r *EngineRequest) LoadClusterResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
// 	return resource.LoadResources(resourceLoader, []byte(r.ClusterResources))
// }

// func (r *EngineRequest) LoadOldResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
// 	return resource.LoadResources(resourceLoader, []byte(r.OldResources))
// }

// func (r *EngineRequest) LoadPolicyExceptions(resourceLoader loader.Loader) ([]*kyvernov2alpha1.PolicyException, error) {
// 	return exception.Load([]byte(r.PolicyExceptions))
// }

// func (r *EngineRequest) LoadConfig(resourceLoader loader.Loader) (*corev1.ConfigMap, error) {
// 	if len(r.Config) == 0 {
// 		return nil, nil
// 	}
// 	return resource.Load[corev1.ConfigMap](resourceLoader, []byte(r.Config))
// }

// func (r *EngineRequest) ResourceLoader(cluster cluster.Cluster, kubeVersion string, config APIConfiguration) (loader.Loader, error) {
// 	var clients []openapi.Client
// 	if cluster != nil && !cluster.IsFake() {
// 		dclient, err := cluster.DClient()
// 		if err != nil {
// 			return nil, err
// 		}
// 		clients = append(clients, dclient.GetKubeClient().Discovery().OpenAPIV3())
// 	} else {
// 		kubeVersion, err := parseKubeVersion(kubeVersion)
// 		if err != nil {
// 			return nil, err
// 		}
// 		clients = append(clients, openapiclient.NewHardcodedBuiltins(kubeVersion))
// 	}
// 	clients = append(clients, openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas"))
// 	if len(r.CustomResourceDefinitions) != 0 {
// 		mapFs := fstest.MapFS{
// 			"crds.yaml": &fstest.MapFile{
// 				Data: []byte(r.CustomResourceDefinitions),
// 			},
// 		}
// 		clients = append(clients, openapiclient.NewLocalCRDFiles(mapFs, "."))
// 	}
// 	for _, crd := range config.LocalCrds {
// 		clients = append(clients, openapiclient.NewLocalCRDFiles(nil, crd))
// 	}
// 	for _, crd := range config.BuiltInCrds {
// 		fs, path := data.BuiltInCrds(crd)
// 		clients = append(clients, openapiclient.NewLocalCRDFiles(fs, path))
// 	}
// 	return loader.New(openapiclient.NewComposite(clients...))
// }
