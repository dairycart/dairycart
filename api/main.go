// +build !test

package main

import (
	// stdlib
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// local storage adapter thing
	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/postgres"

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
func migrateDatabase(db *sql.DB, migrationCount int) {
	numberOfUnsuccessfulAttempts := 0
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err := migratePG.WithInstance(db, &migratePG.Config{})
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

	var (
		storageClient storage.Storer
		dbx           *sqlx.DB
		rawDB         *sql.DB
		err           error
	)

	// Connect to the database
	dbChoice := strings.ToLower(os.Getenv("DB_TO_USE"))
	switch dbChoice {
	case "postgres":
		dbURL := os.Getenv("DAIRYCART_DB_URL")
		rawDB, err = sql.Open("postgres", dbURL)
		if err != nil {
			logrus.Fatalf("error encountered connecting to database: %v", err)
		}
		storageClient = &postgres.Postgres{}
		dbx = sqlx.NewDb(rawDB, "postgres")
		dbx.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	default:
		log.Fatalf("invalid database choice: '%s'", dbChoice)
	}

	// migrate the database
	migrationCount := determineMigrationCount()
	migrateDatabase(rawDB, migrationCount)

	secret := os.Getenv("DAIRYSECRET")
	if len(secret) < 32 {
		logrus.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	store := sessions.NewCookieStore([]byte(secret))

	v1APIRouter := chi.NewRouter()

	v1APIRouter.Use(middleware.RequestID)
	v1APIRouter.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))

	SetupAPIRoutes(v1APIRouter, rawDB, dbx, store, storageClient)

	// serve 'em up a lil' sauce
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "healthy!") })

	http.Handle("/", context.ClearHandler(v1APIRouter))
	port := 4321
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
