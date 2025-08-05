package pkg

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	message   string = "メッセージ待機中"
	mu        sync.RWMutex
	clients   = make(map[chan string]bool)
	clientsMu sync.RWMutex
)

func GetStatus(c echo.Context) error {
	mu.RLock()
	currentMessage := message
	mu.RUnlock()

	data := struct {
		Message string
	}{
		Message: currentMessage,
	}

	return c.Render(http.StatusOK, "index", data)
}

// Polling用のAPI - 現在のメッセージを返す
func GetMessage(c echo.Context) error {
	mu.RLock()
	currentMessage := message
	mu.RUnlock()

	return c.JSON(http.StatusOK, map[string]string{
		"message": currentMessage,
	})
}

func SSEHandler(c echo.Context) error {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // Nginxのバッファリング無効化

	clientChan := make(chan string, 10)

	clientsMu.Lock()
	clients[clientChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, clientChan)
		clientsMu.Unlock()
		close(clientChan)
	}()

	mu.RLock()
	currentMessage := message
	mu.RUnlock()

	// 初期メッセージ送信
	fmt.Fprintf(w, "data: %s\n\n", currentMessage)
	w.Flush()

	// keepalive用のticker（30秒間隔）
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.Flush()
		case <-ticker.C:
			// keepalive ping（クライアントには表示されない）
			fmt.Fprintf(w, ": keepalive\n\n")
			w.Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}

func BroadcastMessage(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	// 切断されたクライアントを削除するためのスライス
	var toDelete []chan string

	for client := range clients {
		select {
		case client <- msg:
			// メッセージ送信成功
		default:
			// チャンネルがブロックされている場合は切断されたクライアント
			toDelete = append(toDelete, client)
		}
	}

	// 切断されたクライアントを削除
	for _, client := range toDelete {
		delete(clients, client)
		close(client)
	}
}

func WebhookHandler(c echo.Context) error {
	type WebhookRequest struct {
		Message string `json:"message" form:"message"`
	}

	req := new(WebhookRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.Message != "" {
		mu.Lock()
		message = req.Message
		mu.Unlock()

		BroadcastMessage(req.Message)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Message updated",
	})
}