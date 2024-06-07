package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// add ftpcrawl.go from ../internal/ftpcrawl.go
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type FormData struct {
	Values, Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	// var text string = "duh"
	// fmt.Printf("This is a test of my go code. %s\n", text)

	// ftpCrawl("92.204.128.116")
	// crawler.FtpCrawl("92.204.128.116")

	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {

		return c.Render(200, "index.html", page)
	})

	e.POST("/ftp", func(c echo.Context) error {
		host := c.FormValue("host")
		user := c.FormValue("user")
		password := c.FormValue("password")
		path := c.FormValue("path")
		terms := c.FormValue("terms")

		return c.Render(200, "results", page)
	})
	e.Logger.Fatal(e.Start(":42069"))
}
