package playground

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server"
	"github.com/kyverno/kyverno-json/pkg/server/ui"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/spf13/cobra"
)

type options struct {
	serverFlags serverFlags
	ginFlags    ginFlags
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

func (c *options) Run(_ *cobra.Command, _ []string) error {
	// initialise gin framework
	gin.SetMode(c.ginFlags.mode)
	tonic.SetBindHook(tonic.DefaultBindingHookMaxBodyBytes(int64(c.ginFlags.maxBodySize)))
	// create router
	router, err := server.New(c.ginFlags.log, c.ginFlags.cors)
	if err != nil {
		return err
	}
	// register api routes
	if err := ui.AddRoutes(router); err != nil {
		return err
	}
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdown := server.Run(ctx, router, c.serverFlags.host, c.serverFlags.port)
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
