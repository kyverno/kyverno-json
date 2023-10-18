package serve

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/client/clientset/versioned"
	"github.com/kyverno/kyverno-json/pkg/server"
	"github.com/kyverno/kyverno-json/pkg/server/api"
	restutils "github.com/kyverno/kyverno-json/pkg/utils/rest"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type options struct {
	serverFlags  serverFlags
	ginFlags     ginFlags
	clusterFlags clusterFlags
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

type clusterFlags struct {
	kubeConfigOverrides clientcmd.ConfigOverrides
}

func (c *options) Run(_ *cobra.Command, _ []string) error {
	// initialise gin framework
	gin.SetMode(c.ginFlags.mode)
	tonic.SetBindHook(tonic.DefaultBindingHookMaxBodyBytes(int64(c.ginFlags.maxBodySize)))
	// create server
	server, err := server.New(c.ginFlags.log, c.ginFlags.cors)
	if err != nil {
		return err
	}
	restConfig, err := restutils.RestConfig(c.clusterFlags.kubeConfigOverrides)
	if err != nil {
		return err
	}
	client, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	config := api.Configuration{
		PolicyProvider: &provider{
			client: client,
		},
	}
	// register api routes
	if err := server.AddApiRoutes(config); err != nil {
		return err
	}
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
