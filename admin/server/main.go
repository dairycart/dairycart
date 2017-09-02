package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	cookieName  = "dairycart"
	templateDir = "templates"
	staticDir   = "dist"
)

var (
	apiURL string
	debug  bool
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

func serveLogin(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Login"}
	lp := filepath.Join(templateDir, "login.html")

	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "login.html", p); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// HTTP middleware setting a value on the request context
func cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookies := req.Cookies()
		if len(cookies) == 0 {
			http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		}

		for _, c := range cookies {
			if c.Name == cookieName {
				next.ServeHTTP(res, req)
				return
			}
		}
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
	})
}

func main() {
	debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
	apiURL = os.Getenv("DAIRYCART_API_URL")
	if apiURL == "" {
		log.Fatal("DAIRYCART_API_URL is not set")
	}

	log.Printf(`

		apiURL: %s

	`, apiURL)

	_, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("DAIRYCART_API_URL (%s) is invalid: %v", apiURL, err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))

	FileServer(r, "/static/", http.Dir(staticDir))
	r.Get("/login", serveLogin)
	r.Route("/", func(r chi.Router) {
		r.Use(cookieMiddleware)
		r.Get("/", serveDashboard)
		r.Get("/products", serveProducts)
		r.Get("/product/{sku}", serveProduct)
		r.Get("/orders", serveOrders)
		r.Get("/order/{orderID}", serveOrder)
	})

	port := 1234
	log.Printf("server is listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
