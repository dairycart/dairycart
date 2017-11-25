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
func SetupAPIRoutes(router *chi.Mux, db *sql.DB, dbxReplaceMePlz *sqlx.DB, store *sessions.CookieStore, client storage.Storer) {
	// Auth
	router.Post("/login", buildUserLoginHandler(db, client, store))
	router.Post("/logout", buildUserLogoutHandler(store))
	router.Post("/user", buildUserCreationHandler(db, client, store))
	router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserInfoUpdateHandler(db, client))
	router.Post("/password_reset", buildUserForgottenPasswordHandler(db, client))
	router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(db, client))

	router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(db, client, store))

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
		r.Get(optionsListRoute, buildProductOptionListHandler(dbxReplaceMePlz))
		r.Post(optionsListRoute, buildProductOptionCreationHandler(dbxReplaceMePlz, client))
		r.Patch(specificOptionRoute, buildProductOptionUpdateHandler(dbxReplaceMePlz))
		r.Delete(specificOptionRoute, buildProductOptionDeletionHandler(dbxReplaceMePlz))

		// Product Option Values
		specificOptionValueRoute := fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern)
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(dbxReplaceMePlz))
		r.Patch(specificOptionValueRoute, buildProductOptionValueUpdateHandler(dbxReplaceMePlz))
		r.Delete(specificOptionValueRoute, buildProductOptionValueDeletionHandler(dbxReplaceMePlz))

		// Discounts
		specificDiscountRoute := fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern)
		r.Get(specificDiscountRoute, buildDiscountRetrievalHandler(dbxReplaceMePlz))
		r.Patch(specificDiscountRoute, buildDiscountUpdateHandler(dbxReplaceMePlz))
		r.Delete(specificDiscountRoute, buildDiscountDeletionHandler(dbxReplaceMePlz))
		r.Get("/discounts", buildDiscountListRetrievalHandler(dbxReplaceMePlz))
		r.Post("/discount", buildDiscountCreationHandler(dbxReplaceMePlz))
	})
}
