package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairycart/storage/images"
	"github.com/dairycart/dairymodels/v1"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
)

const (
	ProductCreatedWebhookEvent  = "product_created"
	ProductUpdatedWebhookEvent  = "product_updated"
	ProductArchivedWebhookEvent = "product_archived"
)

// newProductFromCreationInput creates a new product from a ProductCreationInput
func newProductFromCreationInput(in *models.ProductCreationInput) *models.Product {
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
	}
	if in.AvailableOn != nil {
		np.AvailableOn = in.AvailableOn.Time
	}
	return np
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	models.Product
	Options []models.ProductOption `json:"options"`
}

func buildProductExistenceHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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

func buildSingleProductHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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

func buildProductListHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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

func buildProductDeletionHandler(db *sql.DB, client database.Storer, webhookExecutor WebhookExecutor) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		// can't delete a product that doesn't exist!
		product, err := client.GetProductBySKU(db, sku)
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

		_, err = client.DeleteProductVariantBridgeByProductID(tx, product.ID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "archive product variant bridges in database")
			return
		}

		archiveTime, err := client.DeleteProduct(tx, product.ID)
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
		product.ArchivedOn = &models.Dairytime{Time: archiveTime}

		webhooks, err := client.GetWebhooksByEventType(db, ProductArchivedWebhookEvent)
		if err != nil && err != sql.ErrNoRows {
			notifyOfInternalIssue(res, err, "retrieve webhooks from database")
			return
		}

		for _, wh := range webhooks {
			go webhookExecutor.CallWebhook(wh, product, db, client)
		}

		json.NewEncoder(res).Encode(product)
	}
}

func buildProductUpdateHandler(db *sql.DB, client database.Storer, webhookExecutor WebhookExecutor) http.HandlerFunc {
	// ProductUpdateHandler is a request handler that can update products
	return func(res http.ResponseWriter, req *http.Request) {
		sku := chi.URLParam(req, "sku")

		updatedProduct := &models.Product{}
		err := validateRequestInput(req, updatedProduct)
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

		mergo.Merge(updatedProduct, existingProduct)

		if !restrictedStringIsValid(updatedProduct.SKU) {
			notifyOfInvalidRequestBody(res, fmt.Errorf("The sku received (%s) is invalid", updatedProduct.SKU))
			return
		}

		updatedTime, err := client.UpdateProduct(db, updatedProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}
		updatedProduct.UpdatedOn = &models.Dairytime{Time: updatedTime}

		webhooks, err := client.GetWebhooksByEventType(db, ProductUpdatedWebhookEvent)
		if err != nil && err != sql.ErrNoRows {
			notifyOfInternalIssue(res, err, "retrieve webhooks from database")
			return
		}

		for _, wh := range webhooks {
			go webhookExecutor.CallWebhook(wh, updatedProduct, db, client)
		}

		json.NewEncoder(res).Encode(updatedProduct)
	}
}

func createProductsInDBFromOptions(client database.Storer, tx *sql.Tx, r *models.ProductRoot, input *models.ProductCreationInput, createdOptions []models.ProductOption) ([]models.Product, error) {
	var err error
	createdProducts := []models.Product{}
	productsToCreate := buildProductsFromOptions(input, createdOptions)
	for _, p := range productsToCreate {
		p.ProductRootID = r.ID
		p.ID, p.CreatedOn, p.AvailableOn, err = client.CreateProduct(tx, p)
		if err != nil {
			return nil, err
		}

		optionIDs := []uint64{}
		for _, o := range p.ApplicableOptionValues {
			optionIDs = append(optionIDs, o.ID)
		}

		err = client.CreateMultipleProductVariantBridgesForProductID(tx, p.ID, optionIDs)
		if err != nil {
			return nil, err
		}
		createdProducts = append(createdProducts, *p)
	}
	return createdProducts, nil
}

func buildProductCreationHandler(db *sql.DB, client database.Storer, imager images.ImageStorer, webhookExecutor WebhookExecutor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput := &models.ProductCreationInput{}
		err := validateRequestInput(req, productInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		if !restrictedStringIsValid(productInput.SKU) {
			notifyOfInvalidRequestBody(res, fmt.Errorf("the sku received (%s) is invalid", productInput.SKU))
			return
		}

		newProduct := newProductFromCreationInput(productInput)
		newProduct.QuantityPerPackage = uint32(math.Max(float64(newProduct.QuantityPerPackage), 1))
		if productInput.AvailableOn == nil {
			newProduct.AvailableOn = time.Now()
		}

		// can't create a product with a sku that already exists!
		exists, err := client.ProductRootWithSKUPrefixExists(db, productInput.SKU)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product with sku '%s' already exists", productInput.SKU))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "create new database transaction")
			return
		}

		productRoot := createProductRootFromProduct(newProduct)
		productRoot.ID, productRoot.CreatedOn, err = client.CreateProductRoot(tx, productRoot)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product options and values in database")
			return
		}

		productRoot.Images, productRoot.PrimaryImageID, err = handleProductCreationImages(tx, client, imager, productInput.Images, productInput.SKU, productRoot.ID)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product images in database")
			return
		}

		if len(productInput.Options) == 0 {
			newProduct.ProductRootID = productRoot.ID
			newProduct.ID, newProduct.CreatedOn, newProduct.AvailableOn, err = client.CreateProduct(tx, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product in database")
				return
			}

			if productRoot.PrimaryImageID != nil {
				newProduct.PrimaryImageID = productRoot.PrimaryImageID
			} else if newProduct.PrimaryImageID == nil && len(productRoot.Images) > 0 {
				productRoot.PrimaryImageID = &productRoot.Images[0].ID
				newProduct.PrimaryImageID = &productRoot.Images[0].ID
			}
			productRoot.Products = []models.Product{*newProduct}

			if len(productRoot.Images) > 0 {
				_, err = client.SetPrimaryProductImageForProduct(tx, newProduct.ID, *newProduct.PrimaryImageID)
				if err != nil {
					tx.Rollback()
					notifyOfInternalIssue(res, err, "set primary image ID")
					return
				}
			}
		} else {
			for _, optionAndValues := range productInput.Options {
				o, err := createProductOptionAndValuesInDBFromInput(tx, optionAndValues, productRoot.ID, client)
				if err != nil {
					tx.Rollback()
					notifyOfInternalIssue(res, err, "insert product options and values in database")
					return
				}
				productRoot.Options = append(productRoot.Options, o)
			}

			productRoot.Products, err = createProductsInDBFromOptions(client, tx, productRoot, productInput, productRoot.Options)
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

		webhooks, err := client.GetWebhooksByEventType(db, ProductCreatedWebhookEvent)
		if err != nil && err != sql.ErrNoRows {
			notifyOfInternalIssue(res, err, "retrieve webhooks from database")
			return
		}

		for _, wh := range webhooks {
			go webhookExecutor.CallWebhook(wh, productRoot, db, client)
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(productRoot)
	}
}
