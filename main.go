package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zedjones/redirectprotect/internal"
	"github.com/zedjones/redirectprotect/routes"
)

/*
TemplateRenderer is exported for the Echo template renderer to use
*/
type TemplateRenderer struct {
	templates *template.Template
}

/*
Render is exported for the Echo template renderer to use
*/
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	go internal.AddChecks()
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("frontend/build/*.html")),
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Renderer = renderer
	e.Static("/", "frontend/build")
	e.POST("/add_redirect", routes.RegisterURL)
	e.POST("/check_pass", routes.CheckPassphrase)
	e.GET("*", routes.GetRedirect)
	e.Start(":1234")
}
