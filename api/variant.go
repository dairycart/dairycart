package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

// Variant describes a product's variation
type Variant struct {
	// Basic Info
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	UPC       string `json:"upc"`
	Quantity  int    `json:"quantity"`

	// Pricing Fields
	OnSale    bool    `json:"on_sale"`
	Price     float32 `json:"price"`
	SalePrice float32 `json:"sale_price"`

	// Housekeeping
	CreatedAt  time.Time `json:"created"`
	ArchivedAt NullTime  `json:"-"`
}

func filterActiveProducts(req *http.Request, q *orm.Query) {
	if req.URL.Query().Get("include_archived") != "true" {
		q.Where("product.archived_at is null")
	}
}

// productExistsInDB will return whether or not a product with a given sku exists in the database
func productExistsInDB(db *sql.DB, sku string) (bool, error) {
	var exists string
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1);", sku).Scan(&exists)
	if err != nil {
		log.Printf("error encountered querying for shit: %v", err)
	}
	return exists == "true", err
}

func buildProductExistenceHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		productExists, err := productExistsInDB(db, sku)
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
}

// retrieveProductFromDB retrieves a product with a given SKU from the database
func retrieveProductFromDB(db *pg.DB, sku string) (*Product, error) {
	product := &Product{}
	err := db.Model(product).
		Where("sku = ?", sku).
		Where("product.archived_at is null").
		Select()

	return product, err
}

func buildSingleProductHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]
		product, err := retrieveProductFromDB(db, sku)

		if err != nil {
			informOfServerIssue(err, "Error encountered querying for product", res)
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func buildProductListHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
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
}

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

func buildProductDeletionHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		product := &Product{}
		err := db.Model(product).Where("sku = ?", sku).Where("archived_at is null").Select()

		if err != nil {
			informOfServerIssue(err, "Error deleting product from database", res)
			return
		}

		db.Model(product).Set("archived_at = now()").Where("sku = ?", sku).Update(product)
		json.NewEncoder(res).Encode(product)
	}
}
