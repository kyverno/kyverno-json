package playground

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup) error {
	handler, err := newHandler()
	if err != nil {
		return err
	}
	group.POST("/scan", handler)
	return nil
}
