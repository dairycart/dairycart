package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
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
	} else {
		np.AvailableOn = time.Now()
	}
	return np
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	models.Product
	Options []models.ProductOption `json:"options"`
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

func buildProductDeletionHandler(db *sql.DB, client storage.Storer, webhookExecutor WebhookExecutor) http.HandlerFunc {
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

func buildProductUpdateHandler(db *sql.DB, client storage.Storer, webhookExecutor WebhookExecutor) http.HandlerFunc {
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

func buildProductsFromOptions(input *models.ProductCreationInput, createdOptions []models.ProductOption) (toCreate []*models.Product) {
	// lovingly borrowed from:
	//     https://stackoverflow.com/questions/29002724/implement-ruby-style-cartesian-product-in-go
	// NextIndex sets ix to the lexicographically next value,
	// such that for each i > 0, 0 <= ix[i] < lens(i).
	next := func(ix []int, sl [][]optionPlaceholder) {
		for j := len(ix) - 1; j >= 0; j-- {
			ix[j]++
			if j == 0 || ix[j] < len(sl[j]) {
				return
			}
			ix[j] = 0
		}
	}

	// meat & potatoes starts here
	var optionData [][]optionPlaceholder
	for _, o := range createdOptions {
		var newOptions []optionPlaceholder
		for _, v := range o.Values {
			summary := fmt.Sprintf("%s: %s", o.Name, v.Value)
			ph := optionPlaceholder{
				ID:            v.ID,
				Summary:       summary,
				Value:         v.Value,
				OriginalValue: v,
			}
			newOptions = append(newOptions, ph)
		}
		optionData = append(optionData, newOptions)
	}

	for ix := make([]int, len(optionData)); ix[0] < len(optionData[0]); next(ix, optionData) {
		var skuPrefixParts, optionSummaryParts []string
		var originalValues []models.ProductOptionValue
		for j, k := range ix {
			optionSummaryParts = append(optionSummaryParts, optionData[j][k].Summary)
			skuPrefixParts = append(skuPrefixParts, strings.ToLower(optionData[j][k].Value))
			originalValues = append(originalValues, optionData[j][k].OriginalValue)
		}

		productTemplate := newProductFromCreationInput(input)
		productTemplate.OptionSummary = strings.Join(optionSummaryParts, ", ")
		productTemplate.SKU = fmt.Sprintf("%s_%s", input.SKU, strings.Join(skuPrefixParts, "_"))
		productTemplate.ApplicableOptionValues = originalValues
		toCreate = append(toCreate, productTemplate)

	}
	return toCreate
}

func createProductsInDBFromOptionRows(client storage.Storer, tx *sql.Tx, r *models.ProductRoot, input *models.ProductCreationInput, createdOptions []models.ProductOption) ([]models.Product, error) {
	var err error
	createdProducts := []models.Product{}
	productsToCreate := buildProductsFromOptions(input, createdOptions)
	for _, p := range productsToCreate {
		p.ID, p.CreatedOn, p.AvailableOn, err = client.CreateProduct(tx, p)
		if err != nil {
			return nil, err
		}
		p.ProductRootID = r.ID

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

func handleProductCreationImages(tx *sql.Tx, client storage.Storer, imager storage.ImageStorer, images []models.ProductImageCreationInput, sku string, rootID uint64) ([]models.ProductImage, error) {
	returnImages := []models.ProductImage{}
	for i, imageInput := range images {
		var img image.Image
		var err error
		imageType := strings.ToLower(imageInput.Type)

		switch imageType {
		case "base64":
			reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageInput.Data))
			img, _, err = image.Decode(reader)
			if err != nil {
				return nil, fmt.Errorf("Image data at index %d is invalid", i)
			}
		case "url":
			// FIXME: this is almost definitely the wrong way to do this,
			// we should support conversion from known data types (mainly JPEGs) to PNGs
			if !strings.HasSuffix(imageInput.Data, "png") {
				return nil, errors.New("only PNG images are supported")
			}
			response, err := http.Get(imageInput.Data)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error retrieving product image from url %s", imageInput.Data))
			} else {
				defer response.Body.Close()
				img, _, err = image.Decode(response.Body)
				if err != nil {
					return nil, fmt.Errorf("Image data at index %d is invalid", i)
				}
			}
		}

		thumbnails := imager.CreateThumbnails(img)
		locations, err := imager.StoreImages(thumbnails, sku, uint(i))
		if err != nil || locations == nil {
			return nil, err
		}

		newImage := &models.ProductImage{
			ProductRootID: rootID,
			ThumbnailURL:  locations.Thumbnail,
			MainURL:       locations.Main,
			OriginalURL:   locations.Original,
		}

		if imageType == "url" {
			newImage.SourceURL = imageInput.Data
		}

		newImage.ID, newImage.CreatedOn, err = client.CreateProductImage(tx, newImage)
		if err != nil {
			return nil, err
		}

		returnImages = append(returnImages, *newImage)
	}
	return returnImages, nil
}

func buildTestProductCreationHandler(db *sql.DB, client storage.Storer, imager storage.ImageStorer, webhookExecutor WebhookExecutor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		/*
			1. Validate input
			2. Create product root
			3. Create images

			If product has options:
				4a. Create options and values
				5a. create products from the created options and values
				6a. create any necessary images and associate them with created products
			Else:
				4b. create product to associate with product root
				5b. save images and associate them with singular created product
		*/

		// 1. Validate Input
		productInput := &models.ProductCreationInput{}
		err := validateRequestInput(req, productInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		if !restrictedStringIsValid(productInput.SKU) {
			notifyOfInvalidRequestBody(res, fmt.Errorf("The sku received (%s) is invalid", productInput.SKU))
			return
		}
		newProduct := newProductFromCreationInput(productInput)
		newProduct.QuantityPerPackage = uint32(math.Max(float64(newProduct.QuantityPerPackage), 1))
		if productInput.AvailableOn == nil {
			newProduct.AvailableOn = time.Now()
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

		// 2. Create product root
		productRoot := createProductRootFromProduct(newProduct)
		productRoot.ID, productRoot.CreatedOn, err = client.CreateProductRoot(tx, productRoot)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product options and values in database")
			return
		}

		if len(productInput.Options) == 0 {
			// 4b. create product to associate with product root
			newProduct.ProductRootID = productRoot.ID
			newProduct.ID, newProduct.CreatedOn, newProduct.AvailableOn, err = client.CreateProduct(tx, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product in database")
				return
			}

			productRoot.Options = []models.ProductOption{} // so this won't be Marshaled as null
			productRoot.Products = []models.Product{*newProduct}

			// 5b. save images and associate them with singular created product
			newProduct.Images, err = handleProductCreationImages(tx, client, imager, productInput.Images, productInput.SKU, productRoot.ID)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product images in database")
				return
			}

			if newProduct.PrimaryImageID == nil && len(newProduct.Images) > 0 {
				newProduct.PrimaryImageID = &newProduct.Images[0].ID
			}

			_, err = client.SetPrimaryProductImageForProduct(tx, newProduct.ID, *newProduct.PrimaryImageID)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "set primary image ID")
				return
			}
		} else {
			// 4a. Create options and values
			for _, optionAndValues := range productInput.Options {
				o, err := createProductOptionAndValuesInDBFromInput(tx, optionAndValues, productRoot.ID, client)
				if err != nil {
					tx.Rollback()
					notifyOfInternalIssue(res, err, "insert product options and values in database")
					return
				}
				productRoot.Options = append(productRoot.Options, o)
			}

			// 5a. create products from the created options and values
			productRoot.Products, err = createProductsInDBFromOptionRows(client, tx, productRoot, productInput, productRoot.Options)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert products in database")
				return
			}

			// 6a. create any necessary images and associate them with created products
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

func buildProductCreationHandler(db *sql.DB, client storage.Storer, webhookExecutor WebhookExecutor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput := &models.ProductCreationInput{}
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
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product options and values in database")
			return
		}

		newProduct.QuantityPerPackage = uint32(math.Max(float64(newProduct.QuantityPerPackage), 1))

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
			newProduct.ID, newProduct.CreatedOn, newProduct.AvailableOn, err = client.CreateProduct(tx, newProduct)
			if err != nil {
				tx.Rollback()
				notifyOfInternalIssue(res, err, "insert product in database")
				return
			}

			productRoot.Options = []models.ProductOption{} // so this won't be Marshaled as null
			productRoot.Products = []models.Product{*newProduct}
		} else {
			productRoot.Products, err = createProductsInDBFromOptionRows(client, tx, productRoot, productInput, productRoot.Options)
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
