package pkg

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var templateFS embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	templates := template.Must(template.ParseFS(templateFS, "templates/*.html"))
	return &Template{
		templates: templates,
	}
}