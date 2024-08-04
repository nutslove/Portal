package main

import (
	"portal/config"
	"portal/controllers"
	"portal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config.OtelInitialSetting()

	// 静的ファイルの設定
	router.Static("/static", "./static")

	// テンプレート内で使うカスタム関数を登録
	controllers.RegisterCustomFunction(router)

	// テンプレートファイルの読み込み
	router.LoadHTMLGlob("templates/*")

	// ルーティングの設定
	routers.SetupRouter(router)

	// サーバーの起動
	router.Run(":8080")
}
