// +build !test

package main

import (
	// stdlib
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// dependencies
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattes/migrate"
	migratePG "github.com/mattes/migrate/database/postgres"
	"github.com/sirupsen/logrus"

	// unnamed dependencies
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
func migrateDatabase(db *sqlx.DB, migrationCount int) {
	numberOfUnsuccessfulAttempts := 0
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err := migratePG.WithInstance(db.DB, &migratePG.Config{})
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
			numberOfUnsuccessfulAttempts++

			if numberOfUnsuccessfulAttempts == maxConnectionAttempts {
				log.Fatal("Failed to connect to the database")
			}
		} else {
			migrationsDir := os.Getenv("DAIRYCART_MIGRATIONS_DIR")
			m, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
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
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		logrus.Fatalf("error encountered connecting to database: %v", err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// migrate the database
	migrationCount := determineMigrationCount()
	migrateDatabase(db, migrationCount)

	secret := os.Getenv("DAIRYSECRET")
	if len(secret) < 32 {
		logrus.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	store := sessions.NewCookieStore([]byte(secret))

	v1APIRouter := chi.NewRouter()

	v1APIRouter.Use(middleware.RequestID)
	v1APIRouter.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))

	SetupAPIRoutes(v1APIRouter, db, store)

	// serve 'em up a lil' sauce
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "healthy!") })

	http.Handle("/", context.ClearHandler(v1APIRouter))
	port := 80
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
