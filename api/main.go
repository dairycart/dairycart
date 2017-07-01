// +build !test

package main

import (
	// stdlib

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	// dependencies

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattes/migrate"
	migratePG "github.com/mattes/migrate/database/postgres"
	log "github.com/sirupsen/logrus"

	// unnamed dependencies
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
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
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err := migratePG.WithInstance(db.DB, &migratePG.Config{})
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error encountered connecting to database: %v", err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// migrate the database
	migrationCount := determineMigrationCount()
	migrateDatabase(db, migrationCount)

	APIRouter := mux.NewRouter()
	SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "👍") })
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening for requests")
	http.ListenAndServe(":80", nil)
}
