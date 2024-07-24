package controllers

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(hashedPassword, requestPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword))
}
