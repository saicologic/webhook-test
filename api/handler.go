package handler

import (
	"net/http"

	"webhook-server/pkg"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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

	e.ServeHTTP(w, r)
}