package main

import (
	"fmt"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
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
func SetupAPIRoutes(router *chi.Mux, db *sqlx.DB, store *sessions.CookieStore) {
	// Auth
	router.Post("/login", buildUserLoginHandler(db, store))
	router.Post("/user", buildUserCreationHandler(db, store))

	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", ValidURLCharactersPattern))
	router.Post("/v1/product", buildProductCreationHandler(db))
	router.Get("/v1/products", buildProductListHandler(db))
	router.Get(productEndpoint, buildSingleProductHandler(db))
	router.Put(productEndpoint, buildProductUpdateHandler(db))
	router.Head(productEndpoint, buildProductExistenceHandler(db))
	router.Delete(productEndpoint, buildProductDeletionHandler(db))

	// Product Options
	productOptionEndpoint := buildRoute("product", "{product_id:[0-9]+}", "options")
	specificOptionEndpoint := buildRoute("product_options", "{option_id:[0-9]+}")
	router.Get(productOptionEndpoint, buildProductOptionListHandler(db))
	router.Post(productOptionEndpoint, buildProductOptionCreationHandler(db))
	router.Put(specificOptionEndpoint, buildProductOptionUpdateHandler(db))

	// Product Option Values
	optionValueEndpoint := buildRoute("product_options", "{option_id:[0-9]+}", "value")
	specificOptionValueEndpoint := buildRoute("product_option_values", "{option_value_id:[0-9]+}")
	router.Post(optionValueEndpoint, buildProductOptionValueCreationHandler(db))
	router.Put(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db))

	// Discounts
	specificDiscountEndpoint := buildRoute("discount", "{discount_id:[0-9]+}")
	router.Get(specificDiscountEndpoint, buildDiscountRetrievalHandler(db))
	router.Put(specificDiscountEndpoint, buildDiscountUpdateHandler(db))
	router.Delete(specificDiscountEndpoint, buildDiscountDeletionHandler(db))
	router.Get(buildRoute("discounts"), buildDiscountListRetrievalHandler(db))
	router.Post(buildRoute("discount"), buildDiscountCreationHandler(db))
	// specificDiscountCodeEndpoint := buildRoute("discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
	// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
}
