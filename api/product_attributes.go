package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID                  int64       `json:"id"`
	Name                string      `json:"name"`
	ProductProgenitorID int64       `json:"product_progenitor_id"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           pq.NullTime `json:"-"`
	ArchivedAt          pq.NullTime `json:"-"`
}

func (a ProductAttribute) generateScanArgs() []interface{} {
	return []interface{}{
		&a.ID,
		&a.Name,
		&a.ProductProgenitorID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.ArchivedAt,
	}
}

// ProductAttributesResponse is a product attribute response struct
type ProductAttributesResponse struct {
	ListResponse
	Data []ProductAttribute `json:"data"`
}

func getProductAttributesForProgenitor(db *sql.DB, progenitorID string, queryFilter *QueryFilter) ([]ProductAttribute, error) {
	var attributes []ProductAttribute

	rows, err := db.Query(buildProductAttributeListQuery(progenitorID, queryFilter))
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	defer rows.Close()
	for rows.Next() {
		var attribute ProductAttribute
		_ = rows.Scan(attribute.generateScanArgs()...)
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

func buildProductAttributeListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		progenitorID := mux.Vars(req)["progenitor_id"]
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		attributes, err := getProductAttributesForProgenitor(db, progenitorID, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		attributesResponse := &ProductAttributesResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: uint64(len(attributes)),
			},
			Data: attributes,
		}
		json.NewEncoder(res).Encode(attributesResponse)
	}
}

func createProductAttributeInDB(db *sql.DB, a *ProductAttribute) error {
	query, args := buildProductAttributeCreationQuery(a)
	err := db.QueryRow(query, args...).Scan(a.generateScanArgs()...)
	return err
}

// func buildProductAttributeCreationHandler(db *sql.DB) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {

// 	}
// }

// func updateProductAttributeInDB(db *sql.DB, a *ProductAttribute) error {
// 	productUpdateQuery, queryArgs := buildProductAttributeUpdateQuery(a)
// 	err := db.QueryRow(productUpdateQuery, queryArgs...).Scan(a.generateScanArgs()...)
// 	return err
// }

// func buildProductAttributeUpdateHandler(db *sql.DB) http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {

// 	}
// }
