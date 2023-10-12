package playground

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/server/playground/scan"
)

func AddRoutes(group *gin.RouterGroup) error {
	if err := scan.AddRoutes(group); err != nil {
		return err
	}
	return nil
}
