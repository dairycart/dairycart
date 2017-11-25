package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

const (
	productRootExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_roots WHERE id = $1 AND archived_on IS NULL)`
	productRootRetrievalQuery = `SELECT id, name, subtitle, description, sku_prefix, manufacturer, brand, taxable, cost, product_weight, product_height, product_width, product_length, package_weight, package_height, package_width, package_length, quantity_per_package, available_on, created_on, updated_on, archived_on FROM product_roots WHERE id = $1`
	productRootDeletionQuery  = `UPDATE product_roots SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`

	productDeletionQueryByRootID              = `UPDATE products SET archived_on = NOW() WHERE product_root_id = $1 AND archived_on IS NULL`
	productOptionDeletionQueryByRootID        = `UPDATE product_options SET archived_on = NOW() WHERE product_root_id = $1 AND archived_on IS NULL`
	productOptionValueDeletionQueryByRootID   = `UPDATE product_option_values SET archived_on = NOW() WHERE product_option_id IN (SELECT id FROM product_options WHERE product_root_id = $1)`
	productVariantBridgeDeletionQueryByRootID = `UPDATE product_variant_bridge SET archived_on = NOW() WHERE product_id IN (SELECT id FROM products WHERE product_root_id = $1)`
)

func createProductRootFromProduct(p *models.Product) *models.ProductRoot {
	r := &models.ProductRoot{
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

// retrieveProductRootFromDB retrieves a product root with a given ID from the database
func retrieveProductRootFromDB(db *sqlx.DB, id uint64) (models.ProductRoot, error) {
	var root models.ProductRoot
	err := db.QueryRowx(productRootRetrievalQuery, id).StructScan(&root)
	return root, err
}

func buildProductRootListHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// productRootListHandler is a request handler that returns a list of product roots
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := client.GetProductRootCount(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of product roots from the database")
			return
		}

		productRoots, err := client.GetProductRootList(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product roots from the database")
			return
		}

		for _, r := range productRoots {
			products, err := client.GetProductsByProductRootID(db, r.ID)
			if err != nil {
				notifyOfInternalIssue(res, err, "retrieve products from the database")
				return
			}
			r.Products = products
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

func buildSingleProductRootHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// SingleProductRootHandler is a request handler that returns a single product root
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		productRootID, err := strconv.ParseUint(productRootIDStr, 10, 64)

		productRoot, err := client.GetProductRoot(db, productRootID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product_root", productRootIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving product root from database")
			return
		}

		products, err := client.GetProductsByProductRootID(db, productRoot.ID)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}
		productRoot.Products = products

		options, err := client.GetProductOptionsByProductRootID(db, productRoot.ID)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product options from the database")
			return
		}
		productRoot.Options = options

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
