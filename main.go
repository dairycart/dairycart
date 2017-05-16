package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

var db *pg.DB

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
	db = pg.Connect(dbOptions)
	router := mux.NewRouter()

	// // https://github.com/go-pg/pg/wiki/FAQ#how-can-i-view-queries-this-library-generates
	// db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
	// 	query, err := event.FormattedQuery()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	log.Printf("%s %s", time.Since(event.StartTime), query)
	// })

	// Base Products
	router.HandleFunc("/base_product/{id:[0-9]+}", SingleBaseProductHandler).Methods("GET")

	// Products
	router.HandleFunc("/products", ProductListHandler).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", ProductExistenceHandler).Methods("HEAD")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", SingleProductHandler).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", ProductUpdateHandler).Methods("PUT")
	router.HandleFunc("/product", ProductCreationHandler).Methods("POST")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", ProductDeletionHandler).Methods("DELETE")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", ProductAttributeValueCreationHandler).Methods("POST")

	// Orders
	router.HandleFunc("/orders", OrderListHandler).Methods("GET")
	router.HandleFunc("/order", OrderCreationHandler).Methods("POST")

	// serve 'em up a lil' sauce
	http.Handle("/", router)
	log.Println("Dairycart now listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
