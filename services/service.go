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
