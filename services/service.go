package services

import (
	"net/http"
	"portal/controllers"
	"portal/models"
	"strconv"
	"time"

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

func GetCareerPosts(db *gorm.DB) []models.CareerBoard {
	// フォーマットと日付文字列
	layout := "2006-01-02 15:04:05"
	date1 := "2024-08-02 02:59:59"
	date2 := "2024-11-30 11:59:59"

	dateafter1, err := time.Parse(layout, date1)
	if err != nil {
		panic(err)
	}

	dateafter2, err := time.Parse(layout, date2)
	if err != nil {
		panic(err)
	}

	Posts := []models.CareerBoard{
		{Number: 1, Title: "なんでもいい", Author: "nutslove", Date: dateafter1, Count: 5},
		{Number: 12345678, Title: "転職したい", Author: "kumi", Date: dateafter2, Count: 50000},
	}
	return Posts
}
