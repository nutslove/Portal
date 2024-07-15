package main

import (
	"github.com/gin-gonic/gin"
	"portal/routers"
)

func main() {
	router := gin.Default()

	// 静的ファイルの設定
	router.Static("/static", "./static")

	// テンプレートファイルの読み込み
	router.LoadHTMLGlob("templates/*")

	// ルーティングの設定
	routers.SetupRouter(router)

	// サーバーの起動
	router.Run(":8080")
}
