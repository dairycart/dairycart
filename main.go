package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

var db *pg.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	rv := struct{ OK bool }{OK: true}
	json.NewEncoder(w).Encode(rv)
}

func main() {
	db = pg.Connect(&pg.Options{
		User:     os.Getenv("DAIRYBASE_USER"),
		Password: os.Getenv("DAIRYBASE_PASS"),
		Database: os.Getenv("DAIRYBASE_DATABASE"),
	})

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
