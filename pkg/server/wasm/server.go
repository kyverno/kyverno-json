//go:build js && wasm

package wasm

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server/api"
	"github.com/kyverno/kyverno-json/pkg/server/playground"
	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

const (
	apiPrefix        = "/api"
	playgroundPrefix = "/playground"
)

type Shutdown = func(context.Context) error

type Server interface {
	AddApiRoutes(api.Configuration) error
	AddPlaygroundRoutes() error
	Run(context.Context) Shutdown
}

type server struct {
	*gin.Engine
}

func New(enableLogger bool, enableCors bool) (Server, error) {
	router := gin.New()
	if enableLogger {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	if enableCors {
		router.Use(cors.New(cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"POST", "GET", "HEAD"},
			AllowHeaders:  []string{"Origin", "Content-Type"},
			ExposeHeaders: []string{"Content-Length"},
		}))
	}
	return server{router}, nil
}

func (s server) AddApiRoutes(config api.Configuration) error {
	return api.AddRoutes(s.Group(apiPrefix), config)
}

func (s server) AddPlaygroundRoutes() error {
	return playground.AddRoutes(s.Group(playgroundPrefix))
}

func (s server) Run(_ context.Context) Shutdown {
	wasmhttp.Serve(s.Engine.Handler())
	return nil
}
