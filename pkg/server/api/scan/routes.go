package scan

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup, policyProvider PolicyProvider) error {
	handler, err := newHandler(policyProvider)
	if err != nil {
		return err
	}
	group.POST("/scan", handler)
	return nil
}
