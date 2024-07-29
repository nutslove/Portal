package routers

import (
	"net/http"
	"os"
	"portal/controllers"
	"portal/middlewares"
	"portal/models"
	"portal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	// 404ハンドラーを設定
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	db := models.ConnectDB()
	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SECRET_KEY")))
	store.Options(sessions.Options{
		MaxAge: 3600, // セッションの有効期限(秒)
		Path:   "/",  // クッキーが有効なパス("/"は全URLパス対象)
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("usersession", store))
	router.Use(middlewares.TracerSetting("Portal"))

	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(http.StatusOK, "index.tpl", gin.H{
			"IsLoggedIn": username != nil,
			"Username":   username,
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tpl", gin.H{
			"title": "Login",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		plainPassword := c.PostForm("password")
		userExistCheckResult := services.UserExistCheck(username, db)
		hashedPassword := services.GetUserPassword(username, db)
		passwordUnmatchErr := controllers.ComparePassword(hashedPassword, plainPassword)
		if userExistCheckResult && passwordUnmatchErr == nil {
			session.Set("username", username)
			session.Save()
			c.Redirect(http.StatusFound, "/")
		} else {
			c.HTML(http.StatusUnauthorized, "signinup.tpl", gin.H{
				"title":  "Login",
				"status": "loginfailed",
			})
		}
	})

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tpl", gin.H{
			"title": "SignUp",
		})
	})

	router.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		if services.UserExistCheck(c.PostForm("username"), db) {
			c.HTML(http.StatusBadRequest, "signinup.tpl", gin.H{
				"userid": username,
				"title":  "SignUp",
				"status": "signupfailed",
				"reason": "existalready",
			})
			return
		}
		services.UserCreate(c, db)
		c.Redirect(http.StatusFound, "/signup_success")
	})

	router.GET("/signup_success", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signinup.tpl", gin.H{
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
