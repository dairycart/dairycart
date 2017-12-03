package main

import (
	"database/sql"
	"fmt"
	"strings"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage"

	// external dependencies
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
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
func SetupAPIRoutes(router *chi.Mux, db *sql.DB, cookies *sessions.CookieStore, client storage.Storer) {
	// Auth
	router.Post("/login", buildUserLoginHandler(db, client, cookies))
	router.Post("/logout", buildUserLogoutHandler(cookies))
	router.Post("/user", buildUserCreationHandler(db, client, cookies))
	router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserInfoUpdateHandler(db, client))
	router.Post("/password_reset", buildUserForgottenPasswordHandler(db, client))
	router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(db, client))

	router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(db, client, cookies))

		// Product Roots
		specificProductRootRoute := fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern)
		r.Get("/product_roots", buildProductRootListHandler(db, client))
		r.Get(specificProductRootRoute, buildSingleProductRootHandler(db, client))
		r.Delete(specificProductRootRoute, buildProductRootDeletionHandler(db, client))

		// Products
		specificProductRoute := fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern)
		r.Get("/products", buildProductListHandler(db, client))
		r.Post("/product", buildProductCreationHandler(db, client))
		r.Get(specificProductRoute, buildSingleProductHandler(db, client))
		r.Patch(specificProductRoute, buildProductUpdateHandler(db, client))
		r.Head(specificProductRoute, buildProductExistenceHandler(db, client))
		r.Delete(specificProductRoute, buildProductDeletionHandler(db, client))

		// Product Options
		optionsListRoute := fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern)
		specificOptionRoute := fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern)
		r.Get(optionsListRoute, buildProductOptionListHandler(db, client))
		r.Post(optionsListRoute, buildProductOptionCreationHandler(db, client))
		r.Patch(specificOptionRoute, buildProductOptionUpdateHandler(db, client))
		r.Delete(specificOptionRoute, buildProductOptionDeletionHandler(db, client))

		// Product Option Values
		specificOptionValueRoute := fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern)
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(db, client))
		r.Patch(specificOptionValueRoute, buildProductOptionValueUpdateHandler(db, client))
		r.Delete(specificOptionValueRoute, buildProductOptionValueDeletionHandler(db, client))

		// Discounts
		specificDiscountRoute := fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern)
		r.Get(specificDiscountRoute, buildDiscountRetrievalHandler(db, client))
		r.Patch(specificDiscountRoute, buildDiscountUpdateHandler(db, client))
		r.Delete(specificDiscountRoute, buildDiscountDeletionHandler(db, client))
		r.Get("/discounts", buildDiscountListRetrievalHandler(db, client))
		r.Post("/discount", buildDiscountCreationHandler(db, client))
	})
}
