package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
)

const (
	productTableHeaders = `id,
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
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`
	// TODO: change this and the other quantity field to a uint32
	QuantityPerPackage int32 `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`
}

// newProductFromCreationInput creates a new product from a ProductCreationInput
func newProductFromCreationInput(in *ProductCreationInput) *Product {
	np := &Product{
		Name:               in.Name,
		Subtitle:           NullString{sql.NullString{String: in.Subtitle, Valid: true}},
		Description:        in.Description,
		SKU:                in.SKU,
		UPC:                NullString{sql.NullString{String: in.UPC, Valid: true}},
		Manufacturer:       NullString{sql.NullString{String: in.Manufacturer, Valid: true}},
		Brand:              NullString{sql.NullString{String: in.Brand, Valid: true}},
		Quantity:           in.Quantity,
		Taxable:            in.Taxable,
		Price:              in.Price,
		OnSale:             in.OnSale,
		SalePrice:          in.SalePrice,
		Cost:               in.Cost,
		ProductWeight:      in.ProductWeight,
		ProductHeight:      in.ProductHeight,
		ProductWidth:       in.ProductWidth,
		ProductLength:      in.ProductLength,
		PackageWeight:      in.PackageWeight,
		PackageHeight:      in.PackageHeight,
		PackageWidth:       in.PackageWidth,
		PackageLength:      in.PackageLength,
		QuantityPerPackage: in.QuantityPerPackage,
		AvailableOn:        in.AvailableOn,
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

func buildProductExistenceHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

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
		sku := chi.URLParam(req, "sku")

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
		sku := chi.URLParam(req, "sku")

		// can't delete a product that doesn't exist!
		exists, err := rowExistsInDB(db, skuExistenceQuery, sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		err = deleteProductBySKU(db, sku)
		if err != nil {
			notifyOfInternalIssue(res, err, "archive product in database")
			return
		}

		io.WriteString(res, fmt.Sprintf("Successfully deleted product `%s`", sku))
	}
}

func updateProductInDatabase(db *sqlx.DB, up *Product) error {
	// FIXME: this update function is not like the others.
	productUpdateQuery, queryArgs := buildProductUpdateQuery(up)
	err := db.QueryRowx(productUpdateQuery, queryArgs...).StructScan(up)
	return err
}

func buildProductUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductUpdateHandler is a request handler that can update products
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		newerProduct := &Product{}
		err := validateRequestInput(req, newerProduct)
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
		mergo.Merge(newerProduct, &existingProduct)

		if !restrictedStringIsValid(newerProduct.SKU) {
			notifyOfInvalidRequestBody(res, fmt.Errorf("The sku received (%s) is invalid", newerProduct.SKU))
			return
		}
		err = updateProductInDatabase(db, newerProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		json.NewEncoder(res).Encode(newerProduct)
	}
}

// createProductInDB takes a marshaled Product object and creates an entry for it and a base_product in the database
func createProductInDB(tx *sql.Tx, np *Product) (uint64, time.Time, error) {
	var newProductID uint64
	var createdOn time.Time
	productCreationQuery, queryArgs := buildProductCreationQuery(np)
	err := tx.QueryRow(productCreationQuery, queryArgs...).Scan(&newProductID, &createdOn)
	return newProductID, createdOn, err
}

func buildProductCreationHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput := &ProductCreationInput{}
		err := validateRequestInput(req, productInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		if !restrictedStringIsValid(productInput.SKU) {
			notifyOfInvalidRequestBody(res, fmt.Errorf("The sku received (%s) is invalid", productInput.SKU))
			return
		}

		// can't create a product with a sku that already exists!
		exists, err := rowExistsInDB(db, skuExistenceQuery, productInput.SKU)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product with sku '%s' already exists", productInput.SKU))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "create new database transaction")
			return
		}

		newProduct := newProductFromCreationInput(productInput)
		newProductID, createdOn, err := createProductInDB(tx, newProduct)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProduct.ID = newProductID
		newProduct.CreatedOn = createdOn

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
