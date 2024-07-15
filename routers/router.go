package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"": "",
		})
	})

	// user := router.Group("/users")

	// {
	// 	user.POST("/", func(c *gin.Context))
	// }
}
