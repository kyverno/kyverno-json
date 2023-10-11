package scan

import (
	"github.com/gin-gonic/gin"
)

type APIConfiguration struct {
	BuiltInCrds []string
	LocalCrds   []string
}

func AddRoutes(group *gin.RouterGroup, config APIConfiguration) error {
	handler, err := newHandler(config)
	if err != nil {
		return err
	}
	group.POST("/scan", handler)
	return nil
}
