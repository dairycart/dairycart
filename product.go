package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
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

// Product describes...well, a product
type Product struct {
	// Basic Info
	ID            int64        `json:"id"`
	BaseProductID int64        `json:"-"`
	BaseProduct   *BaseProduct `json:"-"`
	SKU           string       `json:"sku"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	UPC           string       `json:"upc"`
	Quantity      int          `json:"quantity"`

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
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created"`
	ArchivedAt time.Time `json:"-"`
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

func filterActiveProducts(req *http.Request, q *orm.Query) {
	if req.URL.Query().Get("include_archived") != "true" {
		q.Where("product.archived_at is null")
	}
}

// ProductListHandler is a request handler that returns a list of products
func ProductListHandler(res http.ResponseWriter, req *http.Request) {
	var products []Product
	productsModel := db.Model(&products).
		Column("product.*", "BaseProduct").
		Where("base_product.archived_at is null")

	pager, err := genericListQueryHandler(req, productsModel, filterActiveProducts)
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
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
	vars := mux.Vars(req)
	sku := vars["sku"]

	var p Product
	product := db.Model(&p).
		Column("products.*", "BaseProduct").
		Where("sku = ?", sku)

	filterActiveProducts(req, product)

	err := product.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(product)
}

// ProductExistenceHandler handles requests to check if a sku exists
func ProductExistenceHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sku := vars["sku"]

	var p Product
	product := db.Model(&p).Where("sku = ?", sku)
	filterActiveProducts(req, product)

	productCount, err := product.Count()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	responseStatus := http.StatusNotFound
	if productCount == 1 {
		responseStatus = http.StatusOK
	}
	res.WriteHeader(responseStatus)
}

// ProductUpdateHandler is a request handler that can update products
func ProductUpdateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(res, "Please send a request body", http.StatusBadRequest)
	}

	vars := mux.Vars(req)
	sku := vars["sku"]
	var existingProduct Product
	existingProductQuery := db.Model(&existingProduct).Where("sku = ?", sku)

	updatedProduct := &Product{}
	err := json.NewDecoder(req.Body).Decode(updatedProduct)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingProductQuery.Select()
	updatedProduct.ID = existingProduct.ID
	if err := mergo.Merge(updatedProduct, existingProduct); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}
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

	err = db.RunInTransaction(func(tx *pg.Tx) error {
		baseProduct := NewBaseProductFromProduct(newProduct)
		newProduct.BaseProduct = baseProduct
		err = tx.Insert(baseProduct)
		if err != nil {
			return err
		}

		newProduct.BaseProductID = baseProduct.ID
		err = tx.Insert(newProduct)
		return err
	})

	if err != nil {
		errorString := fmt.Sprintf("error inserting product into database: %v", err)
		log.Println(errorString)
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}
}

// ProductDeletionHandler is a request handler that deletes a single product
func ProductDeletionHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sku := vars["sku"]

	var p Product
	product := db.Model(&p).Where("sku = ?", sku)

	err := product.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Model(&p).Set("archived_at = now()").Set("active = false").Where("sku = ?", sku).Update(&p)

	json.NewEncoder(res).Encode(product)
}
