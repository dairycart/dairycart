package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Order describes, well... orders.
type Order struct {
	ID     int64       `json:"id"`
	ShipTo Customer    `json:"ship_to"`
	BillTo Customer    `json:"bill_to"`
	Lines  []OrderLine `json:"lines"`
}

// OrderLine represents a product in an order
type OrderLine struct {
	ID       int64   `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

// Customer describes a user that places an order
type Customer struct {
	ID int64
}

// OrdersResponse is a order response struct
type OrdersResponse struct {
	ListResponse
	Data []Order `json:"data"`
}

// OrderListHandler is a generic order list request handler
func OrderListHandler(res http.ResponseWriter, req *http.Request) {
	var orders []Order
	ordersModel := db.Model(&orders)

	pager, err := genericListQueryHandler(req, ordersModel, genericActiveFilter)
	if err != nil {
		informOfServerIssue(err, "Error encountered querying for orders", res)
		return
	}

	ordersResponse := &OrdersResponse{
		ListResponse: ListResponse{
			Page:  pager.Page(),
			Limit: pager.Limit(),
			Count: len(orders),
		},
		Data: orders,
	}
	json.NewEncoder(res).Encode(ordersResponse)
}

// OrderCreationHandler is a order creation handler
func OrderCreationHandler(res http.ResponseWriter, req *http.Request) {
	newOrder := &Order{}
	bodyIsInvalid := ensureRequestBodyValidity(res, req, newOrder)
	if bodyIsInvalid {
		return
	}

	err := db.Insert(newOrder)
	if err != nil {
		log.Printf("error inserting product into database: %v", err)
	}
}
