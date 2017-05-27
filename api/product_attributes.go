package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/lib/pq"
)

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID                  int64       `json:"id"`
	Name                string      `json:"Name"`
	ProductProgenitorID int64       `json:"product_progenitor_id"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           pq.NullTime `json:"-"`
	ArchivedAt          pq.NullTime `json:"-"`
}

func getProductAttributesForProduct(db *sql.DB, progenitorID int64) ([]*ProductAttribute, error) {
	return nil, nil
}

func buildProductAttributeListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}

func createProductAttributeInDB(db *sql.DB, pa *ProductAttribute) error {
	return nil
}

func buildProductAttributeCreationHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}

func updateProductAttributeInDB(db *sql.DB, pa *ProductAttribute) error {
	return nil
}

func buildProductAttributeUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}
