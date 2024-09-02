package routers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/config"
	"portal/controllers"
	"portal/middlewares"
	"portal/models"
	"portal/services"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

// ユーザからのPostの投稿内容
type RequestData struct {
	Title   string `json:"title"`
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
			"TopPage":    true,
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

		// Post一覧読み込み。
		career.GET("/:page", func(c *gin.Context) {
			// pageに文字列(string)が入っていたり、存在しないページを直接指定した場合404ページを返す
			session := sessions.Default(c)
			username := session.Get("username")

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

		// 新規Post書き込み画面
		career.GET("/posting", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			// ログインせず直でURL指定してアクセスした場合、１ページにリダイレクトする
			if username == nil {
				c.Redirect(http.StatusFound, "/career/")
			}

			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"IsLoggedIn": username != nil,
				"Username":   username,
				"PostWrite":  true,
				"BoardType":  "career",
			})
		})

		// 既存Post読み込み/修正画面
		career.GET("/posting/:postId", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")

			postId := c.Param("postId")
			PostContent, PostTitle, Author, err := services.GetCareerPostContent(postId)
			if err != nil {
				controllers.NotFoundResponse(c)
				return
			}

			postIdInt, err := strconv.Atoi(postId)
			if err != nil {
				log.Fatal("cannot convert postId(string) to int:", err)
				return
			}
			postdate := services.GetCareerPostDate(postIdInt, db)

			// URLパラメータに"modify=true"が設定されている場合は既存Postの修正
			if is_modify := c.Query("modify"); is_modify == "true" {
				if username == nil || username != Author {
					c.Redirect(http.StatusFound, "/career/")
				}

				c.HTML(http.StatusOK, "index.tpl", gin.H{
					"IsLoggedIn":  username != nil,
					"Username":    username,
					"PostWrite":   true,
					"BoardType":   "career",
					"PostContent": PostContent,
					"PostTitle":   PostTitle,
					"PostId":      postId,
					"Modify":      true,
				})
				return
			}

			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"IsLoggedIn":  username != nil,
				"Username":    username,
				"CreatedAt":   postdate["CreatedAt"],
				"ModifiedAt":  postdate["ModifiedAt"],
				"Author":      Author,
				"PostId":      postId,
				"PostRead":    true,
				"PostContent": PostContent,
				"PostTitle":   PostTitle,
				"BoardType":   "career",
			})
		})

		// 新規Post登録
		career.POST("/posting", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				controllers.UnAuthorizedResponse(c)
				return
			}

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
				Title:      requestData.Title,
				Author:     username.(string),
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
				// Count:      0,
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

			///// OpenSearchにデータを入れる処理
			type Post struct {
				Title  string `json:"title"`
				Author string `json:"author"`
				Post   string `json:"post"`
			}

			post := Post{
				Title:  requestData.Title,
				Author: username.(string),
				Post:   requestData.Content,
			}

			client, err := config.OpensearchNewClient()
			if err != nil {
				log.Fatal("cannot initialize", err)
			}

			// ドキュメントの挿入
			insertResp, err := client.Index(
				context.Background(),
				opensearchapi.IndexReq{
					Index:      "career",
					DocumentID: strconv.Itoa(addedPost.Number),
					Body:       opensearchutil.NewJSONReader(&post),
					Params: opensearchapi.IndexParams{
						Refresh: "true",
					},
				})

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": err,
				})
				return
			}
			fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"redirectUrl": fmt.Sprintf("/career/posting/%d", addedPost.Number),
			})
		})

		career.DELETE("/posting/:postId", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			if username == nil {
				controllers.UnAuthorizedResponse(c)
				return
			}
			postId := c.Param("postId")
			if postId == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"seccess": false,
					"message": "post number is missing!",
				})
				return
			}

			postIdInt, err := strconv.Atoi(postId)
			if err != nil {
				log.Fatal("cannot convert postId(string) to int:", err)
				return
			}
			if !(services.PostExistCheck(postIdInt, db)) {
				c.JSON(http.StatusBadRequest, gin.H{
					"seccess": false,
					"message": "該当Postが存在しません",
				})
				return
			}

			services.DeleteCareerPost(postIdInt, db)

			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"redirectUrl": "/career",
			})
		})

		// 既存Post内容修正
		career.PATCH("/posting/:postId", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")
			postId := c.Param("postId")
			if username == nil {
				controllers.UnAuthorizedResponse(c)
				return
			}

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

			postIdInt, err := strconv.Atoi(postId)
			if err != nil {
				log.Fatal("cannot convert postId(string) to int:", err)
				return
			}

			postForDB := models.CareerBoard{
				Number:     postIdInt,
				Title:      requestData.Title,
				ModifiedAt: time.Now(),
			}

			result := db.Updates(&postForDB)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": result.Error,
				})
				return
			}

			// fmt.Println("挿入されたレコードのNumber:", addedPost.Number)

			///// OpenSearchにデータを入れる処理
			type Post struct {
				Title  string `json:"title"`
				Author string `json:"author"`
				Post   string `json:"post"`
			}

			type UpdateRequest struct {
				Doc Post `json:"doc"`
			}

			postForOpensearch := Post{
				Title:  requestData.Title,
				Author: username.(string),
				Post:   requestData.Content,
			}

			updatedPost := UpdateRequest{
				Doc: postForOpensearch,
			}

			client, err := config.OpensearchNewClient()
			if err != nil {
				log.Fatal("cannot initialize", err)
			}

			// ドキュメントの更新
			updateResp, err := client.Update(
				context.Background(),
				opensearchapi.UpdateReq{
					Index:      "career",
					DocumentID: postId,
					Body:       opensearchutil.NewJSONReader(&updatedPost),
					Params: opensearchapi.UpdateParams{
						Refresh: "true",
					},
				})

			// fmt.Println("OpenSearch updateResp:", updateResp)
			if err != nil {
				fmt.Println("OpenSearch Update err:", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": err,
				})
				return
			}
			fmt.Printf("Updated document in %s\n  ID: %s\n", updateResp.Index, updateResp.ID)

			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"redirectUrl": fmt.Sprintf("/career/posting/%d", postForDB.Number),
			})
		})
	}

	// Post検索、検索にヒットしたPost一覧表示
	career.GET("/search", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		query := c.Query("query") // inputのnameの値を指定して取得

		posts := services.SearchCareerPost(query, db)

		c.HTML(http.StatusOK, "index.tpl", gin.H{
			"IsLoggedIn":        username != nil,
			"Username":          username,
			"posts":             posts,
			"page":              1,
			"pageSlice":         []int{1},
			"pageNum":           1,
			"PostList":          true,
			"BoardType":         "career",
			"BoardName":         "キャリア相談",
			"query":             query,
			"SearchResultCount": len(posts),
		})
	})

	// user := router.Group("/users")

	// {
	// 	user.POST("/", func(c *gin.Context))
	// }
}
