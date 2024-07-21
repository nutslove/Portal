package routers

import (
	"net/http"
	"portal/middlewares"
	"portal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	_ = models.ConnectDB()
	store := cookie.NewStore([]byte("secret"))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("mysession", store))
	router.Use(middlewares.TracerSetting("Portal"))

	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"IsLoggedIn": username != nil,
			"Username":   username,
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tmpl", gin.H{
			"title": "Login",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		session.Set("username", username)
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tmpl", gin.H{
			"title": "SignUp",
		})
	})

	router.POST("/signup", func(c *gin.Context) {
		// signup処理を実装する
		c.Redirect(http.StatusFound, "/signup_success")
	})

	router.GET("/signup_success", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tmpl", gin.H{
			"title": "SignUp Success!",
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("username")
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	// user := router.Group("/users")

	// {
	// 	user.POST("/", func(c *gin.Context))
	// }
}
