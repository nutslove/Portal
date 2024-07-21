package models

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserData struct {
	Nickname string `gorm:"size:30;primaryKey;column:userid"`
	Password string `gorm:"size:30;column:password;not null"`
	Age      int    `gorm:"size:3;column:age"`
	Company  string `gorm:"size:50;column:company"`
	Role     string `gorm:"size:50;column:role"`
}

var (
	DBUser   string = os.Getenv("MYSQL_USER")
	Password string = os.Getenv("MYSQL_PASSWORD")
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/portal?charset=utf8mb4&parseTime=True&loc=Local", DBUser, Password)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&UserData{})
	return db
}
