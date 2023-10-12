package ui

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) error {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(fs))
	router.NoRoute(func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
