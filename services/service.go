package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"portal/config"
	"portal/controllers"
	"portal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const postNumPerPage int = 10

func UserExistCheck(userid string, db *gorm.DB) bool {
	user := models.UserData{Nickname: userid}
	selectResult := db.Take(&user)
	if errors.Is(selectResult.Error, gorm.ErrRecordNotFound) {
		return false
	} else {
		return true
	}
}

func GetUserPassword(userid string, db *gorm.DB) string {
	var user models.UserData
	db.Where("userid = ?", userid).Take(&user)
	return user.Password
}

func UserCreate(c *gin.Context, db *gorm.DB) {
	var age *int
	if c.PostForm("age") != "" {
		ageValue, err := strconv.Atoi(c.PostForm("age"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "signinup.tpl", gin.H{
				"title":  "SignUp",
				"status": "signupfailed",
				"reason": "agestring",
			})
			return
		}
		age = &ageValue
	}
	pw, err := controllers.EncryptPassword(c.PostForm("password"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signinup.tpl", gin.H{
			"title":  "SignUp",
			"status": "signupfailed",
			"reason": "internalservererror",
		})
		return
	}
	user := models.UserData{
		Nickname: c.PostForm("username"),
		Password: pw,
		Age:      age,
		Company:  c.PostForm("company"),
		Role:     c.PostForm("role"),
	}
	result := db.Create(&user)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "signinup.tpl", gin.H{
			"title":  "SignUp",
			"status": "signupfailed",
			"reason": "internalservererror",
		})
	}
}

// func UserDelete(c *gin.Context, db *gorm.DB) {

// }

func GetCareerPostsList(page int, db *gorm.DB) ([]map[string]interface{}, int, int, []int) {
	var posts []models.CareerBoard
	db.Order("num desc").Offset((page - 1) * postNumPerPage).Limit(postNumPerPage).Find(&posts) // 1ページ内にpostはpostNumPerPageに定義した個数まで表示

	postNum := len(posts)
	var count int64
	db.Model(&models.CareerBoard{}).Count(&count)
	pageNum := int(math.Ceil((float64(count) / float64(postNumPerPage)))) // 総ページ数 ( 1ページ内にpostはpostNumPerPageに定義した個数まで表示)
	if pageNum == 0 {
		pageNum = 1
	}
	var pageSlice []int
	if pageNum > 10 && page > pageNum-9 { // 総ページ数が10個より多くてユーザがアクセスしたページが 最後のページ - 10 〜 最後のページの場合、最後のページ - 10 〜 最後のページを表示する
		for i := pageNum - 9; i <= pageNum; i++ {
			pageSlice = append(pageSlice, i)
		}
	} else if pageNum > 10 && page <= pageNum-9 { // 総ページ数が10個より多くてユーザがアクセスしたページが 最後のページ - 10 より前の場合、ユーザがアクセスしたページ 〜 ユーザがアクセスしたページ + 10ページを表示する
		for i := page; i < page+10; i++ {
			pageSlice = append(pageSlice, i)
		}
	} else { // 総ページ数が10個以下の場合
		for i := 1; i <= pageNum; i++ {
			pageSlice = append(pageSlice, i)
		}
	}

	var formattedPosts []map[string]interface{}

	for _, post := range posts {
		formattedPosts = append(formattedPosts, map[string]interface{}{
			"Number":     post.Number,
			"Title":      post.Title,
			"Author":     post.Author,
			"CreatedAt":  post.CreatedAt.Format("2006年01月02日 15:04"),
			"ModifiedAt": post.ModifiedAt.Format("2006年01月02日 15:04"),
			// "CreatedAt":  post.CreatedAt.Format("2006-01-02 15:04"),
			// "ModifiedAt": post.ModifiedAt.Format("2006-01-02 15:04"),
			// "Count":      post.Count,
		})
	}

	return formattedPosts, postNum, pageNum, pageSlice
}

func GetCareerPostContent(postId string) (string, string, string, error) {
	client, err := config.OpensearchNewClient()
	if err != nil {
		log.Fatal("cannot initialize", err)
	}

	ctx := context.Background()
	// infoResp, err := client.Info(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }
	// fmt.Printf("Cluster INFO:\n  Cluster Name: %s\n  Cluster UUID: %s\n  Version Number: %s\n", infoResp.ClusterName, infoResp.ClusterUUID, infoResp.Version.Number)
	var post struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Post   string `json:"post"`
	}
	content := strings.NewReader(fmt.Sprintf(`{
    "query": {
        "ids": {
            "values": ["%s"]
        }
    }
	}`, postId))
	searchResp, err := client.Search(
		ctx,
		&opensearchapi.SearchReq{
			Body: content,
		},
	)
	if err != nil {
		log.Fatal("failed to search document ", err)
	}
	fmt.Printf("Search hits: %v\n", searchResp.Hits.Total.Value)

	if searchResp.Hits.Total.Value > 0 {
		// indices := make([]string, 0)
		for _, hit := range searchResp.Hits.Hits {
			err := json.Unmarshal(hit.Source, &post) // post変数に検索内容で上書き
			if err != nil {
				log.Printf("Error unmarshaling document: %v", err)
				continue
			}

			fmt.Printf("Document ID: %s\n", hit.ID)
			fmt.Printf("  Title: %s\n", post.Title)
			fmt.Printf("  Author: %s\n", post.Author)
			fmt.Printf("  Post: %s\n", post.Post)

			// add := true
			// for _, index := range indices {
			// 	if index == hit.Index {
			// 		add = false
			// 	}
			// }
			// if add {
			// 	indices = append(indices, hit.Index)
			// }
		}
		// fmt.Printf("Search indices: %s\n", strings.Join(indices, ","))
		return post.Post, post.Title, post.Author, nil
	} else {
		return "", "", "", errors.New("not found")
	}
}

func AddCareerPost(db *gorm.DB) {

	// type Posts struct {
	// 	Title  string `json:"title"`
	// 	Author string `json:"author"`
	// 	Post   string `json:"post"`
	// }
	// post := Posts{
	// 	Title:  "What are you doing?",
	// 	Author: "Lee",
	// 	Post:   "I am about to go to the bed.",
	// }
	// docId := postId
	// req := opensearchapi.IndexReq{
	// 	Index:      "career",
	// 	DocumentID: docId, // DocumentIDは内部的に `_id` フィールドとして扱われ、通常 keywordタイプと同様に動作する
	// 	Body:       opensearchutil.NewJSONReader(&post),
	// 	Params: opensearchapi.IndexParams{
	// 		Refresh: "true", // ドキュメントをインデックスに追加した直後に検索可能にする（デフォルトは非同期更新）
	// 	},
	// }
	// insertResp, err := client.Index(ctx, req)
	// if err != nil {
	// 	log.Fatal("failed to insert document ", err)
	// }
	// fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

}

func DeleteCareerPost(postId int, db *gorm.DB) {
	post := models.CareerBoard{Number: postId}
	result := db.Delete(&post)
	fmt.Println("DBから削除されたレコード数:", result.RowsAffected)

	client, err := config.OpensearchNewClient()
	if err != nil {
		log.Fatal("cannot initialize", err)
	}

	ctx := context.Background()
	postIdstr := strconv.Itoa(postId)

	deleteReq := opensearchapi.DocumentDeleteReq{
		Index:      "career",
		DocumentID: postIdstr,
	}

	deleteResponse, err := client.Document.Delete(ctx, deleteReq)
	if err != nil {
		log.Printf("WARNING: OpenSearchからDocumentID[%d]のデータ削除に失敗しました\n", postId)
	}
	fmt.Println("deleteResponse:", deleteResponse)
}

func PostExistCheck(postId int, db *gorm.DB) bool {
	post := models.CareerBoard{Number: postId}
	selectResult := db.Take(&post)
	if errors.Is(selectResult.Error, gorm.ErrRecordNotFound) {
		return false
	} else {
		return true
	}
}

func GetCareerPostDate(postId int, db *gorm.DB) map[string]interface{} {
	var post models.CareerBoard
	db.Select("createdat", "modifiedat").Where("num = ?", postId).Take(&post) // createdatとmodifiedatカラムだけ取得
	// fmt.Println(post.CreatedAt)
	// fmt.Println(post.ModifiedAt)f
	formattedPost := map[string]interface{}{
		"CreatedAt":  post.CreatedAt.Format("2006年01月02日 15:04"),
		"ModifiedAt": post.ModifiedAt.Format("2006年01月02日 15:04"),
	}
	return formattedPost
}

func SearchCareerPost(searchKeywords string, db *gorm.DB) {
	// searchtextを全角/半角スペースで分割し、AND条件で探す
	// OpenSearchからTitleとPostで検索して、ひっかかったもののNumberをスライスに追加し、それをpostNumPerPageに指定した数の分表示する
	// Keyword検索はGoroutineで並列化
	fmt.Println("searchKeyword:", searchKeywords)
	searchKeywordList := controllers.SplitBySpaces(searchKeywords)

	client, err := config.OpensearchNewClient()
	if err != nil {
		log.Fatal("cannot initialize", err)
	}

	ctx := context.Background()

	var HitIDs []string

	for _, searchKeyword := range searchKeywordList {
		// content := strings.NewReader(fmt.Sprintf(`{
		// 	"query": {
		// 			"bool": {
		// 					"should": [
		// 							{ "match": { "title": "%s" }},
		// 							{ "match": { "post": "%s" }}
		// 					]
		// 			}
		// 	}
		// }`, searchKeyword, searchKeyword))
		content := strings.NewReader(fmt.Sprintf(`{
			"query": {
					"bool": {
							"should": [
									{ "match_phrase": { "title": "%s" }},
									{ "match_phrase": { "post": "%s" }}
							]
					}
			}
		}`, searchKeyword, searchKeyword))

		searchResp, err := client.Search(
			ctx,
			&opensearchapi.SearchReq{
				Body: content,
			},
		)
		if err != nil {
			log.Fatal("failed to search document ", err)
		}
		fmt.Printf("Search hits: %v\n", searchResp.Hits.Total.Value)

		if searchResp.Hits.Total.Value > 0 {
			for _, hit := range searchResp.Hits.Hits {

				HitIDs = append(HitIDs, hit.ID)
				fmt.Printf("Document ID: %s\n", hit.ID)

			}
		}
	}

	fmt.Println("HitIDs:", HitIDs)

	// var posts []models.CareerBoard
	// db.Order("num desc").Offset((page - 1) * postNumPerPage).Limit(postNumPerPage).Find(&posts) // 1ページ内にpostはpostNumPerPageに定義した個数まで表示

	// postNum := len(posts)
	// var count int64
	// db.Model(&models.CareerBoard{}).Count(&count)
	// pageNum := int(math.Ceil((float64(count) / float64(postNumPerPage)))) // 総ページ数 ( 1ページ内にpostはpostNumPerPageに定義した個数まで表示)
	// if pageNum == 0 {
	// 	pageNum = 1
	// }
	// var pageSlice []int
	// if pageNum > 10 && page > pageNum-9 { // 総ページ数が10個より多くてユーザがアクセスしたページが 最後のページ - 10 〜 最後のページの場合、最後のページ - 10 〜 最後のページを表示する
	// 	for i := pageNum - 9; i <= pageNum; i++ {
	// 		pageSlice = append(pageSlice, i)
	// 	}
	// } else if pageNum > 10 && page <= pageNum-9 { // 総ページ数が10個より多くてユーザがアクセスしたページが 最後のページ - 10 より前の場合、ユーザがアクセスしたページ 〜 ユーザがアクセスしたページ + 10ページを表示する
	// 	for i := page; i < page+10; i++ {
	// 		pageSlice = append(pageSlice, i)
	// 	}
	// } else { // 総ページ数が10個以下の場合
	// 	for i := 1; i <= pageNum; i++ {
	// 		pageSlice = append(pageSlice, i)
	// 	}
	// }

	// var formattedPosts []map[string]interface{}

	// for _, post := range posts {
	// 	formattedPosts = append(formattedPosts, map[string]interface{}{
	// 		"Number":     post.Number,
	// 		"Title":      post.Title,
	// 		"Author":     post.Author,
	// 		"CreatedAt":  post.CreatedAt.Format("2006年01月02日 15:04"),
	// 		"ModifiedAt": post.ModifiedAt.Format("2006年01月02日 15:04"),
	// 		// "CreatedAt":  post.CreatedAt.Format("2006-01-02 15:04"),
	// 		// "ModifiedAt": post.ModifiedAt.Format("2006-01-02 15:04"),
	// 		// "Count":      post.Count,
	// 	})
	// }

	// return formattedPosts, postNum, pageNum, pageSlice
}
