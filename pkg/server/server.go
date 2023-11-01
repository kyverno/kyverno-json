package server

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	ApiPrefix        = "/api"
	PlaygroundPrefix = "/playground"
)

type Shutdown = func(context.Context) error

type Server = *gin.Engine

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
	return router, nil
}
