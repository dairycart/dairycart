package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairycart/storage/images"
	"github.com/dairycart/dairycart/storage/images/local"
	"github.com/dairycart/postgres"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

type PluginConfig struct {
	Name       string `json:"name,omitempty"`
	PluginPath string `json:"plugin_path,omitempty"`
}

type RouterConfig struct {
	Router          *chi.Mux
	DB              *sql.DB
	CookieStore     *sessions.CookieStore
	DatabaseClient  database.Storer
	WebhookExecutor WebhookExecutor
	ImageStorer     images.ImageStorer
}

func setupCookieStorage(secret string) *sessions.CookieStore {
	if len(secret) < 32 {
		log.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	return sessions.NewCookieStore([]byte(secret))
}

// buildDefaultConfig generates a RouterConfig object based on defaults for Dairycart.
// it should only be called in the even the user doesn't pass any valid config files
func buildDefaultConfig() *RouterConfig {
	secret := os.Getenv("DAIRYSECRET")

	return &RouterConfig{
		Router:          chi.NewRouter(),
		DatabaseClient:  postgres.NewPostgres(),
		CookieStore:     setupCookieStorage(secret),
		WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
		ImageStorer:     &local.LocalImageStorer{BaseURL: "http://localhost:4321"},
	}
}
