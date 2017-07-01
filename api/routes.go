package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	// ValidURLCharactersPattern represents the valid characters a sku can contain
	ValidURLCharactersPattern = `[a-zA-Z\-_]+`
)

func buildRoute(routeParts ...string) string {
	return fmt.Sprintf("/v1/%s", strings.Join(routeParts, "/"))
}

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sqlx.DB) {
	// Auth
	router.HandleFunc("/login", buildUserLoginHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/user", buildUserCreationHandler(db)).Methods(http.MethodPost)

	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", ValidURLCharactersPattern))
	router.HandleFunc("/v1/product", validateTokenMiddleware(buildProductCreationHandler(db))).Methods(http.MethodPost)
	router.HandleFunc("/v1/products", validateTokenMiddleware(buildProductListHandler(db))).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, validateTokenMiddleware(buildSingleProductHandler(db))).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, validateTokenMiddleware(buildProductUpdateHandler(db))).Methods(http.MethodPut)
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods(http.MethodHead)
	router.HandleFunc(productEndpoint, validateTokenMiddleware(buildProductDeletionHandler(db))).Methods(http.MethodDelete)

	// Product Options
	productOptionEndpoint := buildRoute("product", "{product_id:[0-9]+}", "options")
	specificOptionEndpoint := buildRoute("product_options", "{option_id:[0-9]+}")
	router.HandleFunc(productOptionEndpoint, validateTokenMiddleware(buildProductOptionListHandler(db))).Methods(http.MethodGet)
	router.HandleFunc(productOptionEndpoint, validateTokenMiddleware(buildProductOptionCreationHandler(db))).Methods(http.MethodPost)
	router.HandleFunc(specificOptionEndpoint, validateTokenMiddleware(buildProductOptionUpdateHandler(db))).Methods(http.MethodPut)

	// Product Option Values
	optionValueEndpoint := buildRoute("product_options", "{option_id:[0-9]+}", "value")
	specificOptionValueEndpoint := buildRoute("product_option_values", "{option_value_id:[0-9]+}")
	router.HandleFunc(optionValueEndpoint, validateTokenMiddleware(buildProductOptionValueCreationHandler(db))).Methods(http.MethodPost)
	router.HandleFunc(specificOptionValueEndpoint, validateTokenMiddleware(buildProductOptionValueUpdateHandler(db))).Methods(http.MethodPut)

	// Discounts
	specificDiscountEndpoint := buildRoute("discount", "{discount_id:[0-9]+}")
	router.HandleFunc(specificDiscountEndpoint, validateTokenMiddleware(buildDiscountRetrievalHandler(db))).Methods(http.MethodGet)
	router.HandleFunc(specificDiscountEndpoint, validateTokenMiddleware(buildDiscountUpdateHandler(db))).Methods(http.MethodPut)
	router.HandleFunc(specificDiscountEndpoint, validateTokenMiddleware(buildDiscountDeletionHandler(db))).Methods(http.MethodDelete)
	router.HandleFunc(buildRoute("discounts"), validateTokenMiddleware(buildDiscountListRetrievalHandler(db))).Methods(http.MethodGet)
	router.HandleFunc(buildRoute("discount"), validateTokenMiddleware(buildDiscountCreationHandler(db))).Methods(http.MethodPost)
	// specificDiscountCodeEndpoint := buildRoute("discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
	// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
}
