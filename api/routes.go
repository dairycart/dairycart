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
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(db, store))

		// Product Roots
		r.Get("/product_roots", buildProductRootListHandler(db))
		r.Get(fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern), buildSingleProductRootHandler(db))
		r.Delete(fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern), buildProductRootDeletionHandler(db))

		// Products
		r.Post("/product", buildProductCreationHandler(db))
		r.Get("/products", buildProductListHandler(db))
		r.Get(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildSingleProductHandler(db))
		r.Patch(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductUpdateHandler(db))
		r.Head(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductExistenceHandler(db))
		r.Delete(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductDeletionHandler(db))

		// Product Options
		r.Get(fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern), buildProductOptionListHandler(db))
		r.Post(fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern), buildProductOptionCreationHandler(db))
		r.Patch(fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern), buildProductOptionUpdateHandler(db))
		r.Delete(fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern), buildProductOptionDeletionHandler(db))

		// Product Option Values
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(db))
		r.Patch(fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern), buildProductOptionValueUpdateHandler(db))
		r.Delete(fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern), buildProductOptionValueDeletionHandler(db))

		// Discounts
		r.Get(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountRetrievalHandler(db))
		r.Patch(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountUpdateHandler(db))
		r.Delete(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountDeletionHandler(db))
		r.Get("/discounts", buildDiscountListRetrievalHandler(db))
		r.Post("/discount", buildDiscountCreationHandler(db))
		// specificDiscountCodeEndpoint := buildRoute("v1", "discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
		// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
	})
}
