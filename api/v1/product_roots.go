package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"

	"github.com/go-chi/chi"
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

func buildProductRootListHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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

func buildSingleProductRootHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	// SingleProductRootHandler is a request handler that returns a single product root
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		// eating this error because the router should have ensured this is an integer
		productRootID, _ := strconv.ParseUint(productRootIDStr, 10, 64)

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

func buildProductRootDeletionHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		// eating this error because the router should have ensured this is an integer
		productRootID, _ := strconv.ParseUint(productRootIDStr, 10, 64)

		// can't delete a product root that doesn't exist!
		productRoot, err := client.GetProductRoot(db, productRootID)
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
		_, err = client.ArchiveProductVariantBridgesWithProductRootID(tx, productRoot.ID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product variant bridge entries in database")
			return
		}

		// delete product option values
		_, err = client.ArchiveProductOptionValuesWithProductRootID(tx, productRoot.ID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product option values in database")
			return
		}

		// delete product options
		_, err = client.ArchiveProductOptionsWithProductRootID(tx, productRoot.ID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product options in database")
			return
		}

		// delete products
		_, err = client.ArchiveProductsWithProductRootID(tx, productRoot.ID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive products in database")
			return
		}

		// delete the actual product root
		archivedOn, err := client.DeleteProductRoot(tx, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product root in database")
			return
		}
		productRoot.ArchivedOn = &models.Dairytime{Time: archivedOn}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		json.NewEncoder(res).Encode(productRoot)
	}
}
