package internal

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

const indexTemplate = `<!DOCTYPE html>
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
        const basePath = window.location.pathname.includes('/api') ? '/api' : '';
        const eventSource = new EventSource(basePath + '/events');
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
        <div class="message">{{.Message}}</div>
        <button class="refresh-btn" onclick="location.reload()">更新</button>
    </div>
</body>
</html>`

func NewTemplate() *Template {
	// ローカル環境では外部ファイルを使用、Vercelでは埋め込みテンプレートを使用
	if templates := tryLoadExternalTemplate(); templates != nil {
		return &Template{templates: templates}
	}
	
	// Vercel環境用の埋め込みテンプレート
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		panic(err)
	}
	return &Template{
		templates: tmpl,
	}
}

func tryLoadExternalTemplate() *template.Template {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil
	}
	return templates
}