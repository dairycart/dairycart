package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
)

const (
	// SKUPattern represents the valid characters a sku can contain
	SKUPattern = `[a-zA-Z\-_]+`
)

var sqlBuilder squirrel.StatementBuilderType

func init() {
	sqlBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

//////////////////////////////////////////////////////////////////////////////////////
//                                                                                  //
//      _______________                                       ||*\_/*|________      //
//     |  ___________  |             .-.     .-.              ||_/-\_|______  |     //
//     | |           | |            .****. .****.             | |           | |     //
//     | |   0   0   | |            .*****.*****.             | |   0   0   | |     //
//     | |     -     | |             .*********.              | |     -     | |     //
//     | |   \___/   | |              .*******.               | |   \___/   | |     //
//     | |___     ___| |               .*****.                | |___________| |     //
//     |_____|\_/|_____|                .***.                 |_______________|     //
//       _|__|/ \|_|_.....................*......................_|________|_       //
//      / ********** \                   ^^^                    / ********** \      //
//    /  ************  \       (Dairycart API Traffic)         / ************* \    //
//   --------------------                                     -------------------   //
//////////////////////////////////////////////////////////////////////////////////////

func buildRoute(routeParts ...string) string {
	allRouteParts := append([]string{"v1"}, routeParts...)
	return fmt.Sprintf("/%s", strings.Join(allRouteParts, "/"))
}

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	productEndpoint := buildRoute("product", fmt.Sprintf("{sku:%s}", SKUPattern))
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods("POST")
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Attributes
	productAttributeEndpoint := buildRoute("product_attributes", "{progenitor_id:[0-9]+}")
	// specificAttributeEndpoint := buildRoute("product_attributes", "{progenitor_id:[0-9]+}", "{attribute_id:[0-9]+}")
	router.HandleFunc(productAttributeEndpoint, buildProductAttributeListHandler(db)).Methods("GET")
	router.HandleFunc(productAttributeEndpoint, buildProductAttributeCreationHandler(db)).Methods("POST")
	// router.HandleFunc(specificAttributeEndpoint, buildProductAttributeUpdateHandler(db)).Methods("PUT")

	// // Product Attribute Values
	attributeValueEndpoint := buildRoute("product_attributes", "{attribute_id:[0-9]+}", "value")
	attributeValuesEndpoint := buildRoute("product_attributes", "{attribute_id:[0-9]+}", "values")
	router.HandleFunc(attributeValueEndpoint, buildProductAttributeValueCreationHandler(db)).Methods("POST")
	router.HandleFunc(attributeValuesEndpoint, buildProductAttributeValueCreationHandler(db)).Methods("GET")
}
