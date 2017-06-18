package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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
func SetupAPIRoutes(router *mux.Router, db *sqlx.DB) {
	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", SKUPattern))
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods(http.MethodPut)
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods(http.MethodHead)
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods(http.MethodDelete)

	// Product Options
	productOptionEndpoint := buildRoute("product_options", "{progenitor_id:[0-9]+}")
	specificOptionEndpoint := buildRoute("product_options", "{option_id:[0-9]+}")
	router.HandleFunc(productOptionEndpoint, buildProductOptionListHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productOptionEndpoint, buildProductOptionCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc(specificOptionEndpoint, buildProductOptionUpdateHandler(db)).Methods(http.MethodPut)

	// Product Option Values
	optionValueEndpoint := buildRoute("product_options", "{option_id:[0-9]+}", "value")
	specificOptionValueEndpoint := buildRoute("product_option_values", "{option_value_id:[0-9]+}")
	router.HandleFunc(optionValueEndpoint, buildProductOptionValueCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db)).Methods(http.MethodPut)

	// Discounts
	specificDiscountEndpoint := buildRoute("discount", "{discount_id:[0-9]+}")
	router.HandleFunc(specificDiscountEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(specificDiscountEndpoint, buildDiscountUpdateHandler(db)).Methods(http.MethodPut)
	router.HandleFunc(specificDiscountEndpoint, buildDiscountDeletionHandler(db)).Methods(http.MethodDelete)
	router.HandleFunc(buildRoute("discounts"), buildDiscountListRetrievalHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(buildRoute("discount"), buildDiscountCreationHandler(db)).Methods(http.MethodPost)
	// specificDiscountCodeEndpoint := buildRoute("discount", fmt.Sprintf("{code:%s}", SKUPattern))
	// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
}
