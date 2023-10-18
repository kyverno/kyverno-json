package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server/api"
	"github.com/kyverno/kyverno-json/pkg/server/playground"
)

const (
	apiPrefix        = "/api"
	playgroundPrefix = "/playground"
)

type Shutdown = func(context.Context) error

type Server interface {
	AddApiRoutes(api.Configuration) error
	AddPlaygroundRoutes() error
	Run(context.Context, string, int) Shutdown
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

func (s server) Run(_ context.Context, host string, port int) Shutdown {
	address := fmt.Sprintf("%v:%v", host, port)
	srv := &http.Server{
		Addr:              address,
		Handler:           s.Engine.Handler(),
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return srv.Shutdown
}
