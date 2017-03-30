package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

var db *pg.DB

// HomeHandler serves up our basic web page
func HomeHandler(res http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(r)
	if val, ok := req.Header["User-Agent"]; ok {
		log.Printf("User-Agent: %v", val)
	}
	indexPage, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Printf("error occurred reading indexPage: %v\n", err)
	}
	res.Write(indexPage)
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

	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/products", ProductsHandler).Methods("GET", "POST")

	// Basic business
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", router)

	log.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}