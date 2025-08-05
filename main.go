package main

import (
	"os"

	"webhook-server/pkg"

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
	e.Renderer = pkg.NewTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", pkg.GetStatus)
	e.GET("/message", pkg.GetMessage)
	e.GET("/events", pkg.SSEHandler)
	e.POST("/webhook", pkg.WebhookHandler)

	e.Logger.Fatal(e.Start(":" + port))
}