package api

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
	"github.com/pkg/errors"
)

const (
	skuExistenceQuery             = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND archived_on IS NULL)`
	productDeletionQuery          = `UPDATE products SET archived_on = NOW() WHERE sku = $1 AND archived_on IS NULL`
	completeProductRetrievalQuery = `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.sku = $1`
)

// Product describes something a user can buy
type Product struct {
	// Basic Info
	ID                  int64      `json:"id"`
	ProductProgenitorID int64      `json:"product_progenitor_id"`
	SKU                 string     `json:"sku"`
	Name                string     `json:"name"`
	UPC                 NullString `json:"upc"`
	Quantity            int        `json:"quantity"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Price   float32 `json:"price"`
	Cost    float32 `json:"cost"`

	// Inheritor
	ProductProgenitor

	// Housekeeping
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  NullTime  `json:"updated_on,omitempty"`
	ArchivedOn NullTime  `json:"archived_on,omitempty"`
}

// generateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (p *Product) generateScanArgs() []interface{} {
	return []interface{}{
		&p.ID,
		&p.ProductProgenitorID,
		&p.SKU,
		&p.Name,
		&p.UPC,
		&p.Quantity,
		&p.Price,
		&p.Cost,
		&p.CreatedOn,
		&p.UpdatedOn,
		&p.ArchivedOn,
	}
}

// generateJoinScanArgs does some stuff TODO: write better docs
func (p *Product) generateJoinScanArgs() []interface{} {
	productScanArgs := p.generateScanArgs()
	progenitorScanArgs := p.ProductProgenitor.generateScanArgs()
	return append(productScanArgs, progenitorScanArgs...)
}

// generateJoinScanArgsWithCount does everything generateJoinScanArgs does,
// only with an added count parameter
func (p *Product) generateJoinScanArgsWithCount(count *uint64) []interface{} {
	scanArgs := []interface{}{count}
	productScanArgs := p.generateScanArgs()
	progenitorScanArgs := p.ProductProgenitor.generateScanArgs()
	scanArgs = append(scanArgs, productScanArgs...)
	return append(scanArgs, progenitorScanArgs...)
}

// newProductFromCreationInputAndProgenitor creates a new product from a ProductProgenitor and a ProductCreationInput
func newProductFromCreationInputAndProgenitor(g *ProductProgenitor, in *ProductCreationInput) *Product {
	np := &Product{
		ProductProgenitor:   *g,
		ProductProgenitorID: g.ID,
		SKU:                 in.SKU,
		Name:                in.Name,
		UPC:                 NullString{sql.NullString{String: in.UPC, Valid: in.UPC != ""}},
		Quantity:            in.Quantity,
		Price:               in.Price,
		Cost:                in.Cost,
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
	Description   string                        `json:"description"`
	Taxable       bool                          `json:"taxable"`
	ProductWeight float32                       `json:"product_weight"`
	ProductHeight float32                       `json:"product_height"`
	ProductWidth  float32                       `json:"product_width"`
	ProductLength float32                       `json:"product_length"`
	PackageWeight float32                       `json:"package_weight"`
	PackageHeight float32                       `json:"package_height"`
	PackageWidth  float32                       `json:"package_width"`
	PackageLength float32                       `json:"package_length"`
	SKU           string                        `json:"sku"`
	Name          string                        `json:"name"`
	UPC           string                        `json:"upc"`
	Quantity      int                           `json:"quantity"`
	Price         float32                       `json:"price"`
	Cost          float32                       `json:"cost"`
	Options       []*ProductOptionCreationInput `json:"options"`
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

func buildProductExistenceHandler(db *sql.DB) http.HandlerFunc {
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
func retrieveProductFromDB(db *sql.DB, sku string) (*Product, error) {
	product := &Product{}
	scanArgs := product.generateJoinScanArgs()
	err := db.QueryRow(completeProductRetrievalQuery, sku).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return product, errors.Wrap(err, "Error querying for product")
	}

	return product, err
}

func buildSingleProductHandler(db *sql.DB) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		productExists, err := rowExistsInDB(db, skuExistenceQuery, sku)
		if err != nil || !productExists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		product, err := retrieveProductFromDB(db, sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func retrieveProductsFromDB(db *sql.DB, queryFilter *QueryFilter) ([]Product, uint64, error) {
	var products []Product
	var count uint64

	query, args := buildProductListQuery(queryFilter)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Error encountered querying for products")
	}

	defer rows.Close()
	for rows.Next() {
		var product Product
		var queryCount uint64

		scanArgs := product.generateJoinScanArgsWithCount(&queryCount)
		_ = rows.Scan(scanArgs...)

		count = queryCount
		products = append(products, product)
	}
	return products, count, nil
}

func buildProductListHandler(db *sql.DB) http.HandlerFunc {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		products, count, err := retrieveProductsFromDB(db, queryFilter)
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

func deleteProductBySKU(db *sql.DB, sku string) error {
	_, err := db.Exec(productDeletionQuery, sku)
	return err
}

func buildProductDeletionHandler(db *sql.DB) http.HandlerFunc {
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

func updateProductInDatabase(db *sql.DB, up *Product) error {
	productUpdateQuery, queryArgs := buildProductUpdateQuery(up)
	scanArgs := up.generateScanArgs()
	err := db.QueryRow(productUpdateQuery, queryArgs...).Scan(scanArgs...)
	return err
}

func buildProductUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductUpdateHandler is a request handler that can update products
		sku := mux.Vars(req)["sku"]

		// can't update a product that doesn't exist!
		exists, err := rowExistsInDB(db, skuExistenceQuery, sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		newerProduct, err := validateProductUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// eating the error here because we're already certain the sku exists
		existingProduct, err := retrieveProductFromDB(db, sku)
		if err != nil {
			notifyOfInternalIssue(res, err, "merge updated product with existing product")
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
func createProductInDB(tx *sql.Tx, np *Product) (int64, error) {
	var newProductID int64
	productCreationQuery, queryArgs := buildProductCreationQuery(np)
	err := tx.QueryRow(productCreationQuery, queryArgs...).Scan(&newProductID)
	return newProductID, err
}

func buildProductCreationHandler(db *sql.DB) http.HandlerFunc {
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

		progenitor := newProductProgenitorFromProductCreationInput(productInput)
		newProgenitorID, err := createProductProgenitorInDB(tx, progenitor)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product progenitor in database")
			return
		}
		progenitor.ID = newProgenitorID

		for _, optionAndValues := range productInput.Options {
			_, err = createProductOptionAndValuesInDBFromInput(tx, optionAndValues, progenitor.ID)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product options and values in database")
				return
			}
		}

		newProduct := newProductFromCreationInputAndProgenitor(progenitor, productInput)
		newProductID, err := createProductInDB(tx, newProduct)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProduct.ID = newProductID

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newProduct)
	}
}
