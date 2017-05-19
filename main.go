package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

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

	router := mux.NewRouter()

	// // https://github.com/go-pg/pg/wiki/FAQ#how-can-i-view-queries-this-library-generates
	// db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
	// 	query, err := event.FormattedQuery()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	log.Printf("%s %s", time.Since(event.StartTime), query)
	// })

	api.SetupAPIRoutes(router, ormDB, properDB)

	// serve 'em up a lil' sauce
	http.Handle("/", router)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
