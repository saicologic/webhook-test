package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	message   string = "メッセージ待機中"
	mu        sync.RWMutex
	clients   = make(map[chan string]bool)
	clientsMu sync.RWMutex
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", getStatus)
	e.GET("/events", sseHandler)
	e.POST("/webhook", webhook)

	e.ServeHTTP(w, r)
}

func getStatus(c echo.Context) error {
	mu.RLock()
	currentMessage := message
	mu.RUnlock()

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Webhook Server</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            background-color: #f0f0f0;
        }
        .container {
            text-align: center;
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .message {
            font-size: 1.5rem;
            color: #333;
            margin-bottom: 1rem;
        }
        .refresh-btn {
            background: #007bff;
            color: white;
            border: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1rem;
        }
        .refresh-btn:hover {
            background: #0056b3;
        }
    </style>
    <script>
        const eventSource = new EventSource('/api/events');
        eventSource.onmessage = function(event) {
            document.querySelector('.message').textContent = event.data;
        };
        eventSource.onerror = function(event) {
            console.log('SSE connection error:', event);
        };
    </script>
</head>
<body>
    <div class="container">
        <div class="message">` + currentMessage + `</div>
        <button class="refresh-btn" onclick="location.reload()">更新</button>
    </div>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}

func sseHandler(c echo.Context) error {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	fmt.Fprintf(w, "data: %s\n\n", currentMessage)
	w.Flush()

	for {
		select {
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}

func broadcastMessage(msg string) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for client := range clients {
		select {
		case client <- msg:
		default:
			close(client)
			delete(clients, client)
		}
	}
}

func webhook(c echo.Context) error {
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

		broadcastMessage(req.Message)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Message updated",
	})
}