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

// ProductExistsInDB will return whether or not a product with a given sku exists in the database
func ProductExistsInDB(sku string) (bool, error) {
	product := db.Model(&Product{}).Where("sku = ?", sku).Where("archived_at is null")

	productCount, err := product.Count()
	return productCount == 1, err
}

// ProductExistenceHandler handles requests to check if a sku exists
func ProductExistenceHandler(res http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

	productExists, err := ProductExistsInDB(sku)
	if err != nil {
		informOfServerIssue(err, "Error encountered querying for product", res)
		return
	}

	responseStatus := http.StatusNotFound
	if productExists {
		responseStatus = http.StatusOK
	}
	res.WriteHeader(responseStatus)
}

// RetrieveProductFromDB retrieves a product with a given SKU from the database
func RetrieveProductFromDB(sku string) (*Product, error) {
	p := &Product{}
	product := db.Model(p).
		Where("sku = ?", sku).
		Where("product.archived_at is null")

	err := product.Select()
	return p, err
}

// SingleProductHandler is a request handler that returns a single Product
func SingleProductHandler(res http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

	product, err := RetrieveProductFromDB(sku)

	if err != nil {
		informOfServerIssue(err, "Error encountered querying for product", res)
		return
	}

	json.NewEncoder(res).Encode(product)
}

// ProductListHandler is a request handler that returns a list of products
func ProductListHandler(res http.ResponseWriter, req *http.Request) {
	var products []Product
	productsModel := db.Model(&products).
		Column("product.*", "BaseProduct").
		Where("base_product.archived_at is null")

	pager, err := genericListQueryHandler(req, productsModel, filterActiveProducts)
	if err != nil {
		informOfServerIssue(err, "Error encountered querying for products", res)
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

// ProductUpdateHandler is a request handler that can update products
func ProductUpdateHandler(res http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]
	var existingProduct Product
	existingProductQuery := db.Model(&existingProduct).Where("sku = ?", sku)

	updatedProduct := &Product{}
	bodyIsInvalid := ensureRequestBodyValidity(res, req, updatedProduct)
	if !bodyIsInvalid {
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

// CreateProduct takes a marshalled Product object and creates an entry for it and a base_product in the database
func CreateProduct(newProduct *Product) error {
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		baseProduct := NewBaseProductFromProduct(newProduct)
		newProduct.BaseProduct = baseProduct
		err := tx.Insert(baseProduct)
		if err != nil {
			return err
		}

		newProduct.BaseProductID = baseProduct.ID
		err = tx.Insert(newProduct)
		return err
	})
	return err
}

// ProductCreationHandler is a product creation handler
func ProductCreationHandler(res http.ResponseWriter, req *http.Request) {
	newProduct := &Product{}
	bodyIsInvalid := ensureRequestBodyValidity(res, req, newProduct)
	if bodyIsInvalid {
		return
	}

	err := CreateProduct(newProduct)
	if err != nil {
		errorString := fmt.Sprintf("error inserting product into database: %v", err)
		log.Println(errorString)
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}
}

// ProductDeletionHandler is a request handler that deletes a single product
func ProductDeletionHandler(res http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

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
