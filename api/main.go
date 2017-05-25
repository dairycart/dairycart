package api

import (
	"database/sql"
	"fmt"

	"github.com/gorilla/mux"
)

const (
	skuRegexValidation = `[a-zA-Z\\-_]+`
)

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	productEndpoint := fmt.Sprintf("/product/{sku:%s}", skuRegexValidation)

	router.HandleFunc("/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Attribute Values
	router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(db)).Methods("POST")
}
