package handler

import (
	"net/http"

	"webhook-server/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	// テンプレートエンジンを設定
	e.Renderer = internal.NewTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", internal.GetStatus)
	e.GET("/events", internal.SSEHandler)
	e.POST("/webhook", internal.WebhookHandler)

	e.ServeHTTP(w, r)
}