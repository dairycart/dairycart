package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairycart/storage/images"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	PluginConfig
	ConnectionString string
}

type PluginConfig struct {
	Name       string            `json:"name,omitempty"`
	PluginPath string            `json:"plugin_path,omitempty"`
	Options    map[string]string `json:"extra_options,omitempty"`
}

type ImageStorageConfig struct {
	PluginConfig
	StoragePath string `json:"storage_path,omitempty"`
}

type DairyConfig struct {
	Secret       string         `json:"-"`
	Domain       string         `json:"domain,omitempty"`
	Port         uint16         `json:"port,omitempty"`
	Database     DatabaseConfig `json:"database,omitempty"`
	ImageStorage PluginConfig   `json:"image_storage,omitempty"`
}

type ServerConfig struct {
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

func setConfigDefaults() {
	viper.SetConfigName("dairyconfig")
	if len(os.Args) >= 2 {
		viper.SetConfigName(os.Args[1])
	}

	viper.AddConfigPath(".")
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("domain", "http://localhost")

	viper.SetDefault("database", DatabaseConfig{
		PluginConfig: PluginConfig{
			Name: DefaultDatabaseProvider,
			Options: map[string]string{
				"migrate_example_data": "no",
			},
		},
	})
	viper.SetDefault("imagestorage", ImageStorageConfig{
		PluginConfig: PluginConfig{
			Name: DefaultImageStorageProvider,
		},
	})
}

func loadConfigFile() DairyConfig {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}

	var config DairyConfig
	viper.Unmarshal(&config)

	return config
}
