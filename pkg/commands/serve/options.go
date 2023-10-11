package serve

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server"
	"github.com/kyverno/kyverno-json/pkg/server/api"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/spf13/cobra"
)

type options struct {
	serverFlags serverFlags
	ginFlags    ginFlags
	// engineFlags engineFlags
}

type serverFlags struct {
	host string
	port int
}

type ginFlags struct {
	mode        string
	log         bool
	cors        bool
	maxBodySize int
}

// type uiFlags struct {
// 	sponsor string
// }

// type engineFlags struct {
// 	builtInCrds []string
// 	localCrds   []string
// }

// type clusterFlags struct {
// 	cluster             bool
// 	kubeConfigOverrides clientcmd.ConfigOverrides
// }

func (c *options) Run(_ *cobra.Command, _ []string) error {
	// initialise gin framework
	gin.SetMode(c.ginFlags.mode)
	tonic.SetBindHook(tonic.DefaultBindingHookMaxBodyBytes(int64(c.ginFlags.maxBodySize)))
	// tonic.SetErrorHook(func(c *gin.Context, err error) (int, interface{}) {
	// 	switch e := err.(type) {
	// 	case engine.PolicyViolationError:
	// 		return http.StatusBadRequest, gin.H{
	// 			"violations": e.Violations,
	// 			"error":      e.Error(),
	// 			"reason":     "POLICY_VALIDATION",
	// 		}
	// 	default:
	// 		return http.StatusBadRequest, gin.H{
	// 			"error":  e.Error(),
	// 			"reason": "ERROR",
	// 		}
	// 	}
	// })
	// create server
	server, err := server.New(c.ginFlags.log, c.ginFlags.cors)
	if err != nil {
		return err
	}
	apiConfig := api.APIConfiguration{
		EngineConfiguration: api.EngineConfiguration{
			// BuiltInCrds: c.engineFlags.builtInCrds,
			// LocalCrds:   c.engineFlags.localCrds,
		},
	}
	// register API routes (with/without cluster support)
	// if c.clusterFlags.cluster {
	// 	// create rest config
	// 	restConfig, err := utils.RestConfig(c.clusterFlags.kubeConfigOverrides)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// create cluster
	// 	cluster, err := cluster.New(restConfig)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// register API routes
	// 	if err := server.AddAPIRoutes(cluster, apiConfig); err != nil {
	// 		return err
	// 	}
	// } else {
	// register API routes
	if err := server.AddAPIRoutes(apiConfig); err != nil {
		return err
	}
	// }
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdown := server.Run(ctx, c.serverFlags.host, c.serverFlags.port)
	<-ctx.Done()
	stop()
	if shutdown != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
