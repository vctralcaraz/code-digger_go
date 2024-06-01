package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

type Count struct {
	Count int
}

func main() {
	// var text string = "duh"
	// fmt.Printf("This is a test of my go code. %s\n", text)

	// ftpCrawl("92.204.128.116")

	e := echo.New()
	e.Use(middleware.Logger())

	count := Count{Count: 0}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", count)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "index.html", count)
	})

	e.Logger.Fatal(e.Start(":42069"))
}

// file struct is used to store the file path and the terms to search for
type file struct {
	path  string   // path to the file
	terms [][]byte // list of terms
}

func ftpCrawl(path string) {
	var found = []file{}
	c, err := ftp.Dial(path + ":21")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	err = c.Login("drpollock", "s9X[DJ$kP+LZ0h97A7")
	if err != nil {
		log.Fatal(err)
	}

	w := c.Walk("/public_html/wp-content/themes")

	for w.Next() {
		if w.Err() != nil {
			log.Fatal(w.Err())
		}
		if w.Stat().Type.String() == "folder" {
			entries, err := c.List(w.Path())
			if err != nil {
				log.Fatal(err)
			}

			for _, entry := range entries {
				if strings.HasSuffix(entry.Name, ".php") {
					fmt.Println(w.Path() + "/" + entry.Name)
					r, err := c.Retr(w.Path() + "/" + entry.Name)
					if err != nil {
						log.Fatal(err)
					}

					buf, err := io.ReadAll(r)
					if err != nil {
						log.Fatal(err)
					}

					// regex to find terms in each file
					// (?i) = case insensitive
					// TODO: add dynamic terms into the regex
					re := regexp.MustCompile("(?i)Ivan|sample|Team")
					if re.FindAll(buf, -1) != nil {
						found = append(found, file{
							path:  w.Path() + "/" + entry.Name,
							terms: re.FindAll(buf, -1),
						})
					}

					err = r.Close()
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
	fmt.Println("Results:")
	fmt.Println("")
	for _, v := range found {
		fmt.Println(v.path)
		for _, term := range v.terms {
			fmt.Println(string(term))
		}
	}
	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}
