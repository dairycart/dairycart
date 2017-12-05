package main

import (
	"fmt"
	"strings"

	// external dependencies
	"github.com/go-chi/chi"
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
func SetupAPIRoutes(config *ServerConfig) {
	// Auth
	config.Router.Post("/login", buildUserLoginHandler(config.DB, config.Dairyclient, config.CookieStore))
	config.Router.Post("/logout", buildUserLogoutHandler(config.CookieStore))
	config.Router.Post("/user", buildUserCreationHandler(config.DB, config.Dairyclient, config.CookieStore))
	config.Router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserInfoUpdateHandler(config.DB, config.Dairyclient))
	config.Router.Post("/password_reset", buildUserForgottenPasswordHandler(config.DB, config.Dairyclient))
	config.Router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(config.DB, config.Dairyclient))

	config.Router.Route("/v1", func(r chi.Router) {
		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(config.DB, config.Dairyclient, config.CookieStore))

		// Product Roots
		specificProductRootRoute := fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern)
		r.Get("/product_roots", buildProductRootListHandler(config.DB, config.Dairyclient))
		r.Get(specificProductRootRoute, buildSingleProductRootHandler(config.DB, config.Dairyclient))
		r.Delete(specificProductRootRoute, buildProductRootDeletionHandler(config.DB, config.Dairyclient))

		// Products
		specificProductRoute := fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern)
		r.Get("/products", buildProductListHandler(config.DB, config.Dairyclient))
		r.Post("/product", buildProductCreationHandler(config.DB, config.Dairyclient))
		r.Get(specificProductRoute, buildSingleProductHandler(config.DB, config.Dairyclient))
		r.Patch(specificProductRoute, buildProductUpdateHandler(config.DB, config.Dairyclient, config.WebhookExecutor))
		r.Head(specificProductRoute, buildProductExistenceHandler(config.DB, config.Dairyclient))
		r.Delete(specificProductRoute, buildProductDeletionHandler(config.DB, config.Dairyclient))

		// Product Options
		optionsListRoute := fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern)
		specificOptionRoute := fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern)
		r.Get(optionsListRoute, buildProductOptionListHandler(config.DB, config.Dairyclient))
		r.Post(optionsListRoute, buildProductOptionCreationHandler(config.DB, config.Dairyclient))
		r.Patch(specificOptionRoute, buildProductOptionUpdateHandler(config.DB, config.Dairyclient))
		r.Delete(specificOptionRoute, buildProductOptionDeletionHandler(config.DB, config.Dairyclient))

		// Product Option Values
		specificOptionValueRoute := fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern)
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(config.DB, config.Dairyclient))
		r.Patch(specificOptionValueRoute, buildProductOptionValueUpdateHandler(config.DB, config.Dairyclient))
		r.Delete(specificOptionValueRoute, buildProductOptionValueDeletionHandler(config.DB, config.Dairyclient))

		// Discounts
		specificDiscountRoute := fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern)
		r.Get(specificDiscountRoute, buildDiscountRetrievalHandler(config.DB, config.Dairyclient))
		r.Patch(specificDiscountRoute, buildDiscountUpdateHandler(config.DB, config.Dairyclient))
		r.Delete(specificDiscountRoute, buildDiscountDeletionHandler(config.DB, config.Dairyclient))
		r.Get("/discounts", buildDiscountListRetrievalHandler(config.DB, config.Dairyclient))
		r.Post("/discount", buildDiscountCreationHandler(config.DB, config.Dairyclient))
	})
}
