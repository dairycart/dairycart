package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	skuExistenceQuery         = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 and archived_at is null);`
	skuDeletionQuery          = `UPDATE products SET archived_at = NOW() WHERE sku = $1 AND archived_at IS NULL;`
	skuRetrievalQuery         = `SELECT * FROM products WHERE sku = $1 AND archived_at IS NULL;`
	skuJoinRetrievalQuery     = `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.sku = $1 AND p.archived_at IS NULL;`
	allProductsRetrievalQuery = `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.id IS NOT NULL AND p.archived_at IS NULL;`
	productUpdateQuery        = `UPDATE products SET "sku"=$1, "name"=$2, "upc"=$3, "quantity"=$4, "on_sale"=$5, "price"=$6, "sale_price"=$7, "updated_at"='NOW()' WHERE "id"=$8;`
	productCreationQuery      = `INSERT INTO products ("product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	skuValidationPattern = `^[a-zA-Z\-_]+$`
)

var skuValidator *regexp.Regexp

func init() {
	skuValidator = regexp.MustCompile(skuValidationPattern)
}

// Product describes something a user can buy
type Product struct {
	ProductProgenitor

	// Basic Info
	ID                  int64      `json:"id"`
	ProductProgenitorID int64      `json:"product_progenitor_id"`
	SKU                 string     `json:"sku"`
	Name                string     `json:"name"`
	UPC                 NullString `json:"upc"`
	Quantity            int        `json:"quantity"`

	// Pricing Fields
	OnSale    bool        `json:"on_sale"`
	Price     float32     `json:"price"`
	SalePrice NullFloat64 `json:"sale_price"`

	// // Housekeeping
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  pq.NullTime `json:"updated_at"`
	ArchivedAt pq.NullTime `json:"-"`
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
		&p.OnSale,
		&p.Price,
		&p.SalePrice,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.ArchivedAt,
	}
}

// generateJoinScanArgs does some stuff TODO: write better docs
func (p *Product) generateJoinScanArgs() []interface{} {
	productScanArgs := p.generateScanArgs()
	progenitorScanArgs := p.ProductProgenitor.generateScanArgs()
	return append(productScanArgs, progenitorScanArgs...)
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

func validateProductUpdateInput(req *http.Request) (*Product, error) {
	product := &Product{}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(product)

	p := structs.New(product)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !skuValidator.MatchString(product.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	// // TODO: revisit this later
	// formatted, err := strconv.ParseFloat(fmt.Sprintf("%.2f", rounded), 64)
	// product.Price =

	return product, err
}

func buildProductExistenceHandler(db *sql.DB) http.HandlerFunc {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		sku := vars["sku"]

		productExists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", "sku", sku)
			return
		}

		responseStatus := http.StatusNotFound
		if productExists {
			responseStatus = http.StatusOK
		}
		res.WriteHeader(responseStatus)
	}
}

// retrievePlainProductFromDB retrieves a product with a given SKU from the database
func retrievePlainProductFromDB(db *sql.DB, sku string) (*Product, error) {
	product := &Product{}
	scanArgs := product.generateScanArgs()

	err := db.QueryRow(skuRetrievalQuery, sku).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return product, errors.Wrap(err, "Error querying for product")
	}

	return product, nil
}

// retrieveProductFromDB retrieves a product with a given SKU from the database
func retrieveProductFromDB(db *sql.DB, sku string) (*Product, error) {
	product := &Product{}
	scanArgs := product.generateJoinScanArgs()

	err := db.QueryRow(skuJoinRetrievalQuery, sku).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return product, errors.Wrap(err, "Error querying for product")
	}

	return product, err
}

func buildSingleProductHandler(db *sql.DB) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		product, err := retrieveProductFromDB(db, sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", "sku", sku)
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func retrieveProductsFromDB(db *sql.DB) ([]Product, error) {
	var products []Product

	rows, err := db.Query(allProductsRetrievalQuery)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		_ = rows.Scan(product.generateJoinScanArgs()...)
		products = append(products, product)
	}
	return products, nil
}

func buildProductListHandler(db *sql.DB) http.HandlerFunc {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		products, err := retrieveProductsFromDB(db)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		productsResponse := &ProductsResponse{
			ListResponse: ListResponse{
				Page:  1,  // TODO: implement proper paging :(
				Limit: 25, // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
				Count: len(products),
			},
			Data: products,
		}
		json.NewEncoder(res).Encode(productsResponse)
	}
}

func deleteProductBySKU(db *sql.DB, sku string) error {
	_, err := db.Exec(skuDeletionQuery, sku)
	return err
}

func buildProductDeletionHandler(db *sql.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		// can't delete a product that doesn't exist!
		exists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", "sku", sku)
			return
		}

		err = deleteProductBySKU(db, sku)
		json.NewEncoder(res).Encode("OK")
	}
}

func updateProductInDatabase(db *sql.DB, up *Product) error {
	_, err := db.Exec(productUpdateQuery, up.SKU, up.Name, up.UPC, up.Quantity, up.OnSale, up.Price, up.SalePrice, up.ID)
	return err
}

func buildProductUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductUpdateHandler is a request handler that can update products
		sku := mux.Vars(req)["sku"]

		// can't update a product that doesn't exist!
		exists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", "sku", sku)
			return
		}
		// eating the error here because we're already certain the sku exists
		existingProduct, _ := retrievePlainProductFromDB(db, sku)

		newerProduct, err := validateProductUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		if err := mergo.Merge(newerProduct, existingProduct); err != nil {
			notifyOfInternalIssue(res, err, "merge updated product with existing product")
			return
		}

		err = updateProductInDatabase(db, newerProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		json.NewEncoder(res).Encode("Product updated")
	}
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	Description   string  `json:"description"`
	Taxable       bool    `json:"taxable"`
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	UPC           string  `json:"upc"`
	Quantity      int     `json:"quantity"`
	OnSale        bool    `json:"on_sale"`
	Price         float32 `json:"price"`
	SalePrice     float64 `json:"sale_price"`
}

func validateProductCreationInput(req *http.Request) (*ProductCreationInput, error) {
	newProduct := &ProductCreationInput{}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(newProduct)

	p := structs.New(newProduct)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !skuValidator.MatchString(newProduct.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	return newProduct, err
}

// createProduct takes a marshalled Product object and creates an entry for it and a base_product in the database
func createProduct(db *sql.DB, np *Product) error {
	_, err := db.Exec(productCreationQuery, np.ProductProgenitorID, np.SKU, np.Name, np.UPC, np.Quantity, np.OnSale, np.Price, np.SalePrice)
	return err
}

func buildProductCreationHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput, err := validateProductCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		sku := productInput.SKU
		// can't create a product with a sku that already exists!
		exists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product with sku `%s` already exists", sku))
			return
		}

		progenitor := newProductProgenitorFromProductCreationInput(productInput)
		newProgenitor, err := createProductProgenitorInDB(db, progenitor)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		newProduct := &Product{
			ProductProgenitor: *newProgenitor,
			SKU:               productInput.SKU,
			UPC:               NullString{sql.NullString{String: productInput.UPC}},
			Quantity:          productInput.Quantity,
			OnSale:            productInput.OnSale,
			SalePrice:         NullFloat64{sql.NullFloat64{Float64: productInput.SalePrice}},
		}

		err = createProduct(db, newProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		json.NewEncoder(res).Encode(newProduct)
	}
}
