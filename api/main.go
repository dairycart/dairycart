// +build !test

package main

import (
	"database/sql"
	"fmt"
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

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

const (
	DefaultPort                 = 4321
	DefaultImageStorageProvider = "local"
	DefaultDatabaseProvider     = "postgres"
)

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

func convertDairyConfigToRouterConfig(in DairyConfig) (*ServerConfig, error) {
	db, err := sql.Open(strings.ToLower(in.Database.Name), in.Database.ConnectionString)
	if err != nil {
		logrus.Fatalf("error encountered connecting to database: %v", err)
	}

	config := &ServerConfig{
		Router:          chi.NewRouter(),
		DB:              db,
		CookieStore:     setupCookieStorage(in.Secret),
		WebhookExecutor: &webhookExecutor{Client: http.DefaultClient},
	}

	if strings.ToLower(in.Database.Name) == DefaultDatabaseProvider && in.Database.PluginPath == "" {
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

	if strings.ToLower(in.ImageStorage.Name) == DefaultImageStorageProvider && in.ImageStorage.PluginPath == "" {
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

func fileServer(r chi.Router, path string, root http.FileSystem) {
	// path := fmt.Sprintf("/%s/", local.LocalProductImagesDirectory)
	// root := http.Dir(local.LocalProductImagesDirectory)

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

	setConfigDefaults()
	cfg := loadConfigFile()
	config, err := convertDairyConfigToRouterConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	config.Router.Use(middleware.RequestID)
	config.Router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
	SetupAPIRouter(config)

	fileServer(config.Router, fmt.Sprintf("/%s/", local.LocalProductImagesDirectory), http.Dir(local.LocalProductImagesDirectory))

	http.Handle("/", context.ClearHandler(config.Router))
	log.Printf("API now listening for requests on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
