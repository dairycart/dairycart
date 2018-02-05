// +build !test

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	dairyserver "github.com/dairycart/dairycart/api"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	cfg, err := dairyserver.LoadServerConfig()
	if err != nil {
		logrus.Fatalf("error validating server configuration: %v\n", err)
	}

	config, err := dairyserver.BuildServerConfig(cfg)
	if err != nil {
		logrus.Fatalf("error configuring server: %v\n", err)
	}

	config.Router.Use(middleware.RequestID)
	config.Router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
	dairyserver.SetupAPIRouter(config)

	err = dairyserver.InitializeServerComponents(cfg, config)
	if err != nil {
		logrus.Fatalf("error initializing server: %v\n", err)
	}

	port := cfg.GetInt("port")
	http.Handle("/", context.ClearHandler(config.Router))
	log.Printf("API now listening for requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
