package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/verygoodsoftwarenotvirus/dairycart/api"
)

// notImplementedHandler is used for endpoints that haven't been implemented yet.
func notImplementedHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
}

func main() {
	// init stuff
	domainName := os.Getenv("DAIRYCART_DOMAIN")
	if domainName == "" {
		domainName = "localhost"
	}

	dbURL := os.Getenv("DAIRYCART_DB_URL")
	dbOptions, err := pg.ParseURL(dbURL)
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}
	ormDB := pg.Connect(dbOptions)

	properDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	APIRouter := mux.NewRouter()
	APIRouter.Host("api.dairycart.com")

	api.SetupAPIRoutes(APIRouter, ormDB, properDB)

	// serve 'em up a lil' sauce
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
