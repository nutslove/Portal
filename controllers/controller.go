package controllers

import (
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
