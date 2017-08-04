package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Page struct {
	Title string
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// yuge thanks to Alex Edwards: http://www.alexedwards.net/blog/serving-static-sites-with-go
	path := filepath.Clean(r.URL.Path)

	pathToTitleMap := map[string]string{
		"/":         "Dashboard",
		"/products": "Products",
	}

	p := &Page{
		Title: pathToTitleMap[path],
	}
	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", fmt.Sprintf("%s.html", filepath.Clean(path)))
	if path == "/" {
		fp = filepath.Join("templates", "index.html")
	}

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", p); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	compiledFileServer := http.FileServer(http.Dir("dist"))
	http.Handle("/static/", http.StripPrefix("/static/", compiledFileServer))

	http.HandleFunc("/", serveTemplate)
	http.ListenAndServe(":3000", nil)
}
