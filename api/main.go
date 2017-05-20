package api

import (
	"database/sql"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, ormDB *pg.DB, properDB *sql.DB) {
	// Products
	router.HandleFunc("/product", buildProductCreationHandler(ormDB)).Methods("POST")
	router.HandleFunc("/products", buildProductListHandler(properDB)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildSingleProductHandler(properDB)).Methods("GET")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductUpdateHandler(ormDB)).Methods("PUT")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductExistenceHandler(properDB)).Methods("HEAD")
	router.HandleFunc("/product/{sku:[a-zA-Z]+}", buildProductDeletionHandler(properDB)).Methods("DELETE")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(ormDB)).Methods("POST")
}
