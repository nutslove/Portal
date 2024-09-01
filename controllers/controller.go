package controllers

import (
	"html/template"
	"net/http"
	"strings"
	"unicode"

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

func SplitBySpaces(s string) []string {
	// 全角と半角のスペースで分割
	f := func(c rune) bool {
		return unicode.IsSpace(c) || c == '　' // '　'は全角スペース
	}
	return strings.FieldsFunc(strings.TrimSpace(s), f) // strings.FieldsFuncはfがtrueを返す文字（今回の場合は半角または全角のスペース）で文字列を分割し、結果を文字列のスライスとして返す
}
