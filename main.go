package main

import (
	"os"

	"webhook-server/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	e := echo.New()

	// テンプレートエンジンを設定
	e.Renderer = internal.NewTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", internal.GetStatus)
	e.GET("/events", internal.SSEHandler)
	e.POST("/webhook", internal.WebhookHandler)

	e.Logger.Fatal(e.Start(":" + port))
}