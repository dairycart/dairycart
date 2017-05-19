package api

import (
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *pg.DB) {
	// Base Products
	router.HandleFunc("/base_product/{id:[0-9]+}", buildSingleBaseProductHandler(db)).Methods("GET")

	// Products
	router.HandleFunc("/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc("/product", buildProductCreationHandler(db)).Methods("POST")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(db)).Methods("POST")

	// Orders
	router.HandleFunc("/orders", buildOrderListHandler(db)).Methods("GET")
	router.HandleFunc("/order", buildOrderCreationHandler(db)).Methods("POST")

}
