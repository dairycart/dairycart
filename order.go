package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Order describes, well... orders.
type Order struct {
	ID int64
}

// Customer describes a user that places an order
type Customer struct {
	ID int64
}

// OrderListHandler is a generic order list request handler
func OrderListHandler(res http.ResponseWriter, req *http.Request) {
	var orders []Order
	err := db.Model(&orders).Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
	}
	json.NewEncoder(res).Encode(orders)
}

// OrderCreationHandler is a order creation handler
func OrderCreationHandler(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(res, "Please send a request body", http.StatusBadRequest)
		return
	}

	newOrder := &Order{}
	err := json.NewDecoder(req.Body).Decode(newOrder)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
		// fmt.Fprintf(w, "Error encountered parsing request: %v", err)
	}

	err = db.Insert(newOrder)
	if err != nil {
		log.Printf("error inserting product into database: %v", err)
	}
}
