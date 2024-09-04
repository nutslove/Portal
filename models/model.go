package models

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserData struct {
	Nickname string `gorm:"size:30;primaryKey;column:userid"`
	Password string `gorm:"size:60;column:password;not null"` // bcryptで生成されたハッシュ値は通常60文字固定長のため
	Age      *int   `gorm:"size:3;column:age"`                // Ageの入力がない場合、0ではなくnullで登録するためにPoint型に変更
	Company  string `gorm:"size:50;column:company"`
	Role     string `gorm:"size:50;column:role"`
}

type CareerBoard struct {
	Number     int       `gorm:"primaryKey;column:num"`
	Title      string    `gorm:"size:100;column:title"`
	Author     string    `gorm:"size:30;column:author"`
	CreatedAt  time.Time `gorm:"type:datetime;column:createdat"`
	ModifiedAt time.Time `gorm:"type:datetime;column:modifiedat"`
	// Count      int       `gorm:"column:count"`
}

type AnythingBoard struct {
	Number     int       `gorm:"primaryKey;column:num"`
	Title      string    `gorm:"size:100;column:title"`
	Author     string    `gorm:"size:30;column:author"`
	CreatedAt  time.Time `gorm:"type:datetime;column:createdat"`
	ModifiedAt time.Time `gorm:"type:datetime;column:modifiedat"`
	// Count      int       `gorm:"column:count"`
}

type CareerBoardComment struct {
	ID            int       `gorm:"primaryKey;column:id"`
	Number        int       `gorm:";column:num"`
	Author        string    `gorm:"size:30;column:author"`
	Date          time.Time `gorm:"type:datetime;column:date"`
	CommentNumber int       `gorm:":column:commentnum"`
}

var (
	DBUser   string = os.Getenv("MYSQL_USER")
	Password string = os.Getenv("MYSQL_PASSWORD")
)

func ConnectDB() *gorm.DB {
	timezone := url.QueryEscape("Asia/Tokyo") // 文字列をURL用にエスケープ（この例だと"Asia%2FTokyo"を返す）
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/portal?charset=utf8mb4&parseTime=True&loc=%s", DBUser, Password, timezone)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&UserData{})
	db.AutoMigrate(&CareerBoard{})
	db.AutoMigrate(&AnythingBoard{})
	db.AutoMigrate(&CareerBoardComment{})
	return db
}
