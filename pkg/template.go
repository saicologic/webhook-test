package pkg

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


func NewTemplate() *Template {
	// 外部ファイルを使用（ローカル・Vercel共通）
	templates := template.Must(template.ParseGlob("templates/*.html"))
	return &Template{
		templates: templates,
	}
}