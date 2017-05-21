package api

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	router.HandleFunc("/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductDeletionHandler(db)).Methods("DELETE")
	router.HandleFunc("/product/{progenitor_id:[0-9]+}", buildProductCreationHandler(db)).Methods("POST")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(db)).Methods("POST")
}
