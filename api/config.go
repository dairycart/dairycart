package main

import (
	"database/sql"
	"github.com/dairycart/dairycart/storage/images/local"
	"log"
	"net/http"
	"os"
	"plugin"
	"strings"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairycart/storage/images"
	"github.com/dairycart/postgres"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	defaultPort                 = 4321
	DefaultImageStorageProvider = "local"
	DefaultDatabaseProvider     = "postgres"

	// Config keys
	// ===========
	// databases
	databaseKey       = "database"
	databaseTypeKey   = "database.type"
	databasePluginKey = "database.plugin_path"
	connectionKey     = "database.connection_details"

	// image storage
	imageStorageKey       = "image_storage"
	imageStorageTypeKey   = "image_storage.type"
	imageStoragePluginKey = "image_storage.plugin_path"

	// defaults
	defaultDatabaseType     = "postgres"
	defaultImageStorageType = "local"
)

type ServerConfig struct {
	Router          *chi.Mux
	DB              *sql.DB
	CookieStore     *sessions.CookieStore
	DatabaseClient  database.Storer
	WebhookExecutor WebhookExecutor
	ImageStorer     images.ImageStorer
}

func loadPlugin(pluginPath string, symbolName string) (plugin.Symbol, error) {
	if pluginPath == "" {
		return nil, errors.New("plugin path cannot be empty")
	}
	if symbolName == "" {
		return nil, errors.New("symbol name cannot be empty")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open plugin")
	}

	symbolToLookup := symbolName
	if symbolName[:1] == strings.ToLower(symbolName[:1]) {
		symbolToLookup = strings.Title(symbolToLookup)
	}

	sym, err := p.Lookup(symbolToLookup)
	if err != nil {
		return nil, errors.Wrap(err, "failed to locate appropriate plugin symbol")
	}
	return sym, nil
}

func setupCookieStorage(secret string) *sessions.CookieStore {
	if len(secret) < 32 {
		log.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	return sessions.NewCookieStore([]byte(secret))
}

func setConfigDefaults(config *viper.Viper) {
	config.SetConfigName("dairyconfig")
	if len(os.Args) >= 2 {
		config.SetConfigName(os.Args[1])
	}
	config.AddConfigPath(".")

	config.SetDefault("port", defaultPort)
	config.SetDefault("domain", "http://localhost")

	// config.SetDefault("database", DatabaseConfig{PluginConfig: PluginConfig{Name: "postgres"}})
	// config.SetDefault("imagestorage", ImageStorageConfig{PluginConfig: PluginConfig{Name: "local"}})
}

func validateServerConfig(config *viper.Viper) error {
	setConfigDefaults(config)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

func buildServerConfig(cfg *viper.Viper) (*ServerConfig, error) {
	db, dbClient, err := buildDatabaseFromConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error configuring database")
	}

	imageStorer, err := buildImageStorerFromConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error configuring image storage")
	}

	return &ServerConfig{
		Router:          chi.NewMux(),
		DB:              db,
		CookieStore:     setupCookieStorage(cfg.GetString("secret")),
		DatabaseClient:  dbClient,
		WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
		ImageStorer:     imageStorer,
	}, nil
}

func buildDatabaseFromConfig(cfg *viper.Viper) (*sql.DB, database.Storer, error) {
	var (
		db     *sql.DB
		client database.Storer
		err    error
	)

	if !cfg.IsSet(databaseKey) || cfg.GetString(connectionKey) == "" {
		return nil, nil, errors.New("no database type specified in config")
	}

	dbType := cfg.GetString(databaseTypeKey)
	connectionStr := cfg.GetString(connectionKey)
	db, err = sql.Open(dbType, connectionStr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "issue opening sql connection")
	}

	if dbType == defaultDatabaseType {
		client = postgres.NewPostgres()
	} else {
		missingPluginErr := errors.New("non-default database selected without complimentary plugin path, please check your configuration file")
		if !cfg.IsSet(databasePluginKey) {
			return nil, nil, missingPluginErr
		}

		pluginPath := cfg.GetString(databasePluginKey)
		if pluginPath == "" {
			return nil, nil, missingPluginErr
		}

		client, err = loadDatabasePlugin(pluginPath, dbType)
	}

	return db, client, err
}

func loadDatabasePlugin(pluginPath string, name string) (database.Storer, error) {
	dbSym, err := loadPlugin(pluginPath, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load plugin")
	}
	if _, ok := dbSym.(database.Storer); !ok {
		return nil, errors.New("Symbol provided in database plugin does not satisfy the database.Storer interface")
	}

	return dbSym.(database.Storer), nil
}

func buildImageStorerFromConfig(cfg *viper.Viper) (images.ImageStorer, error) {
	var (
		storer images.ImageStorer
		err    error
	)

	if !cfg.IsSet(databaseKey) || cfg.GetString(connectionKey) == "" {
		return nil, errors.New("no database type specified in config")
	}

	storageType := cfg.GetString(imageStorageTypeKey)
	if storageType == defaultImageStorageType {
		storer = local.NewLocalImageStorer()
	} else {
		missingPluginErr := errors.New("non-default database selected without complimentary plugin path, please check your configuration file")
		if !cfg.IsSet(imageStoragePluginKey) {
			return nil, missingPluginErr
		}

		pluginPath := cfg.GetString(imageStoragePluginKey)
		if pluginPath == "" {
			return nil, missingPluginErr
		}

		storer, err = loadImageStoragePlugin(pluginPath, storageType)
	}

	return storer, err
}

func loadImageStoragePlugin(pluginPath string, name string) (images.ImageStorer, error) {
	imgSym, err := loadPlugin(pluginPath, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load plugin")
	}
	if _, ok := imgSym.(images.ImageStorer); !ok {
		return nil, errors.New("Symbol provided in database plugin does not satisfy the database.Storer interface")
	}

	return imgSym.(images.ImageStorer), nil
}