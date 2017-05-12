package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Archived time.Time `json:"-"`
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

// ProductListHandler is a request handler that returns a list of products
func ProductListHandler(res http.ResponseWriter, req *http.Request) {
	var products []Product
	productsModel := db.Model(&products)

	pager, err := GenericListQueryHandler(req, productsModel)
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	productsResponse := &ProductsResponse{
		ListResponse: ListResponse{
			Page:  pager.Page(),
			Limit: pager.Limit(),
			Count: len(products),
		},
		Data: products,
	}
	json.NewEncoder(res).Encode(productsResponse)
}

// SingleProductHandler is a request handler that returns a single product
func SingleProductHandler(res http.ResponseWriter, req *http.Request) {
	var p Product
	product := db.Model(&p)

	vars := mux.Vars(req)
	sku := vars["sku"]

	SelectActiveRows(req, product)

	err := product.Where("sku = ?", sku).Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(res).Encode(product)
}

// ProductUpdateHandler is a request handler that can update products
func ProductUpdateHandler(res http.ResponseWriter, req *http.Request) {

	// TODO: This works for now, but users have to provide the entire
	// 		 product object, instead of being able to update just parts
	//		 of it. I consider this lame, and this should be fixed.

	if req.Body == nil {
		http.Error(res, "Please send a request body", http.StatusBadRequest)
	}

	vars := mux.Vars(req)
	sku := vars["sku"]
	log.Printf("sku: %v", sku)

	updatedProduct := &Product{}
	err := json.NewDecoder(req.Body).Decode(updatedProduct)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	var p Product
	_ = db.Model(&p).Where("sku = ?", sku).Returning("id").Select()

	updatedProduct.ID = p.ID
	db.Update(updatedProduct)

	json.NewEncoder(res).Encode(updatedProduct)
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
	}

	err = db.Insert(newProduct)
	if err != nil {
		log.Printf("error inserting product into database: %v", err)
	}
}

// ProductDeletionHandler is a request handler that deletes a single product
func ProductDeletionHandler(res http.ResponseWriter, req *http.Request) {
	var p Product
	product := db.Model(&p)

	vars := mux.Vars(req)
	sku := vars["sku"]

	SelectActiveRows(req, product)

	err := product.Where("sku = ?", sku).Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	db.Model(&p).Set("archived = now()").Set("active = false").Where("sku = ?", sku).Update(&p)

	json.NewEncoder(res).Encode(product)
}
