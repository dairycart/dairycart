package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/verygoodsoftwarenotvirus/dairycart/api"
)

// notImplementedHandler is used for endpoints that haven't been implemented yet.
func notImplementedHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
}

func main() {
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	APIRouter := mux.NewRouter()
	APIRouter.Host("api.dairycart.com")
	api.SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
