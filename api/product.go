package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

const (
	skuRetrievalQuery      = "SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.sku = $1 AND p.archived_at IS NULL;"
	productsRetrievalQuery = "SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.id IS NOT NULL AND p.archived_at IS NULL;"
)

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
	CreatedAt  time.Time `json:"created"`
	ArchivedAt NullTime  `json:"-"`
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

// productExistsInDB will return whether or not a product with a given sku exists in the database
func productExistsInDB(db *sql.DB, sku string) (bool, error) {
	var exists string

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 and archived_at is null);", sku).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, errors.Wrap(err, "Error querying for product")
	}

	return exists == "true", err
}

func buildProductExistenceHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		productExists, err := productExistsInDB(db, sku)
		if err != nil {
			log.Printf(`informing user that the product they were looking for (sku %s) does not exist`, sku)
			http.NotFound(res, req)
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

	err := db.QueryRow(skuRetrievalQuery, sku).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return product, errors.Wrap(err, "Error querying for product")
	}

	return product, nil
}

func buildSingleProductHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		product, err := retrieveProductFromDB(db, sku)
		if err != nil {
			log.Printf(`informing user that the product they were looking for (sku %s) does not exist`, sku)
			http.NotFound(res, req)
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func retrieveProductsFromDB(db *sql.DB) ([]Product, error) {
	var products []Product

	rows, err := db.Query(productsRetrievalQuery)
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

func buildProductListHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		products, err := retrieveProductsFromDB(db)
		if err != nil {
			informOfServerIssue(err, "Error encountered querying for product", res)
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

func buildProductDeletionHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]
		// can't delete a product that doesn't exist!
		productExists, err := productExistsInDB(db, sku)
		if err != nil {
			informOfServerIssue(err, "Error encountered querying for product", res)
			return
		}

		if !productExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = db.Exec("UPDATE products SET archived_at = NOW() WHERE sku = $1 and archived_at is null", sku)
		if err != nil {
			informOfServerIssue(err, "Error deleting product from database", res)
			return
		}

		json.NewEncoder(res).Encode("OK")
	}
}

// Unmigrated functions start here

func buildProductUpdateHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductUpdateHandler is a request handler that can update products
		sku := mux.Vars(req)["sku"]
		existingProduct := &Product{}
		existingProductQuery := db.Model(existingProduct).Where("sku = ?", sku).Where("archived_at is null")

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
}

// createProduct takes a marshalled Product object and creates an entry for it and a base_product in the database
func createProduct(db *pg.DB, newProduct *Product) error {
	err := db.Insert(newProduct)
	if err != nil {
		return err
	}
	return err
}

func buildProductCreationHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		newProduct := &Product{}
		bodyIsInvalid := ensureRequestBodyValidity(res, req, newProduct)
		if bodyIsInvalid {
			return
		}

		err := createProduct(db, newProduct)
		if err != nil {
			informOfServerIssue(err, "Error inserting product into database", res)
			return
		}
	}
}
