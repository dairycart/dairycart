package main

import (
	// stdlib
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// dependencies
	"github.com/gorilla/mux"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"

	// unnamed dependencies
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"

	// homegrown
	"github.com/verygoodsoftwarenotvirus/dairycart/api"
)

// notImplementedHandler is used for endpoints that haven't been implemented yet.
func notImplementedHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
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
			m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
			if err != nil {
				log.Fatalf("error encountered setting up new migration client: %v", err)
			}

			for i := 0; i < migrationCount; i++ {
				err = m.Steps(1)
				if err != nil {
					log.Printf("error encountered migrating database: %v", err)
				}
			}
			databaseIsNotMigrated = false
		}
	}
}

func main() {
	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error encountered connecting to database: %v", err)
	}

	// migrate the database
	migrationCount := determineMigrationCount()
	log.Printf("determined migrationCount to be %d", migrationCount)
	migrateDatabase(db, migrationCount)

	// setup all our API routes
	APIRouter := mux.NewRouter()
	APIRouter.Host("api.dairycart.com")
	api.SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
