package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
)

// BaseProduct is the parent product for every product
type BaseProduct struct {
	// Basic Info
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Pricing Fields
	Taxable               bool    `json:"taxable"`
	CustomerCanSetPricing bool    `json:"customer_can_set_pricing"`
	BasePrice             float32 `json:"base_price"`

	// Product Dimensions
	BaseProductWeight float32 `json:"base_product_weight"`
	BaseProductHeight float32 `json:"base_product_height"`
	BaseProductWidth  float32 `json:"base_product_width"`
	BaseProductLength float32 `json:"base_product_length"`

	// Package dimensions
	BasePackageWeight float32 `json:"base_package_weight"`
	BasePackageHeight float32 `json:"base_package_height"`
	BasePackageWidth  float32 `json:"base_package_width"`
	BasePackageLength float32 `json:"base_package_length"`

	// Other Tables
	ProductAttributes []*ProductAttribute `json:"product_attributes"`
	ChildProducts     []*Product          `json:"child_products"`

	// Housekeeping
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created"`
	ArchivedAt time.Time `json:"-"`
}

// NewBaseProductFromProduct takes a Product object and create a BaseProduct from it
func NewBaseProductFromProduct(p *Product) *BaseProduct {
	bp := &BaseProduct{
		Name:                  p.Name,
		Description:           p.Description,
		Taxable:               p.Taxable,
		CustomerCanSetPricing: p.CustomerCanSetPricing,
		BasePrice:             p.Price,
		BaseProductWeight:     p.ProductWeight,
		BaseProductHeight:     p.ProductHeight,
		BaseProductWidth:      p.ProductWidth,
		BaseProductLength:     p.ProductLength,
		BasePackageWeight:     p.PackageWeight,
		BasePackageHeight:     p.PackageHeight,
		BasePackageWidth:      p.PackageWidth,
		BasePackageLength:     p.PackageLength,
	}

	return bp
}

// RetrieveBaseProductFromDB retrieves a product with a given SKU from the database
func RetrieveBaseProductFromDB(id int64) (*BaseProduct, error) {
	bp := &BaseProduct{}
	baseProduct := db.Model(bp).
		Where("id = ?", id).
		Relation("ChildProducts", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("base_product_id = ?", id), nil
		}).
		Relation("ProductAttributes", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("base_product_id = ?", id), nil
		}).
		Where("base_product.archived_at is null")

	err := baseProduct.Select()
	return bp, err
}

// SingleBaseProductHandler is a request handler that returns a single BaseProduct
func SingleBaseProductHandler(res http.ResponseWriter, req *http.Request) {
	baseProductID := mux.Vars(req)["id"]

	// we can eat this error because Mux takes care of validating route params for us
	actualID, _ := strconv.ParseInt(baseProductID, 10, 64)

	baseProduct, err := RetrieveBaseProductFromDB(actualID)
	if err != nil {
		informOfServerIssue(err, "Error encountered querying for base_product", res)
		return
	}

	json.NewEncoder(res).Encode(baseProduct)
}
