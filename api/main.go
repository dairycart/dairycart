// +build !test

package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dairycart/dairycart/api/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	migratePG "github.com/mattes/migrate/database/postgres"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

const (
	maxConnectionAttempts = 25
)

func determineMigrationCount() int {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatalf("missing migrations files")
	}

	migrationCount := 0
	migrateExampleData := os.Getenv("MIGRATE_EXAMPLE_DATA")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			migrationCount++
		}
	}
	if migrateExampleData != "YES" {
		migrationCount--
	}

	return migrationCount
}

// this function not only waits for the database to accept its incoming connection, but also performs any necessary migrations
func migrateDatabase(db *sql.DB) {
	var err error
	var driver database.Driver

	migrationCount := determineMigrationCount()
	numberOfUnsuccessfulAttempts := 0
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err = migratePG.WithInstance(db, &migratePG.Config{})
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
			numberOfUnsuccessfulAttempts++

			if numberOfUnsuccessfulAttempts == maxConnectionAttempts {
				log.Fatal("Failed to connect to the database")
			}
		} else {
			var m *migrate.Migrate
			migrationsDir := os.Getenv("DAIRYCART_MIGRATIONS_DIR")
			m, err = migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
			if err != nil {
				log.Fatalf("error encountered setting up new migration client: %v", err)
			}

			for i := 0; i < migrationCount; i++ {
				err = m.Steps(1)
				if err != nil {
					log.Printf("error encountered migrating database: %v", err)
					break
				}
			}
			databaseIsNotMigrated = false
		}
	}
	if err != nil {
		log.Printf("error encountered migrating database: %v", err)
	}
	log.Println("database migrated!")
}

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

		return &ServerConfig{
			DB:              db,
			Dairyclient:     postgres.NewPostgres(),
			WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
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

	// migrate the database
	migrateDatabase(config.DB)

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
