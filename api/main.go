// +build !test

package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dairycart/dairycart/storage/images/local"
	"github.com/dairycart/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

func buildServerConfig() *ServerConfig {
	// Connect to the database
	dbChoice := strings.ToLower(os.Getenv("DB_TO_USE"))
	switch dbChoice {
	case "postgres":
		dbURL := os.Getenv("DAIRYCART_DB_URL")
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			logrus.Fatalf("error encountered connecting to database: %v", err)
		}

		loadExampleData := os.Getenv("MIGRATE_EXAMPLE_DATA") == "YES"
		pg := postgres.NewPostgres()
		if err = pg.Migrate(db, dbURL, loadExampleData); err != nil {
			logrus.Fatalf("error encountered migrating database: %v", err)
		}

		return &ServerConfig{
			DB:              db,
			Dairyclient:     postgres.NewPostgres(),
			WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
			ImageStorer:     &local.LocalImageStorer{BaseURL: "http://localhost:4321"},
		}
	default:
		logrus.Fatalf("invalid database choice: '%s'", dbChoice)
	}
	return nil
}

func setupCookieStorage() *sessions.CookieStore {
	secret := os.Getenv("DAIRYSECRET")
	if len(secret) < 32 {
		logrus.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	return sessions.NewCookieStore([]byte(secret))
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
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

func createPhotoDirectory(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	config := buildServerConfig()
	config.CookieStore = setupCookieStorage()
	config.Router = chi.NewRouter()

	config.Router.Use(middleware.RequestID)
	config.Router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
	SetupAPIRoutes(config)

	photoDir := "product_images"
	createPhotoDirectory(photoDir)
	fileServer(config.Router, fmt.Sprintf("/%s/", photoDir), http.Dir(photoDir))

	port := 4321
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "healthy!") })
	http.Handle("/", context.ClearHandler(config.Router))
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
