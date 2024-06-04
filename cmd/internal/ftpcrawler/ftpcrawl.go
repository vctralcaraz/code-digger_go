package crawler

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/jlaffaye/ftp"
)

// file struct is used to store the file path and the terms to search for
type file struct {
	path  string   // path to the file
	terms [][]byte // list of terms
}

func FtpCrawl(path string) {
	var found = []file{}
	c, err := ftp.Dial(path + ":21")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// TODO: Add dynamic login through htmx
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
