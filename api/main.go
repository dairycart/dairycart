package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gorilla/mux"
)

const (
	// SKUPattern represents the valid characters a sku can contain
	SKUPattern = `[a-zA-Z\-_]+`
)

func buildRoute(routeParts ...string) string {
	allRouteParts := append([]string{"v1"}, routeParts...)
	return fmt.Sprintf("/%s", strings.Join(allRouteParts, "/"))
}

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", SKUPattern))
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods("POST")
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Options
	productOptionEndpoint := buildRoute("product_options", "{progenitor_id:[0-9]+}")
	specificOptionEndpoint := buildRoute("product_options", "{option_id:[0-9]+}")
	router.HandleFunc(productOptionEndpoint, buildProductOptionListHandler(db)).Methods("GET")
	router.HandleFunc(productOptionEndpoint, buildProductOptionCreationHandler(db)).Methods("POST")
	router.HandleFunc(specificOptionEndpoint, buildProductOptionUpdateHandler(db)).Methods("PUT")

	// Product Option Values
	optionValueEndpoint := buildRoute("product_options", "{option_id:[0-9]+}", "value")
	specificOptionValueEndpoint := buildRoute("product_option_values", "{option_value_id:[0-9]+}")
	router.HandleFunc(optionValueEndpoint, buildProductOptionValueCreationHandler(db)).Methods("POST")
	router.HandleFunc(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db)).Methods("PUT")

	// Discounts
	specificDiscountEndpoint := buildRoute("discount", "{discount_id:[0-9]+}")
	router.HandleFunc(specificDiscountEndpoint, buildDiscountRetrievalHandler(db)).Methods("GET")
	router.HandleFunc(buildRoute("discounts"), buildDiscountListRetrievalHandler(db)).Methods("GET")
	router.HandleFunc(buildRoute("discount"), buildDiscountCreationHandler(db)).Methods("POST")
	// specificDiscountCodeEndpoint := buildRoute("discount", fmt.Sprintf("{code:%s}", SKUPattern))
	// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods("HEAD")
}
