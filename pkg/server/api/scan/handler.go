package scan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
)

// import (
// 	"context"
// 	"fmt"
// 	"net/http"

// 	"github.com/Masterminds/semver/v3"
// 	"github.com/gin-gonic/gin"
// 	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
// 	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
// 	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
// 	"github.com/loopfz/gadgeto/tonic"
// 	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

// 	"github.com/kyverno/playground/backend/data"
// 	"github.com/kyverno/playground/backend/pkg/cluster"
// 	"github.com/kyverno/playground/backend/pkg/engine"
// 	"github.com/kyverno/playground/backend/pkg/engine/models"
// )

func newHandler(config APIConfiguration) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *Request) (*EngineResponse, error) {
		return &EngineResponse{
			// Policies:  policies,
			// Resources: resources,
			// Results:   results,
		}, nil
	}, http.StatusOK), nil
}

// func parseKubeVersion(kubeVersion string) (string, error) {
// 	if kubeVersion == "" {
// 		return "1.28", nil
// 	}
// 	version, err := semver.NewVersion(kubeVersion)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprint(version.Major(), ".", version.Minor()), nil
// }

// func validateParams(params *models.Parameters, cmResolver engineapi.ConfigmapResolver, policies []kyvernov1.PolicyInterface) error {
// 	if params == nil {
// 		return nil
// 	}

// 	for _, policy := range policies {
// 		for _, rule := range policy.GetSpec().Rules {
// 			for _, variable := range rule.Context {
// 				if variable.APICall == nil && variable.ConfigMap == nil {
// 					continue
// 				}
// 				if _, ok := params.Variables[variable.Name]; ok {
// 					continue
// 				}
// 				if variable.ConfigMap != nil {
// 					_, err := cmResolver.Get(context.Background(), variable.ConfigMap.Namespace, variable.ConfigMap.Name)
// 					if err == nil {
// 						continue
// 					}
// 				}

// 				return fmt.Errorf("Variable %s is not defined in the context", variable.Name)
// 			}
// 		}
// 	}

// 	return nil
// }
