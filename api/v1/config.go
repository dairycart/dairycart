package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"plugin"
	"strings"

	"github.com/dairycart/dairycart/storage/v1/database"
	"github.com/dairycart/dairycart/storage/v1/database/postgres"
	"github.com/dairycart/dairycart/storage/v1/images"
	"github.com/dairycart/dairycart/storage/v1/images/local"

	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	mandatorySecretLength       = 32
	defaultPort                 = 4321
	DefaultImageStorageProvider = "local"
	DefaultDatabaseProvider     = "postgres"

	// Config keys //
	// =========== //

	// basic stuff
	portKey   = "port"
	domainKey = "domain"
	secretKey = "secret"

	// databases
	databaseKey           = "database"
	migrateExampleDataKey = "database.migrate_example_data"
	databaseTypeKey       = "database.type"
	databasePluginKey     = "database.plugin_path"
	databaseConnectionKey = "database.connection_details"

	// image storage
	imageStorageKey       = "image_storage"
	imageStorageTypeKey   = "image_storage.type"
	imageStoragePluginKey = "image_storage.plugin_path"
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

	if symbolName[:1] == strings.ToLower(symbolName[:1]) {
		symbolName = strings.Title(symbolName)
	}

	sym, err := p.Lookup(symbolName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to locate appropriate plugin symbol")
	}
	return sym, nil
}

func setupCookieStorage(secret string) (*sessions.CookieStore, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("Something is up with your app secret: `%s`", secret)
	}
	return sessions.NewCookieStore([]byte(secret)), nil
}

func setConfigDefaults(config *viper.Viper, altCfgFilePath string) {
	config.SetConfigName("dairyconfig")
	if altCfgFilePath != "" {
		config.SetConfigFile(altCfgFilePath)
	}
	config.AddConfigPath(".")

	config.SetDefault(portKey, defaultPort)
	config.SetDefault(domainKey, fmt.Sprintf("http://localhost:%d", defaultPort))
	config.SetDefault(migrateExampleDataKey, false)

	config.SetDefault(databaseTypeKey, DefaultDatabaseProvider)
	config.SetDefault(imageStorageTypeKey, DefaultImageStorageProvider)

	// Secret stuff
	config.BindEnv(secretKey, "DAIRYSECRET")
	config.SetDefault(secretKey, uniuri.NewLen(mandatorySecretLength))
	// TODO: generate secret as default in case user doesn't provide one
}

func LoadServerConfig(altCfgFilePath string) (*viper.Viper, error) {
	config := viper.New()
	setConfigDefaults(config, altCfgFilePath)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "error validating server config")
		}
	}
	return config, nil
}

func BuildServerConfig(config *viper.Viper) (*ServerConfig, error) {
	db, dbClient, err := buildDatabaseFromConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "error configuring database")
	}

	imageStorer, err := buildImageStorerFromConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "error configuring image storage")
	}

	cookieStorer, err := setupCookieStorage(config.GetString(secretKey))
	if err != nil {
		return nil, errors.Wrap(err, "error configuring cookie storage")
	}

	return &ServerConfig{
		Router:          chi.NewMux(),
		DB:              db,
		CookieStore:     cookieStorer,
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

	dbConnStr := cfg.GetString(databaseConnectionKey)

	if !cfg.IsSet(databaseKey) || dbConnStr == "" {
		return nil, nil, errors.New("no database type specified in config")
	}

	dbType := cfg.GetString(databaseTypeKey)
	db, err = sql.Open(dbType, dbConnStr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "issue opening sql connection")
	}

	if dbType == DefaultDatabaseProvider {
		client = postgres.NewPostgres()
	} else {
		missingPluginErr := errors.New("non-default database selected without complimentary plugin path, please check your configuration file")
		pluginPath := cfg.GetString(databasePluginKey)

		if !cfg.IsSet(databasePluginKey) || pluginPath == "" {
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

	if !cfg.IsSet(imageStorageKey) {
		return nil, errors.New("no image storer type specified in config")
	}

	storageType := cfg.GetString(imageStorageTypeKey)
	if storageType == DefaultImageStorageProvider {
		storer = local.NewLocalImageStorer()
	} else {
		missingPluginErr := errors.New("non-default image storer selected without complimentary plugin path, please check your configuration file")
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
		return nil, errors.New("Symbol provided in image storage plugin does not satisfy the images.ImageStorer interface")
	}

	return imgSym.(images.ImageStorer), nil
}

// InitializeServerComponents calls Init on all the relevant server components, and migrates the database.
func InitializeServerComponents(cfg *viper.Viper, config *ServerConfig) error {
	var err error
	config.Router.Route("/product_images", func(r chi.Router) {
		imageConfig := cfg.Sub(imageStorageKey)
		err = config.ImageStorer.Init(imageConfig, r)
	})
	if err != nil {
		return errors.Wrap(err, "error migrating database")
	}

	dbConfig := cfg.Sub(databaseKey)
	err = config.DatabaseClient.Migrate(config.DB, dbConfig)
	if err != nil {
		return errors.Wrap(err, "error migrating database")
	}

	return nil
}
