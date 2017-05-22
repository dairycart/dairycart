package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/verygoodsoftwarenotvirus/dairycart/api"
)

// notImplementedHandler is used for endpoints that haven't been implemented yet.
func notImplementedHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
}

func main() {
	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// migrate the database
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error encountered setting up migration instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("error encountered setting up new migration client: %v", err)
	}
	m.Steps(1)

	// setup all our API routes
	APIRouter := mux.NewRouter()
	APIRouter.Host("api.dairycart.com")
	api.SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
