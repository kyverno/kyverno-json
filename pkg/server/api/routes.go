package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server/api/scan"
)

func AddRoutes(group *gin.RouterGroup, config Configuration) error {
	if err := scan.AddRoutes(group, config.PolicyProvider); err != nil {
		return err
	}
	return nil
}
