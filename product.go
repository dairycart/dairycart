package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg/orm"
)

// Product describes...well, a product
type Product struct {
	// Basic Info
	ID          int64  `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UPC         string `json:"upc"`
	Quantity    int32  `json:"quantity"`

	// Pricing Fields
	OnSale                bool    `json:"on_sale"`
	Price                 float32 `json:"price"`
	SalePrice             float32 `json:"sale_price"`
	Taxable               bool    `json:"taxable"`
	CustomerCanSetPricing bool    `json:"customer_can_set_pricing"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`

	// Housekeeping
	Active   bool       `json:"-"`
	Created  *time.Time `json:"created"`
	Archived *time.Time `json:"archived,omitempty"`
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	Count int64     `json:"count"`
	Limit int64     `json:"limit"`
	Data  []Product `json:"data"`
}

// ProductListHandler is a generic product list request handler
func ProductListHandler(res http.ResponseWriter, req *http.Request) {
	actualLimit := DetermineRequestLimits(req)

	var products []Product
	productsModel := db.Model(&products)

	productsModel.Apply(orm.URLFilters(req.URL.Query()))
	err := productsModel.Limit(actualLimit).Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
	}
	productsResponse := &ProductsResponse{
		Limit: int64(actualLimit),
		Count: int64(len(products)),
		Data:  products,
	}
	json.NewEncoder(res).Encode(productsResponse)
}

// ProductCreationHandler is a product creation handler
func ProductCreationHandler(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(res, "Please send a request body", http.StatusBadRequest)
		return
	}

	newProduct := &Product{}
	err := json.NewDecoder(req.Body).Decode(newProduct)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
		// fmt.Fprintf(w, "Error encountered parsing request: %v", err)
	}

	err = db.Insert(newProduct)
	if err != nil {
		log.Printf("error inserting product into database: %v", err)
	}
}
