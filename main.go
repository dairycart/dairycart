package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

var db *pg.DB

// HomeHandler serves up our basic web page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	rv := struct{ OK bool }{OK: true}
	json.NewEncoder(w).Encode(rv)
}

// ProductsHandler is a generic product list response handler
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		var products []Product
		err := db.Model(&products).Select()
		if err != nil {
			log.Fatalf("Error encountered querying for products: %v", err)
		}
		json.NewEncoder(w).Encode(products)
	case "POST":
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}

		newProduct := &Product{}
		err := json.NewDecoder(r.Body).Decode(newProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
			// fmt.Fprintf(w, "Error encountered parsing request: %v", err)
		}

		err = db.Insert(newProduct)
		if err != nil {
			log.Printf("error inserting product into database: %v", err)
		}
	}
}

func main() {
	// init stuff
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	dbOptions, err := pg.ParseURL(dbURL)
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}
	db = pg.Connect(dbOptions)
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/products", ProductsHandler).Methods("GET", "POST")

	// Basic business
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
