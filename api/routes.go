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

func buildRoute(routeVersion string, routeParts ...string) string {
	return fmt.Sprintf("/%s/%s", routeVersion, strings.Join(routeParts, "/"))
}

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *chi.Mux, db *sqlx.DB, store *sessions.CookieStore) {
	// Auth
	router.Post("/login", buildUserLoginHandler(db, store))
	router.Post("/logout", buildUserLogoutHandler(store))
	router.Post("/user", buildUserCreationHandler(db, store))
	router.Put("/user/{user_id:[0-9]+}", buildUserInfoUpdateHandler(db))
	router.Post("/password_reset", buildUserForgottenPasswordHandler(db))
	router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(db))
	//router.Head("/password_reset/{reset_token:[a-zA-Z0-9]{}}", buildUserPasswordResetTokenValidationHandler(db))

	router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete("/user/{user_id:[0-9]+}", buildUserDeletionHandler(db))

		// Products
		productEndpoint := fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern)
		r.Post("/product", buildProductCreationHandler(db))
		r.Get("/products", buildProductListHandler(db))
		r.Get(productEndpoint, buildSingleProductHandler(db))
		r.Put(productEndpoint, buildProductUpdateHandler(db))
		r.Head(productEndpoint, buildProductExistenceHandler(db))
		r.Delete(productEndpoint, buildProductDeletionHandler(db))

		// Product Options
		productOptionEndpoint := "/product/{product_id:[0-9]+}/options"
		specificOptionEndpoint := "/product_options/{option_id:[0-9]+}"
		r.Get(productOptionEndpoint, buildProductOptionListHandler(db))
		r.Post(productOptionEndpoint, buildProductOptionCreationHandler(db))
		r.Put(specificOptionEndpoint, buildProductOptionUpdateHandler(db))

		// Product Option Values
		optionValueEndpoint := "/product_options/{option_id:[0-9]+}/value"
		specificOptionValueEndpoint := "/product_option_values/{option_value_id:[0-9]+}"
		r.Post(optionValueEndpoint, buildProductOptionValueCreationHandler(db))
		r.Put(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db))

		// Discounts
		specificDiscountEndpoint := "/discount/{discount_id:[0-9]+}"
		r.Get(specificDiscountEndpoint, buildDiscountRetrievalHandler(db))
		r.Put(specificDiscountEndpoint, buildDiscountUpdateHandler(db))
		r.Delete(specificDiscountEndpoint, buildDiscountDeletionHandler(db))
		r.Get("/discounts", buildDiscountListRetrievalHandler(db))
		r.Post("/discount", buildDiscountCreationHandler(db))
		// specificDiscountCodeEndpoint := buildRoute("v1", "discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
		// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
	})
}
