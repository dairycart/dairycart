package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

const (
	productRootSkuExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_roots WHERE sku_prefix = $1 AND archived_on IS NULL)`
	productRootExistenceQuery    = `SELECT EXISTS(SELECT 1 FROM product_roots WHERE id = $1 AND archived_on IS NULL)`
	productRootRetrievalQuery    = `SELECT id, name, subtitle, description, sku_prefix, manufacturer, brand, taxable, cost, product_weight, product_height, product_width, product_length, package_weight, package_height, package_width, package_length, quantity_per_package, available_on, created_on, updated_on, archived_on FROM product_roots WHERE id = $1`
	productRootDeletionQuery     = `UPDATE product_roots SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`

	productDeletionQueryByRootID              = `UPDATE products SET archived_on = NOW() WHERE product_root_id = $1 AND archived_on IS NULL`
	productOptionDeletionQueryByRootID        = `UPDATE product_options SET archived_on = NOW() WHERE product_root_id = $1 AND archived_on IS NULL`
	productOptionValueDeletionQueryByRootID   = `UPDATE product_option_values SET archived_on = NOW() WHERE product_option_id IN (SELECT id FROM product_options WHERE product_root_id = $1)`
	productVariantBridgeDeletionQueryByRootID = `UPDATE product_variant_bridge SET archived_on = NOW() WHERE product_id IN (SELECT id FROM products WHERE product_root_id = $1)`
)

// ProductRoot represents the object that products inherit from
type ProductRoot struct {
	DBRow

	// Basic Info
	Name               string     `json:"name"`
	Subtitle           NullString `json:"subtitle"`
	Description        string     `json:"description"`
	SKUPrefix          string     `json:"sku_prefix"`
	Manufacturer       NullString `json:"manufacturer"`
	Brand              NullString `json:"brand"`
	AvailableOn        time.Time  `json:"available_on"`
	QuantityPerPackage uint32     `json:"quantity_per_package"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Cost    float32 `json:"cost"`

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

	Options  []*ProductOption `json:"options"`
	Products []Product        `json:"products"`
}

func createProductRootFromProduct(p *Product) *ProductRoot {
	r := &ProductRoot{
		Name:               p.Name,
		Subtitle:           p.Subtitle,
		Description:        p.Description,
		SKUPrefix:          p.SKU,
		Manufacturer:       p.Manufacturer,
		Brand:              p.Brand,
		QuantityPerPackage: p.QuantityPerPackage,
		Taxable:            p.Taxable,
		Cost:               p.Cost,
		ProductWeight:      p.ProductWeight,
		ProductHeight:      p.ProductHeight,
		ProductWidth:       p.ProductWidth,
		ProductLength:      p.ProductLength,
		PackageWeight:      p.PackageWeight,
		PackageHeight:      p.PackageHeight,
		PackageWidth:       p.PackageWidth,
		PackageLength:      p.PackageLength,
		AvailableOn:        p.AvailableOn,
	}
	return r
}

func createProductRootInDB(tx *sql.Tx, r *ProductRoot) (uint64, time.Time, error) {
	var newRootID uint64
	var createdOn time.Time
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	query, queryArgs := buildProductRootCreationQuery(r)
	err := tx.QueryRow(query, queryArgs...).Scan(&newRootID, &createdOn)

	return newRootID, createdOn, err
}

// retrieveProductRootFromDB retrieves a product root with a given ID from the database
func retrieveProductRootFromDB(db *sqlx.DB, id uint64) (ProductRoot, error) {
	var root ProductRoot
	err := db.QueryRowx(productRootRetrievalQuery, id).StructScan(&root)
	return root, err
}

func buildProductRootListHandler(db *sqlx.DB) http.HandlerFunc {
	// productRootListHandler is a request handler that returns a list of product roots
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := getRowCount(db, "product_roots", queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of product roots from the database")
			return
		}

		var productRoots []ProductRoot
		query, args := buildProductRootListQuery(queryFilter)
		err = retrieveListOfRowsFromDB(db, query, args, &productRoots)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product roots from the database")
			return
		}

		for _, r := range productRoots {
			query, args := buildProductAssociatedWithRootListQuery(r.ID)
			err = retrieveListOfRowsFromDB(db, query, args, &r.Products)
			if err != nil {
				notifyOfInternalIssue(res, err, "retrieve products from the database")
				return
			}
		}

		productsResponse := &ListResponse{
			Page:  queryFilter.Page,
			Limit: queryFilter.Limit,
			Count: count,
			Data:  productRoots,
		}
		json.NewEncoder(res).Encode(productsResponse)
	}
}

func buildSingleProductRootHandler(db *sqlx.DB) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		productRootID, err := strconv.ParseUint(productRootIDStr, 10, 64)

		productRoot, err := retrieveProductRootFromDB(db, productRootID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product_root", productRootIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving product root from database")
			return
		}

		query, args := buildProductAssociatedWithRootListQuery(productRoot.ID)
		err = retrieveListOfRowsFromDB(db, query, args, &productRoot.Products)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		productRoot.Options, err = getProductOptionsForProductRoot(db, productRoot.ID, nil)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product options from the database")
			return
		}

		json.NewEncoder(res).Encode(productRoot)
	}
}

func deleteProductRoot(tx *sql.Tx, rootID uint64) error {
	_, err := tx.Exec(productRootDeletionQuery, rootID)
	return err
}

func deleteProductsAssociatedWithRoot(tx *sql.Tx, rootID uint64) error {
	_, err := tx.Exec(productDeletionQueryByRootID, rootID)
	return err
}

func deleteProductOptionsAssociatedWithRoot(tx *sql.Tx, rootID uint64) error {
	_, err := tx.Exec(productOptionDeletionQueryByRootID, rootID)
	return err
}

func deleteProductOptionValuesAssociatedWithRoot(tx *sql.Tx, rootID uint64) error {
	_, err := tx.Exec(productOptionValueDeletionQueryByRootID, rootID)
	return err
}

func deleteVariantBridgeEntriesAssociatedWithRoot(tx *sql.Tx, rootID uint64) error {
	_, err := tx.Exec(productVariantBridgeDeletionQueryByRootID, rootID)
	return err
}

func buildProductRootDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		productRootID, err := strconv.ParseUint(productRootIDStr, 10, 64)

		// can't delete a product root that doesn't exist!
		productRoot, err := retrieveProductRootFromDB(db, productRootID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product_root", productRootIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving product root from database")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "create new database transaction")
			return
		}

		// delete product variant bridge entries
		err = deleteVariantBridgeEntriesAssociatedWithRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product variant bridge entries in database")
			return
		}

		// delete product option values
		err = deleteProductOptionValuesAssociatedWithRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product option values in database")
			return
		}

		// delete product options
		err = deleteProductOptionsAssociatedWithRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product options in database")
			return
		}

		// delete products
		err = deleteProductsAssociatedWithRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive products in database")
			return
		}

		// delete the actual product root
		err = deleteProductRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product root in database")
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
