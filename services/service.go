package services

import (
	"portal/models"
)

func UserExistCheck(userid string) bool {
	db := models.ConnectDB()
	user := models.UserData{Nickname: userid}
	selectResult := db.First(&user)
	if selectResult.Error != nil {
		return false
	} else {
		return true
	}
}
