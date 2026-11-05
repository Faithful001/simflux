package ui

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed dashboard.html
var dashboardHTML []byte

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ui", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", dashboardHTML)
	})
}
