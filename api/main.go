package main

import (
	// stdlib
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	// dependencies
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	log "github.com/sirupsen/logrus"

	// unnamed dependencies
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

const (
	// SKUPattern represents the valid characters a sku can contain
	SKUPattern = `[a-zA-Z\-_]+`
)

func buildRoute(routeParts ...string) string {
	allRouteParts := append([]string{"v1"}, routeParts...)
	return fmt.Sprintf("/%s", strings.Join(allRouteParts, "/"))
}

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sqlx.DB) {
	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", SKUPattern))
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods(http.MethodPut)
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods(http.MethodHead)
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods(http.MethodDelete)

	// Product Options
	productOptionEndpoint := buildRoute("product_options", "{progenitor_id:[0-9]+}")
	specificOptionEndpoint := buildRoute("product_options", "{option_id:[0-9]+}")
	router.HandleFunc(productOptionEndpoint, buildProductOptionListHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(productOptionEndpoint, buildProductOptionCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc(specificOptionEndpoint, buildProductOptionUpdateHandler(db)).Methods(http.MethodPut)

	// Product Option Values
	optionValueEndpoint := buildRoute("product_options", "{option_id:[0-9]+}", "value")
	specificOptionValueEndpoint := buildRoute("product_option_values", "{option_value_id:[0-9]+}")
	router.HandleFunc(optionValueEndpoint, buildProductOptionValueCreationHandler(db)).Methods(http.MethodPost)
	router.HandleFunc(specificOptionValueEndpoint, buildProductOptionValueUpdateHandler(db)).Methods(http.MethodPut)

	// Discounts
	specificDiscountEndpoint := buildRoute("discount", "{discount_id:[0-9]+}")
	router.HandleFunc(specificDiscountEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(specificDiscountEndpoint, buildDiscountUpdateHandler(db)).Methods(http.MethodPut)
	router.HandleFunc(specificDiscountEndpoint, buildDiscountDeletionHandler(db)).Methods(http.MethodDelete)
	router.HandleFunc(buildRoute("discounts"), buildDiscountListRetrievalHandler(db)).Methods(http.MethodGet)
	router.HandleFunc(buildRoute("discount"), buildDiscountCreationHandler(db)).Methods(http.MethodPost)
	// specificDiscountCodeEndpoint := buildRoute("discount", fmt.Sprintf("{code:%s}", SKUPattern))
	// router.HandleFunc(specificDiscountCodeEndpoint, buildDiscountRetrievalHandler(db)).Methods(http.MethodHead)
}

func determineMigrationCount() int {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatalf("missing migrations files")
	}

	migrationCount := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			migrationCount++
		}
	}
	return migrationCount
}

// this function not only waits for the database to accept its incoming connection, but also performs any necessary migrations
func migrateDatabase(db *sql.DB, migrationCount int) {
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
		} else {
			migrationsDir := "file://migrations" // os.Getenv("DAIRYCART_MIGRATIONS_DIR")
			m, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
			if err != nil {
				log.Fatalf("error encountered setting up new migration client: %v", err)
			}

			for i := 0; i < migrationCount; i++ {
				err = m.Steps(1)
				if err != nil {
					log.Printf("error encountered migrating database: %v", err)
					break
				}
			}
			databaseIsNotMigrated = false
		}
	}
}

func main() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	var err error
	pg, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error encountered connecting to database: %v", err)
	}

	// migrate the database
	migrationCount := determineMigrationCount()
	migrateDatabase(pg, migrationCount)

	// setup sqlx
	db := sqlx.NewDb(pg, "postgres")
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// setup all our API routes
	APIRouter := mux.NewRouter()
	// APIRouter.Host("api.dairycart.com")
	SetupAPIRoutes(APIRouter, db)

	// serve 'em up a lil' sauce
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "I live!") })
	http.Handle("/", APIRouter)
	log.Println("Dairycart now listening for requests")
	http.ListenAndServe(":80", nil)
}
