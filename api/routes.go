package main

import (
	"fmt"
	"strings"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage"

	// external dependencies
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
func SetupAPIRoutes(router *chi.Mux, dbx *sqlx.DB, store *sessions.CookieStore, db storage.Storage) {
	// Auth
	router.Post("/login", buildUserLoginHandler(dbx, store))
	router.Post("/logout", buildUserLogoutHandler(store))
	router.Post("/user", buildUserCreationHandler(dbx, store))
	router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserInfoUpdateHandler(dbx))
	router.Post("/password_reset", buildUserForgottenPasswordHandler(dbx))
	router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(dbx))
	//router.Head("/password_reset/{reset_token:[a-zA-Z0-9]{}}", buildUserPasswordResetTokenValidationHandler(dbx))

	router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(dbx, store))

		// Product Roots
		r.Get("/product_roots", buildProductRootListHandler(dbx))
		r.Get(fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern), buildSingleProductRootHandler(dbx))
		r.Delete(fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern), buildProductRootDeletionHandler(dbx))

		// Products
		r.Post("/product", buildProductCreationHandler(dbx))
		r.Get("/products", buildProductListHandler(dbx))
		r.Get(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildSingleProductHandler(db))
		r.Patch(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductUpdateHandler(dbx))
		r.Head(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductExistenceHandler(dbx))
		r.Delete(fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern), buildProductDeletionHandler(dbx))

		// Product Options
		r.Get(fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern), buildProductOptionListHandler(dbx))
		r.Post(fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern), buildProductOptionCreationHandler(dbx))
		r.Patch(fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern), buildProductOptionUpdateHandler(dbx))
		r.Delete(fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern), buildProductOptionDeletionHandler(dbx))

		// Product Option Values
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(dbx))
		r.Patch(fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern), buildProductOptionValueUpdateHandler(dbx))
		r.Delete(fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern), buildProductOptionValueDeletionHandler(dbx))

		// Discounts
		r.Get(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountRetrievalHandler(dbx))
		r.Patch(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountUpdateHandler(dbx))
		r.Delete(fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern), buildDiscountDeletionHandler(dbx))
		r.Get("/discounts", buildDiscountListRetrievalHandler(dbx))
		r.Post("/discount", buildDiscountCreationHandler(dbx))
		// specificDiscountCodeEndpoint := buildRoute("v1", "discount", fmt.Sprintf("{code:%s}", ValidURLCharactersPattern))
		// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(dbx)).Methods(http.MethodHead)
	})
}
