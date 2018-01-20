// +build !test

package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"plugin"
	"strings"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairycart/storage/images"
	"github.com/dairycart/dairycart/storage/images/local"
	"github.com/dairycart/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

const (
	DefaultPort                 = 4321
	DefaultPhotoDir             = "product_images"
	DefaultImageStorageProvider = "local"
	DefaultDatabaseProvider     = "postgres"
)

type ImageStorageConfig struct {
	PluginConfig
	Domain string `json:"domain,omitempty"`
}

type DairyConfig struct {
	Secret       string             `json:"-"`
	Port         uint16             `json:"port,omitempty"`
	Database     PluginConfig       `json:"database,omitempty"`
	ImageStorage ImageStorageConfig `json:"image_storage,omitempty"`
}

func loadPlugin(pluginPath string, symbolName string) (plugin.Symbol, error) {
	if pluginPath == "" {
		return nil, errors.New("plugin path may not be empty")
	}
	if symbolName == "" {
		return nil, errors.New("symbol name may not be empty")
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

func convertDairyConfigToRouterConfig(in DairyConfig) (*RouterConfig, error) {
	config := &RouterConfig{
		Router:          chi.NewRouter(),
		CookieStore:     setupCookieStorage(in.Secret),
		WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
	}

	if strings.ToLower(in.Database.Name) == "postgres" && in.Database.PluginPath == "" {
		config.DatabaseClient = postgres.NewPostgres()
	} else if in.Database.PluginPath != "" && in.Database.Name != "" {
		dbSym, err := loadPlugin(in.Database.PluginPath, in.Database.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load plugin")
		}
		if _, ok := dbSym.(database.Storer); !ok {
			return nil, errors.New("Symbol provided in database plugin does not satisfy the database.Storer interface")
		}

		config.DatabaseClient = dbSym.(database.Storer)
	}

	if strings.ToLower(in.ImageStorage.Name) == "local" && in.ImageStorage.PluginPath == "" {
		config.ImageStorer = &local.LocalImageStorer{BaseURL: "http://localhost:4321"}
	} else if in.ImageStorage.PluginPath != "" && in.ImageStorage.Name != "" {
		imgSym, err := loadPlugin(in.Database.PluginPath, in.Database.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load plugin")
		}
		if _, ok := imgSym.(images.ImageStorer); !ok {
			return nil, errors.New("Symbol provided in database plugin does not satisfy the database.Storer interface")
		}

		config.ImageStorer = imgSym.(images.ImageStorer)
	}

	return config, nil
}

func validateServerConfig() DairyConfig {
	viper.SetConfigName("dairyconfig")
	if len(os.Args) >= 2 {
		viper.SetConfigName(os.Args[1])
	}
	viper.AddConfigPath(".")

	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("domain", "http://localhost")

	viper.SetDefault("database", PluginConfig{Name: DefaultDatabaseProvider})
	viper.SetDefault("imagestorage", ImageStorageConfig{
		PluginConfig: PluginConfig{
			Name: DefaultImageStorageProvider,
		},
		Domain: fmt.Sprintf("http://localhost:%d", DefaultPort),
	})
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}

	var config DairyConfig
	viper.Unmarshal(&config)

	return config
}

func buildServerConfig() *RouterConfig {
	// Connect to the database
	dbChoice := strings.ToLower(os.Getenv("DB_TO_USE"))
	switch dbChoice {
	case "postgres":
		dbURL := os.Getenv("DAIRYCART_DB_URL")
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			logrus.Fatalf("error encountered connecting to database: %v", err)
		}

		loadExampleData := os.Getenv("MIGRATE_EXAMPLE_DATA") == "YES"
		pg := postgres.NewPostgres()
		if err = pg.Migrate(db, dbURL, loadExampleData); err != nil {
			logrus.Fatalf("error encountered migrating database: %v", err)
		}

		return &RouterConfig{
			DB:              db,
			DatabaseClient:  postgres.NewPostgres(),
			WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
			ImageStorer:     &local.LocalImageStorer{BaseURL: "http://localhost:4321"},
		}
	default:
		logrus.Fatalf("invalid database choice: '%s'", dbChoice)
	}
	return nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	cfg := validateServerConfig()

	config := buildDefaultConfig()
	config.Router.Use(middleware.RequestID)
	config.Router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
	SetupAPIRouter(config)

	fileServer(config.Router, fmt.Sprintf("/%s/", DefaultPhotoDir), http.Dir(DefaultPhotoDir))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "healthy!") })
	http.Handle("/", context.ClearHandler(config.Router))
	log.Printf("API now listening for requests on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
