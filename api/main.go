// +build !test

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dairycart/dairycart/storage/images/local"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

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

	cfg := viper.New()
	err := validateServerConfig(cfg)
	if err != nil {
		logrus.Fatalf("error validating server configuration: %v\n", err)
	}

	config, err := buildServerConfig(cfg)
	if err != nil {
		logrus.Fatalf("error configuring server: %v\n", err)
	}

	config.Router.Use(middleware.RequestID)
	config.Router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
	SetupAPIRouter(config)

	// config.Router.Route("/v1", func(r chi.Router) {
	// 	config.ImageStorer.Init(r)
	// })

	fileServer(config.Router, fmt.Sprintf("/%s/", local.LocalProductImagesDirectory), http.Dir(local.LocalProductImagesDirectory))
	http.Handle("/", context.ClearHandler(config.Router))

	port := cfg.GetInt("port")
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
