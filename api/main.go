// +build !test

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

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

	config.Router.Route("/product_images", func(r chi.Router) {
		config.ImageStorer.Init(cfg, r)
	})

	port := cfg.GetInt("port")
	http.Handle("/", context.ClearHandler(config.Router))
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
