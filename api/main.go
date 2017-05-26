package api

import (
	"database/sql"
	"fmt"

	"github.com/gorilla/mux"
)

const (
	// SKUPattern represents the valid characters a sku can contain
	SKUPattern = `[a-zA-Z\-_]+`
)

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	productEndpoint := fmt.Sprintf("/v1/product/{sku:%s}", SKUPattern)
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods("POST")
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(db)).Methods("POST")
}
