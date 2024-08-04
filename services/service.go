package services

import (
	"net/http"
	"portal/controllers"
	"portal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserExistCheck(userid string, db *gorm.DB) bool {
	user := models.UserData{Nickname: userid}
	selectResult := db.Take(&user)
	if selectResult.Error != nil {
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

func GetCareerPosts(page int, db *gorm.DB) ([]map[string]interface{}, int) {

	// var posts []models.CareerBoard
	var posts []models.AnythingBoard
	db.Order("num desc").Offset((page - 1) * 15).Limit(15).Find(&posts)

	postNum := len(posts)
	var formattedPosts []map[string]interface{}

	for _, post := range posts {
		formattedPosts = append(formattedPosts, map[string]interface{}{
			"Number": post.Number,
			"Title":  post.Title,
			"Author": post.Author,
			"Date":   post.Date.Format("2006-01-02 15:04"),
			"Count":  post.Count,
		})
	}

	return formattedPosts, postNum
}

func AddCareerPost(db *gorm.DB) {

}

func DeleteCareerPost() {

}
