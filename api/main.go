package api

import (
	"database/sql"
	"fmt"

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

// SetupAPIRoutes takes a mux router and a database connection and creates all the API routes for the API
func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	// Products
	productEndpoint := fmt.Sprintf("/v1/product/{sku:%s}", SKUPattern)
	router.HandleFunc("/v1/product", buildProductCreationHandler(db)).Methods("POST")
	router.HandleFunc("/v1/products", buildProductListHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildSingleProductHandler(db)).Methods("GET")
	router.HandleFunc(productEndpoint, buildProductUpdateHandler(db)).Methods("PUT")
	router.HandleFunc(productEndpoint, buildProductExistenceHandler(db)).Methods("HEAD")
	router.HandleFunc(productEndpoint, buildProductDeletionHandler(db)).Methods("DELETE")

	// Product Attributes
	router.HandleFunc("/product_attributes/{progenitor_id:[0-9]+}", buildProductAttributeListHandler(db)).Methods("GET")
	router.HandleFunc("/product_attributes/{progenitor_id:[0-9]+}", buildProductAttributeCreationHandler(db)).Methods("POST")
	router.HandleFunc("/product_attributes/{progenitor_id:[0-9]+}/{attribute_id:[0-9]+}", buildProductAttributeUpdateHandler(db)).Methods("PUT")

	// // Product Attribute Values
	// router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/values", buildProductAttributeValueCreationHandler(db)).Methods("GET")
	// router.HandleFunc("/product_attributes/{attribute_id:[0-9]+}/value", buildProductAttributeValueCreationHandler(db)).Methods("POST")
}
