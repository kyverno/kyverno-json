package scan

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup, policyProvider PolicyProvider) error {
	handler, err := newHandler(policyProvider)
	if err != nil {
		return err
	}
	group.POST("/scan", handler)
	log.Default().Printf("configured route %s/%s", group.BasePath(), "scan")
	return nil
}
