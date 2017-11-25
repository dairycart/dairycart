package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
)

// newProductFromCreationInput creates a new product from a ProductCreationInput
func newProductFromCreationInput(in *ProductCreationInput) *models.Product {
	np := &models.Product{
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
	models.Product
	Options []*ProductOptionCreationInput `json:"options"`
}

func buildProductExistenceHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		productExists, err := client.ProductWithSKUExists(db, sku)
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

func buildSingleProductHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		product, err := client.GetProductBySKU(db, sku)
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

func buildProductListHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := client.GetProductCount(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of products from the database")
			return
		}

		products, err := client.GetProductList(db, queryFilter)
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

func buildProductDeletionHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		// can't delete a product that doesn't exist!
		existingProduct, err := client.GetProductBySKU(db, sku)
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

		_, err = client.DeleteProductVariantBridgeByProductID(tx, existingProduct.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product in database")
			return
		}

		archiveTime, err := client.DeleteProduct(tx, existingProduct.ID)
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
		existingProduct.ArchivedOn = models.NullTime{NullTime: pq.NullTime{Time: archiveTime, Valid: true}}
		json.NewEncoder(res).Encode(existingProduct)
	}
}

func buildProductUpdateHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductUpdateHandler is a request handler that can update products
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		newerProduct := &models.Product{}
		err := validateRequestInput(req, newerProduct)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingProduct, err := client.GetProductBySKU(db, sku)
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

		updatedTime, err := client.UpdateProduct(db, newerProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}
		newerProduct.UpdatedOn = models.NullTime{NullTime: pq.NullTime{Time: updatedTime, Valid: true}}

		json.NewEncoder(res).Encode(newerProduct)
	}
}

func createProductsInDBFromOptionRows(client storage.Storer, tx *sql.Tx, r *models.ProductRoot, np *models.Product) ([]models.Product, error) {
	createdProducts := []models.Product{}
	productOptionData := generateCartesianProductForOptions(r.Options)
	for _, option := range productOptionData {
		p := &models.Product{}
		*p = *np // solved: http://www.claymath.org/millennium-problems/p-vs-np-problem

		p.ProductRootID = r.ID
		p.ApplicableOptionValues = option.OriginalValues
		p.OptionSummary = option.OptionSummary
		p.SKU = fmt.Sprintf("%s_%s", r.SKUPrefix, option.SKUPostfix)

		var err error
		p.ID, p.AvailableOn, p.CreatedOn, err = client.CreateProduct(tx, p)
		if err != nil {
			return nil, err
		}

		err = client.CreateMultipleProductVariantBridgesForProductID(tx, p.ID, option.IDs)
		// err = createBridgeEntryForProductValues(tx, p.ID, option.IDs)
		if err != nil {
			return nil, err
		}
		createdProducts = append(createdProducts, *p)
	}
	return createdProducts, nil
}

func buildProductCreationHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
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
		exists, err := client.ProductRootWithSKUPrefixExists(db, productInput.SKU)
		// exists, err := rowExistsInDB(db, productRootSkuExistenceQuery, productInput.SKU)
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
		productRoot.ID, productRoot.CreatedOn, err = client.CreateProductRoot(tx, productRoot)
		// productRoot.ID, productRoot.CreatedOn, err = createProductRootInDB(tx, productRoot)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product options and values in database")
			return
		}

		for _, optionAndValues := range productInput.Options {
			o, err := createProductOptionAndValuesInDBFromInput(tx, optionAndValues, productRoot.ID, client)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product options and values in database")
				return
			}
			productRoot.Options = append(productRoot.Options, o)
		}

		if len(productInput.Options) == 0 {
			newProduct.ProductRootID = productRoot.ID
			newProduct.ID, newProduct.AvailableOn, newProduct.CreatedOn, err = client.CreateProduct(tx, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product in database")
				return
			}
			productRoot.Options = []models.ProductOption{} // so this won't be Marshaled as null
			productRoot.Products = []models.Product{*newProduct}
		} else {
			productRoot.Products, err = createProductsInDBFromOptionRows(client, tx, productRoot, newProduct)
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
