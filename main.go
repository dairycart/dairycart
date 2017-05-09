package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

var db *pg.DB
var templates = template.Must(template.ParseGlob("templates/*"))

// HomeHandler serves up our basic web page
func HomeHandler(res http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(r)
	if val, ok := req.Header["User-Agent"]; ok {
		log.Printf("User-Agent: %v", val)
	}
	indexPage, err := ioutil.ReadFile("templates/home.html")
	if err != nil {
		log.Printf("error occurred reading indexPage: %v\n", err)
	}
	renderTemplates(res, "Dairycart", string(indexPage))
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

	// Basic business
	router.HandleFunc("/", HomeHandler).Methods("GET")

	// Products
	router.HandleFunc("/product", ProductCreationHandler).Methods("POST")
	router.HandleFunc("/products", ProductListHandler).Methods("GET")

	// Orders
	router.HandleFunc("/order", OrderCreationHandler).Methods("POST")
	router.HandleFunc("/orders", OrderListHandler).Methods("GET")

	// let's start servin'
	http.Handle("/", router)
	log.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
