//go:build js && wasm

package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server"
	"github.com/kyverno/kyverno-json/pkg/server/playground"
)

func main() {
	// initialise gin framework
	gin.SetMode(gin.DebugMode)
	// create server
	router, err := server.New(true, true)
	if err != nil {
		panic(err)
	}
	// register playground routes
	if err := playground.AddRoutes(router.Group(server.PlaygroundPrefix)); err != nil {
		panic(err)
	}
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	server.RunWasm(ctx, router)
	<-ctx.Done()
	stop()
}
