package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(hashedPassword, requestPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword))
}

func NotFoundResponse(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func UnAuthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"result":  false,
		"message": "you need to login first!",
	})
}

func RegisterCustomFunction(r *gin.Engine) {
	funcMap := template.FuncMap{
		"add": func(x, y int) int {
			return x + y
		},
		"subtract": func(x, y int) int {
			return x - y
		},
	}

	r.SetFuncMap(funcMap)
}
