package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

const (
	templateDir = "../templates"
	staticDir   = "../dist"
)

type Page struct {
	Title string
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	// yuge thanks to Alex Edwards: http://www.alexedwards.net/blog/serving-static-sites-with-go
	p := &Page{Title: "Dashboard"}
	lp := filepath.Join(templateDir, "base.html")
	fp := filepath.Join(templateDir, "index.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", p); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	r := chi.NewRouter()
	FileServer(r, "/static/", http.Dir(staticDir))
	r.Get("/", serveDashboard)
	r.Get("/products", serveProducts)
	r.Get("/products/{sku}", serveProduct)
	r.Get("/orders", serveOrders)
	r.Get("/order/{orderID}", serveOrder)
	http.ListenAndServe(":3000", r)
}
