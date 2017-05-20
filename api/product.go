package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
)

// Product is the parent product for every product
type Product struct {
	// Basic Info
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Price   float32 `json:"base_price"`

	// Product Dimensions
	ProductWeight float32 `json:"base_product_weight"`
	ProductHeight float32 `json:"base_product_height"`
	ProductWidth  float32 `json:"base_product_width"`
	ProductLength float32 `json:"base_product_length"`

	// Package dimensions
	PackageWeight float32 `json:"base_package_weight"`
	PackageHeight float32 `json:"base_package_height"`
	PackageWidth  float32 `json:"base_package_width"`
	PackageLength float32 `json:"base_package_length"`

	// Other Tables
	ProductAttributes []*ProductAttribute `json:"product_attributes"`
	Variants          []*Variant          `json:"variants"`

	// Housekeeping
	CreatedAt  time.Time `json:"created"`
	ArchivedAt NullTime  `json:"-"`
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

// retrieveBaseProductFromDB retrieves a product with a given SKU from the database
func retrieveBaseProductFromDB(db *pg.DB, id int64) (*Product, error) {
	p := &Product{}
	product := db.Model(p).
		Where("id = ?", id).
		Relation("ChildProducts", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("base_product_id = ?", id), nil
		}).
		Relation("ProductAttributes", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("base_product_id = ?", id), nil
		}).
		Where("base_product.archived_at is null")

	err := product.Select()
	return p, err
}

func buildSingleBaseProductHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// singleBaseProductHandler is a request handler that returns a single BaseProduct
	return func(res http.ResponseWriter, req *http.Request) {
		baseProductID := mux.Vars(req)["id"]

		// we can eat this error because Mux takes care of validating route params for us
		actualID, _ := strconv.ParseInt(baseProductID, 10, 64)

		baseProduct, err := retrieveBaseProductFromDB(db, actualID)
		if err != nil {
			informOfServerIssue(err, "Error encountered querying for base_product", res)
			return
		}

		json.NewEncoder(res).Encode(baseProduct)
	}
}
