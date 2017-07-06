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
	// NumericPattern repesents numeric values
	NumericPattern = `[0-9]+`
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
	router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserInfoUpdateHandler(db))
	router.Post("/password_reset", buildUserForgottenPasswordHandler(db))
	router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(db))
	//router.Head("/password_reset/{reset_token:[a-zA-Z0-9]{}}", buildUserPasswordResetTokenValidationHandler(db))

	router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(db))

		// Products
		productEndpoint := fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern)
		r.Post("/product", buildProductCreationHandler(db))
		r.Get("/products", buildProductListHandler(db))
		r.Get(productEndpoint, buildSingleProductHandler(db))
		r.Patch(productEndpoint, buildProductUpdateHandler(db))
		r.Head(productEndpoint, buildProductExistenceHandler(db))
		r.Delete(productEndpoint, buildProductDeletionHandler(db))

		// Product Options
		productOptionEndpoint := fmt.Sprintf("/product/{product_id:%s}/options", NumericPattern)
		specificOptionEndpoint := fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern)
		r.Get(productOptionEndpoint, buildProductOptionListHandler(db))
		r.Post(productOptionEndpoint, buildProductOptionCreationHandler(db))
		r.Patch(specificOptionEndpoint, buildProductOptionUpdateHandler(db))

		// Product Option Values
		optionValueEndpoint := fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern)
		specificOptionValueEndpoint := fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern)
		r.Post(optionValueEndpoint, buildProductOptionValueCreationHandler(db))
		r.Patch(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db))
		r.Delete(specificOptionValueEndpoint, buildProductOptionValueDeletionHandler(db))

		// Discounts
		specificDiscountEndpoint := fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern)
		r.Get(specificDiscountEndpoint, buildDiscountRetrievalHandler(db))
		r.Patch(specificDiscountEndpoint, buildDiscountUpdateHandler(db))
		r.Delete(specificDiscountEndpoint, buildDiscountDeletionHandler(db))
		r.Get("/discounts", buildDiscountListRetrievalHandler(db))
		r.Post("/discount", buildDiscountCreationHandler(db))
		// specificDiscountCodeEndpoint := buildRoute("v1", "discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
		// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
	})
}
