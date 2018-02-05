package api

import (
	"fmt"
	"io"
	"net/http"
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

// SetupAPIRouter takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRouter(config *ServerConfig) {
	// health check
	config.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "healthy!") })

	// Auth
	config.Router.Post("/login", buildUserLoginHandler(config.DB, config.DatabaseClient, config.CookieStore))
	config.Router.Post("/logout", buildUserLogoutHandler(config.CookieStore))
	config.Router.Post("/user", buildUserCreationHandler(config.DB, config.DatabaseClient, config.CookieStore))
	config.Router.Patch(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserUpdateHandler(config.DB, config.DatabaseClient))
	config.Router.Post("/password_reset", buildUserForgottenPasswordHandler(config.DB, config.DatabaseClient))
	config.Router.Head("/password_reset/{reset_token}", buildUserPasswordResetTokenValidationHandler(config.DB, config.DatabaseClient))

	config.Router.Route("/v1", func(r chi.Router) {
		// r.Use(middleware.AllowContentType("application/json"))

		// Users
		r.Delete(fmt.Sprintf("/user/{user_id:%s}", NumericPattern), buildUserDeletionHandler(config.DB, config.DatabaseClient, config.CookieStore))

		// Product Roots
		specificProductRootRoute := fmt.Sprintf("/product_root/{product_root_id:%s}", NumericPattern)
		r.Get("/product_roots", buildProductRootListHandler(config.DB, config.DatabaseClient))
		r.Get(specificProductRootRoute, buildSingleProductRootHandler(config.DB, config.DatabaseClient))
		r.Delete(specificProductRootRoute, buildProductRootDeletionHandler(config.DB, config.DatabaseClient))

		// Products
		specificProductRoute := fmt.Sprintf("/product/{sku:%s}", ValidURLCharactersPattern)
		r.Get("/products", buildProductListHandler(config.DB, config.DatabaseClient))
		r.Post("/product", buildProductCreationHandler(config.DB, config.DatabaseClient, config.ImageStorer, config.WebhookExecutor))
		r.Get(specificProductRoute, buildSingleProductHandler(config.DB, config.DatabaseClient))
		r.Patch(specificProductRoute, buildProductUpdateHandler(config.DB, config.DatabaseClient, config.WebhookExecutor))
		r.Head(specificProductRoute, buildProductExistenceHandler(config.DB, config.DatabaseClient))
		r.Delete(specificProductRoute, buildProductDeletionHandler(config.DB, config.DatabaseClient, config.WebhookExecutor))

		// Product Options
		optionsListRoute := fmt.Sprintf("/product/{product_root_id:%s}/options", NumericPattern)
		specificOptionRoute := fmt.Sprintf("/product_options/{option_id:%s}", NumericPattern)
		r.Get(optionsListRoute, buildProductOptionListHandler(config.DB, config.DatabaseClient))
		r.Post(optionsListRoute, buildProductOptionCreationHandler(config.DB, config.DatabaseClient))
		r.Patch(specificOptionRoute, buildProductOptionUpdateHandler(config.DB, config.DatabaseClient))
		r.Delete(specificOptionRoute, buildProductOptionDeletionHandler(config.DB, config.DatabaseClient))

		// Product Option Values
		specificOptionValueRoute := fmt.Sprintf("/product_option_values/{option_value_id:%s}", NumericPattern)
		// r.Get(fmt.Sprintf("/product_options/{option_id:%s}/values", NumericPattern), buildProductOptionValueListRetrievalHandler(config.DB, config.DatabaseClient))
		r.Post(fmt.Sprintf("/product_options/{option_id:%s}/value", NumericPattern), buildProductOptionValueCreationHandler(config.DB, config.DatabaseClient))
		r.Patch(specificOptionValueRoute, buildProductOptionValueUpdateHandler(config.DB, config.DatabaseClient))
		r.Delete(specificOptionValueRoute, buildProductOptionValueDeletionHandler(config.DB, config.DatabaseClient))

		// Discounts
		specificDiscountRoute := fmt.Sprintf("/discount/{discount_id:%s}", NumericPattern)
		r.Get(specificDiscountRoute, buildDiscountRetrievalHandler(config.DB, config.DatabaseClient))
		r.Patch(specificDiscountRoute, buildDiscountUpdateHandler(config.DB, config.DatabaseClient))
		r.Delete(specificDiscountRoute, buildDiscountDeletionHandler(config.DB, config.DatabaseClient))
		r.Get("/discounts", buildDiscountListRetrievalHandler(config.DB, config.DatabaseClient))
		r.Post("/discount", buildDiscountCreationHandler(config.DB, config.DatabaseClient))

		// Webhooks
		specificWebhookRoute := fmt.Sprintf("/webhook/{webhook_id:%s}", NumericPattern)
		r.Get(fmt.Sprintf("/webhooks/{event_type:%s}", ValidURLCharactersPattern), buildWebhookListRetrievalByEventTypeHandler(config.DB, config.DatabaseClient))
		r.Get("/webhooks", buildWebhookListRetrievalHandler(config.DB, config.DatabaseClient))
		r.Post("/webhook", buildWebhookCreationHandler(config.DB, config.DatabaseClient))
		r.Patch(specificWebhookRoute, buildWebhookUpdateHandler(config.DB, config.DatabaseClient))
		r.Delete(specificWebhookRoute, buildWebhookDeletionHandler(config.DB, config.DatabaseClient))
	})
}
