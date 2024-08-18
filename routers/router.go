package routers

import (
	"fmt"
	"net/http"
	"os"
	"portal/controllers"
	"portal/middlewares"
	"portal/models"
	"portal/services"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// ユーザからのPostの投稿内容
type RequestData struct {
	Content string `json:"content"`
}

func SetupRouter(router *gin.Engine) {

	// 404ハンドラーを設定
	router.NoRoute(func(c *gin.Context) {
		controllers.NotFoundResponse(c)
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

	career := router.Group("/career")
	{
		career.GET("", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			posts, _, pageNum, pageSlice := services.GetCareerPostsList(1, db)
			// 最初はpost数が0のケースもあり得るので、トップページではpost数による判定はせずにTableのヘッダーだけ表示する
			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"IsLoggedIn": username != nil,
				"Username":   username,
				"posts":      posts,
				"page":       1,
				"pageSlice":  pageSlice,
				"pageNum":    pageNum,
				"PostList":   true,
				"BoardType":  "career",
				"BoardName":  "キャリア相談",
			})
		})

		career.GET("/:page", func(c *gin.Context) {

			// pageに文字列(string)が入っていたり、存在しないページを直接指定した場合404ページを返す
			session := sessions.Default(c)
			username := session.Get("username")
			postId := c.Query("post")
			if postId != "" {
				PostContent, err := services.GetCareerPostContent(postId)
				if err != nil {
					controllers.NotFoundResponse(c)
					return
				}
				// fmt.Println("post番号:", post)
				c.HTML(http.StatusOK, "index.tpl", gin.H{
					"IsLoggedIn":  username != nil,
					"Username":    username,
					"PostRead":    true,
					"PostContent": PostContent,
				})
				return
			}
			page := c.Param("page")
			pageInt, err := strconv.Atoi(page)
			if err != nil {
				controllers.NotFoundResponse(c)
				return
			}
			posts, postNum, pageNum, pageSlice := services.GetCareerPostsList(pageInt, db)
			if postNum == 0 {
				controllers.NotFoundResponse(c)
				return
			}

			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"IsLoggedIn": username != nil,
				"Username":   username,
				"posts":      posts,
				"page":       pageInt,
				"pageSlice":  pageSlice,
				"pageNum":    pageNum,
				"PostList":   true,
				"BoardType":  "career",
				"BoardName":  "キャリア相談",
			})
		})

		career.GET("/posting", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"IsLoggedIn": username != nil,
				"Username":   username,
				"PostWrite":  true,
				"BoardType":  "career",
			})
		})

		career.POST("/posting", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")

			/////// 以下の処理は後でserviceとかに全部移す
			var requestData RequestData

			// リクエストボディをJSONとしてバインド
			if err := c.ShouldBindJSON(&requestData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid request data",
				})
				return
			}
			fmt.Println("requestData:", requestData)

			addedPost := models.CareerBoard{
				// Numberは主キーでautoIncrementオプションが有効になっていて指定しなくても自動で最後のレコードのNumber＋１でInsertしてくれる
				Title:  "テスト", // ポスト入力ページにTitleの入力欄を追加し、入力必須とする
				Author: username.(string),
				Date:   time.Now(),
				Count:  0,
			}

			result := db.Create(&addedPost)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": result.Error,
				})
				return
			}

			// fmt.Println("挿入されたレコードのNumber:", addedPost.Number)

			///// OpenSearchにデータを入れる処理追加

			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"redirectUrl": fmt.Sprintf("/career/1?post=%d", addedPost.Number),
			})
		})
	}

	// user := router.Group("/users")

	// {
	// 	user.POST("/", func(c *gin.Context))
	// }
}
