package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server/api/scan"
)

type EngineConfiguration = scan.APIConfiguration

type APIConfiguration struct {
	EngineConfiguration
	Sponsor string
}

func AddRoutes(group *gin.RouterGroup, config APIConfiguration) error {
	if err := scan.AddRoutes(group, config.EngineConfiguration); err != nil {
		return err
	}
	return nil
}
