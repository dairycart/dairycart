package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/dairycart/dairycart/cmd/admin/v1/server/html/statik"
	"github.com/rakyll/statik/fs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	statikFS     http.FileSystem
	debug        bool
	apiServerURL *url.URL
)

const (
	cookieName = "dairycart"
	staticDir  = "assets"
)

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

func informUserOfFileReadError(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(res).Encode(struct {
		Response string `json:"error"`
	}{fmt.Sprintf("Error encountered reading local file: %v", err)})
}

func serveHomePage(res http.ResponseWriter, req *http.Request) {
	homepage, err := ioutil.ReadFile("html/base.html")
	if err != nil {
		informUserOfFileReadError(res, err)
		return
	}
	res.Write(homepage)
}

func serveLogin(res http.ResponseWriter, req *http.Request) {
	homepage, err := ioutil.ReadFile("html/login.html")
	if err != nil {
		informUserOfFileReadError(res, err)
		return
	}
	res.Write(homepage)
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

func informUserOfForwardingError(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(res).Encode(struct {
		Response string `json:"error"`
	}{fmt.Sprintf("Error encountered forwarding request to API server: %v", err)})
}

func apiForwarder(res http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(fmt.Sprintf("%s?%s", strings.Replace(req.URL.Path, "/api", "", 1), req.URL.Query().Encode()))
	toForwardTo := apiServerURL.ResolveReference(u)

	req, err := http.NewRequest(req.Method, toForwardTo.String(), req.Body)
	if err != nil {
		informUserOfForwardingError(res, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		informUserOfForwardingError(res, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		informUserOfForwardingError(res, err)
		return
	}

	res.WriteHeader(resp.StatusCode)
	res.Write(body)
}

func main() {
	debug = strings.ToLower(os.Getenv("DEBUG")) == "true"

	apiURL := os.Getenv("DAIRYCART_API_URL")

	var err error
	statikFS, err = fs.New()
	if err != nil {
		log.Fatalf("Error initializing static files: %v\n", err)
	}

	apiServerURL, err = url.Parse(apiURL)
	if err != nil {
		log.Fatal("API server URL is invalid")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))

	FileServer(r, "/assets/", statikFS)
	r.Get("/login", serveLogin)
	r.Route("/", func(r chi.Router) {
		// commented out currently for debugging reasons
		if !debug {
			r.Use(cookieMiddleware)
		}
		r.Get("/", serveHomePage)
		r.HandleFunc("/api/*", apiForwarder)
	})

	port := ":1234"
	log.Printf("server is listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
