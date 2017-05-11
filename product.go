package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-pg/pg/orm"
)

// Variant describes children of products with different attributes from the parent
type Variant struct {
	ID int64

	SKU   string
	Type  string
	Value string

	HasSpecialPrice bool
	Price           float32
}

// Product describes...well, a product
type Product struct {
	ID                    int64   `json:"id"`
	SKU                   string  `json:"sku"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	UPC                   string  `json:"upc"`
	OnSale                bool    `json:"on_sale"`
	Taxable               bool    `json:"taxable"`
	CustomerCanSetPricing bool    `json:"customer_can_set_pricing"`
	Price                 float32 `json:"price"`
	Weight                float32 `json:"weight"`
	Height                float32 `json:"height"`
	Width                 float32 `json:"width"`
	Length                float32 `json:"length"`
	Quantity              int32   `json:"quantity"`

	ChildOf *Product `json:"child_of,omitempty"`
	// SalePrice float32
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	Count int64     `json:"count"`
	Limit int64     `json:"limit"`
	Data  []Product `json:"data"`
}

// ProductListHandler is a generic product list request handler
func ProductListHandler(res http.ResponseWriter, req *http.Request) {
	var products []Product
	productsModel := db.Model(&products)
	productsModel.Apply(orm.URLFilters(req.URL.Query()))
	err := productsModel.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
	}
	productsResponse := &ProductsResponse{
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
