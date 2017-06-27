package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	productTableHeaders = `
		id,
		name,
		subtitle,
		description,
		sku,
		upc,
		manufacturer,
		brand,
		quantity,
		taxable,
		price,
		on_sale,
		sale_price,
		cost,
		product_weight,
		product_height,
		product_width,
		product_length,
		package_weight,
		package_height,
		package_width,
		package_length,
		quantity_per_package,
		available_on,
		created_on,
		updated_on,
		archived_on
	`

	skuExistenceQuery             = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND archived_on IS NULL)`
	productExistenceQuery         = `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND archived_on IS NULL)`
	productDeletionQuery          = `UPDATE products SET archived_on = NOW() WHERE sku = $1 AND archived_on IS NULL`
	completeProductRetrievalQuery = `SELECT * FROM products WHERE sku = $1`
)

// Product describes something a user can buy
type Product struct {
	DBRow
	// Basic Info
	Name         string     `json:"name"`
	Subtitle     NullString `json:"subtitle"`
	Description  string     `json:"description"`
	SKU          string     `json:"sku"`
	UPC          NullString `json:"upc"`
	Manufacturer NullString `json:"manufacturer"`
	Brand        NullString `json:"brand"`
	Quantity     int        `json:"quantity"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight      float32 `json:"package_weight"`
	PackageHeight      float32 `json:"package_height"`
	PackageWidth       float32 `json:"package_width"`
	PackageLength      float32 `json:"package_length"`
	QuantityPerPackage int32   `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`
}

// newProductFromCreationInput creates a new product from a ProductCreationInput
func newProductFromCreationInput(in *ProductCreationInput) *Product {
	np := &Product{
		SKU:      in.SKU,
		Name:     in.Name,
		UPC:      NullString{sql.NullString{String: in.UPC, Valid: true}},
		Quantity: in.Quantity,
		Price:    in.Price,
		Cost:     in.Cost,
	}
	return np
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	// Core Product stuff
	Name         string `json:"name"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
	SKU          string `json:"sku"`
	UPC          string `json:"upc"`
	Manufacturer string `json:"manufacturer"`
	Brand        string `json:"brand"`
	Quantity     int    `json:"quantity"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight      float32 `json:"package_weight"`
	PackageHeight      float32 `json:"package_height"`
	PackageWidth       float32 `json:"package_width"`
	PackageLength      float32 `json:"package_length"`
	QuantityPerPackage int32   `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`

	// Other things
	Options []*ProductOptionCreationInput `json:"options"`
}

func validateProductUpdateInput(req *http.Request) (*Product, error) {
	product := &Product{}
	err := json.NewDecoder(req.Body).Decode(product)
	if err != nil {
		return nil, err
	}

	p := structs.New(product)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !dataValueIsValid(product.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	return product, err
}

func buildProductExistenceHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		sku := vars["sku"]

		productExists, err := rowExistsInDB(db, skuExistenceQuery, sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		responseStatus := http.StatusNotFound
		if productExists {
			responseStatus = http.StatusOK
		}
		res.WriteHeader(responseStatus)
	}
}

// retrieveProductFromDB retrieves a product with a given SKU from the database
func retrieveProductFromDB(db *sqlx.DB, sku string) (Product, error) {
	var p Product
	err := db.Get(&p, completeProductRetrievalQuery, sku)
	return p, err
}

func buildSingleProductHandler(db *sqlx.DB) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		product, err := retrieveProductFromDB(db, sku)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving product from database")
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func buildProductListHandler(db *sqlx.DB) http.HandlerFunc {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := getRowCount(db, "products", queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of products from the database")
			return
		}

		var products []Product
		query, args := buildProductListQuery(queryFilter)
		err = retrieveListOfRowsFromDB(db, query, args, &products)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		productsResponse := &ProductsResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: count,
			},
			Data: products,
		}
		json.NewEncoder(res).Encode(productsResponse)
	}
}

func deleteProductBySKU(db *sqlx.DB, sku string) error {
	_, err := db.Exec(productDeletionQuery, sku)
	return err
}

func buildProductDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		// can't delete a product that doesn't exist!
		exists, err := rowExistsInDB(db, skuExistenceQuery, sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		err = deleteProductBySKU(db, sku)
		io.WriteString(res, fmt.Sprintf("Successfully deleted product `%s`", sku))
	}
}

func updateProductInDatabase(db *sqlx.DB, up *Product) error {
	productUpdateQuery, queryArgs := buildProductUpdateQuery(up)
	err := db.QueryRowx(productUpdateQuery, queryArgs...).StructScan(up)
	return err
}

func buildProductUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductUpdateHandler is a request handler that can update products
		sku := mux.Vars(req)["sku"]

		newerProduct, err := validateProductUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingProduct, err := retrieveProductFromDB(db, sku)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(newerProduct, existingProduct)

		err = updateProductInDatabase(db, newerProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		json.NewEncoder(res).Encode(newerProduct)
	}
}

func validateProductCreationInput(req *http.Request) (*ProductCreationInput, error) {
	pci := &ProductCreationInput{}
	err := json.NewDecoder(req.Body).Decode(pci)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}

	p := structs.New(pci)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !dataValueIsValid(pci.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	return pci, err
}

// createProductInDB takes a marshaled Product object and creates an entry for it and a base_product in the database
func createProductInDB(tx *sql.Tx, np *Product) (uint64, error) {
	var newProductID uint64
	productCreationQuery, queryArgs := buildProductCreationQuery(np)
	err := tx.QueryRow(productCreationQuery, queryArgs...).Scan(&newProductID)
	return newProductID, err
}

func noop(things ...interface{}) {
	return
}

func buildProductCreationHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput, err := validateProductCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create a product with a sku that already exists!
		exists, err := rowExistsInDB(db, skuExistenceQuery, productInput.SKU)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product with sku `%s` already exists", productInput.SKU))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "create new database transaction")
			return
		}

		newProduct := newProductFromCreationInput(productInput)
		newProductID, err := createProductInDB(tx, newProduct)
		if err != nil {
			errStr := err.Error()
			noop(errStr)

			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProduct.ID = newProductID

		for _, optionAndValues := range productInput.Options {
			_, err = createProductOptionAndValuesInDBFromInput(tx, optionAndValues, newProduct.ID)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product options and values in database")
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newProduct)
	}
}
