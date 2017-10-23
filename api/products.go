package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	// "github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
)

const (
	skuExistenceQuery             = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND archived_on IS NULL)`
	productExistenceQuery         = `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND archived_on IS NULL)`
	productDeletionQuery          = `UPDATE products SET archived_on = NOW() WHERE sku = $1 AND archived_on IS NULL`
	completeProductRetrievalQuery = `SELECT * FROM products WHERE sku = $1`
)

// Product describes something a user can buy
type Product struct {
	DBRow
	// Basic Info
	ProductRootID      uint64 `json:"product_root_id"`
	Name               string `json:"name"`
	Subtitle           string `json:"subtitle"`
	Description        string `json:"description"`
	OptionSummary      string `json:"option_summary"`
	SKU                string `json:"sku"`
	UPC                string `json:"upc"`
	Manufacturer       string `json:"manufacturer"`
	Brand              string `json:"brand"`
	Quantity           uint32 `json:"quantity"`
	QuantityPerPackage uint32 `json:"quantity_per_package"`

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

	ApplicableOptionValues []ProductOptionValue `json:"applicable_options,omitempty"`

	AvailableOn time.Time `json:"available_on"`
}

// newProductFromCreationInput creates a new product from a ProductCreationInput
func newProductFromCreationInput(in *ProductCreationInput) *Product {
	np := &Product{
		Name:               in.Name,
		Subtitle:           in.Subtitle,
		Description:        in.Description,
		SKU:                in.SKU,
		UPC:                in.UPC,
		Manufacturer:       in.Manufacturer,
		Brand:              in.Brand,
		Quantity:           in.Quantity,
		QuantityPerPackage: in.QuantityPerPackage,
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
		AvailableOn:        in.AvailableOn,
	}
	return np
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
	Quantity     uint32 `json:"quantity"`

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
	QuantityPerPackage uint32  `json:"quantity_per_package"`

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

func buildSingleProductHandler(dbx *sqlx.DB, db storage.Storage) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		// product, err := retrieveProductFromDB(db, sku)
		product, err := db.GetProductBySKU(sku)
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

		productsResponse := &ListResponse{
			Page:  queryFilter.Page,
			Limit: queryFilter.Limit,
			Count: count,
			Data:  products,
		}
		json.NewEncoder(res).Encode(productsResponse)
	}
}

func deleteProductBySKU(tx *sql.Tx, sku string) error {
	_, err := tx.Exec(productDeletionQuery, sku)
	return err
}

func buildProductDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		// can't delete a product that doesn't exist!
		existingProduct, err := retrieveProductFromDB(db, sku)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "create new database transaction")
			return
		}

		err = deleteProductVariantBridgeEntriesByProductID(tx, existingProduct.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product in database")
			return
		}

		err = deleteProductBySKU(tx, sku)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product in database")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		res.WriteHeader(http.StatusOK)
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
func createProductInDB(tx *sql.Tx, np *Product) (uint64, time.Time, time.Time, error) {
	var newProductID uint64
	var availableOn time.Time
	var createdOn time.Time
	productCreationQuery, queryArgs := buildProductCreationQuery(np)
	err := tx.QueryRow(productCreationQuery, queryArgs...).Scan(&newProductID, &availableOn, &createdOn)
	return newProductID, availableOn, createdOn, err
}

func createProductsInDBFromOptionRows(tx *sql.Tx, r *ProductRoot, np *Product) ([]Product, error) {
	createdProducts := []Product{}
	productOptionData := generateCartesianProductForOptions(r.Options)
	for _, option := range productOptionData {
		p := &Product{}
		*p = *np // solved: http://www.claymath.org/millennium-problems/p-vs-np-problem

		p.ProductRootID = r.ID
		p.ApplicableOptionValues = option.OriginalValues
		p.OptionSummary = option.OptionSummary
		p.SKU = fmt.Sprintf("%s_%s", r.SKUPrefix, option.SKUPostfix)

		var err error
		p.ID, p.AvailableOn, p.CreatedOn, err = createProductInDB(tx, p)
		if err != nil {
			return nil, err
		}

		err = createBridgeEntryForProductValues(tx, p.ID, option.IDs)
		if err != nil {
			return nil, err
		}
		createdProducts = append(createdProducts, *p)
	}
	return createdProducts, nil
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
		exists, err := rowExistsInDB(db, productRootSkuExistenceQuery, productInput.SKU)
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
		productRoot := createProductRootFromProduct(newProduct)
		productRoot.ID, productRoot.CreatedOn, err = createProductRootInDB(tx, productRoot)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product options and values in database")
			return
		}

		for _, optionAndValues := range productInput.Options {
			o, err := createProductOptionAndValuesInDBFromInput(tx, optionAndValues, productRoot.ID)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product options and values in database")
				return
			}
			productRoot.Options = append(productRoot.Options, o)
		}

		if len(productInput.Options) == 0 {
			newProduct.ProductRootID = productRoot.ID
			newProduct.ID, newProduct.AvailableOn, newProduct.CreatedOn, err = createProductInDB(tx, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product in database")
				return
			}
			productRoot.Options = []*ProductOption{} // so this won't be Marshaled as null
			productRoot.Products = []Product{*newProduct}
		} else {
			productRoot.Products, err = createProductsInDBFromOptionRows(tx, productRoot, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert products in database")
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(productRoot)
	}
}
