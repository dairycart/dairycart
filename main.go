package main

import (
	// stdlib
	"database/sql"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	// dependencies
	"github.com/gorilla/mux"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	log "github.com/sirupsen/logrus"

	// unnamed dependencies
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"

	// homegrown
	"github.com/verygoodsoftwarenotvirus/dairycart/api"
)

var db *sql.DB

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func determineMigrationCount() int {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatalf("missing migrations files")
	}

	migrationCount := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			migrationCount++
		}
	}
	return migrationCount
}

// this function not only waits for the database to accept its incoming connection, but also performs any necessary migrations
func migrateDatabase(db *sql.DB, migrationCount int) {
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
		} else {
			migrationsDir := "file://migrations" // os.Getenv("DAIRYCART_MIGRATIONS_DIR")
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
	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error encountered connecting to database: %v", err)
	}

	// migrate the database
	migrationCount := determineMigrationCount()
	migrateDatabase(db, migrationCount)

	// setup all our API routes
	APIRouter := mux.NewRouter()
	// APIRouter.Host("api.dairycart.com")
	api.SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "I live!") })
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening for requests")
	http.ListenAndServe(":80", nil)
}
