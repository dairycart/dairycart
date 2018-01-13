package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/storage/images"
	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	badSKUUpdateJSON = `{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`
)

func TestNewProductFromCreationInput(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		_t.Parallel()
		exampleInput := &models.ProductCreationInput{
			Name:               "name",
			Subtitle:           "sub",
			Description:        "desc",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "mfgr",
			Brand:              "brand",
			Quantity:           123,
			QuantityPerPackage: 123,
			Taxable:            true,
			Price:              12.34,
			OnSale:             true,
			SalePrice:          12.34,
			Cost:               12.34,
			ProductWeight:      123,
			ProductHeight:      123,
			ProductWidth:       123,
			ProductLength:      123,
			PackageWeight:      123,
			PackageHeight:      123,
			PackageWidth:       123,
			PackageLength:      123,
		}

		expected := &models.Product{
			Name:               "name",
			Subtitle:           "sub",
			Description:        "desc",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "mfgr",
			Brand:              "brand",
			Quantity:           123,
			QuantityPerPackage: 123,
			Taxable:            true,
			Price:              12.34,
			OnSale:             true,
			SalePrice:          12.34,
			Cost:               12.34,
			ProductWeight:      123,
			ProductHeight:      123,
			ProductWidth:       123,
			ProductLength:      123,
			PackageWeight:      123,
			PackageHeight:      123,
			PackageWidth:       123,
			PackageLength:      123,
		}

		actual := newProductFromCreationInput(exampleInput)
		assert.Equal(t, expected, actual)
	})

	t.Run("with available on value", func(_t *testing.T) {
		_t.Parallel()
		exampleInput := &models.ProductCreationInput{
			AvailableOn: buildTestDairytime(),
		}

		expected := &models.Product{
			AvailableOn: buildTestTime(),
		}

		actual := newProductFromCreationInput(exampleInput)
		assert.Equal(t, expected, actual)
	})
}

func TestCreateProductsInDBFromOptions(t *testing.T) {
	exampleID := uint64(1)
	exampleProductRoot := &models.ProductRoot{
		Options: []models.ProductOption{
			{
				ID:            exampleID,
				Name:          "name",
				ProductRootID: exampleID,
				CreatedOn:     buildTestTime(),
				Values:        []models.ProductOptionValue{{Value: "one"}, {Value: "two"}, {Value: "three"}},
			},
		},
	}
	exampleProductCreationInput := &models.ProductCreationInput{
		Options: []models.ProductOptionCreationInput{
			{
				Name:   "name",
				Values: []string{"one", "two", "three"},
			},
		},
	}
	exampleCreatedProductOptions := []models.ProductOption{
		{
			Name: "name",
			Values: []models.ProductOptionValue{
				{
					Value: "one",
				},
				{
					Value: "two",
				},
				{
					Value: "three",
				},
			},
		},
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actual, err := createProductsInDBFromOptions(testUtil.MockDB, tx, exampleProductRoot, exampleProductCreationInput, exampleCreatedProductOptions)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("with error creating products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), buildTestTime(), generateArbitraryError())

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		_, err = createProductsInDBFromOptions(testUtil.MockDB, tx, exampleProductRoot, exampleProductCreationInput, exampleCreatedProductOptions)
		assert.NotNil(t, err)
	})

}

func TestProductExistenceHandler(t *testing.T) {
	exampleSKU := "example"
	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(true, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
	})
	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
	t.Run("with error performing check", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestProductRetrievalHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     buildTestTime(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   buildTestTime(),
	}
	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with DB error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestProductListHandler(t *testing.T) {
	exampleProduct := models.Product{
		ID:            2,
		CreatedOn:     buildTestTime(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   buildTestTime(),
	}
	exampleLength := uint64(3)

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).
			Return(exampleLength, nil)
		testUtil.MockDB.On("GetProductList", mock.Anything, mock.Anything).
			Return([]models.Product{exampleProduct, exampleProduct, exampleProduct}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error retrieving count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).Return(exampleLength, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with database error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).Return(exampleLength, nil)
		testUtil.MockDB.On("GetProductList", mock.Anything, mock.Anything).
			Return([]models.Product{exampleProduct, exampleProduct, exampleProduct}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductUpdateHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     buildTestTime(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   buildTestTime(),
	}
	exampleProductUpdateInput := `
		{
			"sku": "example",
			"name": "Test",
			"quantity": 666,
			"upc": "1234567890",
			"price": 12.34
		}
	`
	exampleWebhook := models.Webhook{
		URL:         "https://dairycart.com",
		ContentType: "application/json",
	}

	t.Run("normal operation", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductUpdatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/v1/product/%s", exampleProduct.SKU),
			strings.NewReader(exampleProductUpdateInput),
		)
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
		ensureExpectationsWereMet(t, testUtil.Mock)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, sql.ErrNoRows).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/v1/product/%s", exampleProduct.SKU),
			strings.NewReader(exampleProductUpdateInput),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with database error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch,
			fmt.Sprintf("/v1/product/%s", exampleProduct.SKU),
			strings.NewReader(exampleProductUpdateInput),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with input validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/v1/product/example",
			strings.NewReader(exampleGarbageInput),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with SKU validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/v1/product/skateboard",
			strings.NewReader(badSKUUpdateJSON),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with database error updating product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).
			Return(buildTestTime(), generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/v1/product/%s", exampleProduct.SKU),
			strings.NewReader(exampleProductUpdateInput),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving webhooks", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductUpdatedWebhookEvent).
			Return([]models.Webhook{}, generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/v1/product/%s", exampleProduct.SKU),
			strings.NewReader(exampleProductUpdateInput),
		)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductDeletionHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:        2,
		CreatedOn: buildTestTime(),
		SKU:       exampleSKU,
		Name:      "Skateboard",
	}
	exampleWebhook := models.Webhook{
		URL:         "https://dairycart.com",
		ContentType: "application/json",
	}

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductArchivedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, sql.ErrNoRows).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error beginning transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting bridge entries", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error encountered deleting product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).
			Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).
			Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), generateArbitraryError()).
			Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(buildTestTime(), nil).Once()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductArchivedWebhookEvent).
			Return([]models.Webhook{}, generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductCreationHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     buildTestTime(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	expectedFirstOption := &models.Product{}
	expectedSecondOption := &models.Product{}
	expectedThirdOption := &models.Product{}

	*expectedFirstOption = *exampleProduct
	*expectedSecondOption = *exampleProduct
	*expectedThirdOption = *exampleProduct

	expectedFirstOption.OptionSummary = "something: one"
	expectedFirstOption.SKU = "skateboard_one"
	expectedSecondOption.OptionSummary = "something: two"
	expectedSecondOption.SKU = "skateboard_two"
	expectedThirdOption.OptionSummary = "something: three"
	expectedThirdOption.SKU = "skateboard_three"

	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &models.ProductOption{
		ID:            exampleProductOption.ID,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []models.ProductOptionValue{
			{
				ID:              128, // == exampleProductOptionValue.ID,
				CreatedOn:       buildTestTime(),
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				ID:              256, // == exampleProductOptionValue.ID,
				CreatedOn:       buildTestTime(),
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				ID:              512, // == exampleProductOptionValue.ID,
				CreatedOn:       buildTestTime(),
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

	exampleProductCreationInput := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1
		}
	`

	exampleProductCreationInputWithImages := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"images": [
				{
					"type": "base64",
					"is_primary": true,
					"data": "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKAQMAAAC3/F3+AAAABlBMVEUA/wAA/wD8J4MxAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAAC0lEQVQImWNgwAcAAB4AAe72cCEAAAAASUVORK5CYII="
				}
			]
		}
	`

	exampleProductCreationInputWithNoPrimaryImages := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"images": [
				{
					"type": "base64",
					"data": "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKAQMAAAC3/F3+AAAABlBMVEUA/wAA/wD8J4MxAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAAC0lEQVQImWNgwAcAAB4AAe72cCEAAAAASUVORK5CYII="
				}
			]
		}
	`

	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": [{
				"name": "something",
				"values": [
					"one",
					"two",
					"three"
				]
			}]
		}
	`

	exampleWebhook := models.Webhook{
		URL:         "https://dairycart.com",
		ContentType: "application/json",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("without options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with error validating input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(false, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with invalid product sku", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(false, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(badSKUUpdateJSON))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with already existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(true, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("when transaction fails to begin", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with issue handling product images", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.Mock.ExpectRollback()

		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(images.ProductImageSet{})
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(&images.ProductImageLocations{}, generateArbitraryError())

		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithImages))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("where product creation fails", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("primary image ID gets set", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()

		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(images.ProductImageSet{})
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(&images.ProductImageLocations{}, nil)
		testUtil.MockDB.On("CreateProductImage", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), nil)
		testUtil.MockDB.On("SetPrimaryProductImageForProduct", mock.Anything, mock.AnythingOfType("uint64"), mock.Anything).
			Return(buildTestTime(), nil)

		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithNoPrimaryImages))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("primary image ID gets set when specified", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()

		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(images.ProductImageSet{})
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(&images.ProductImageLocations{}, nil)
		testUtil.MockDB.On("CreateProductImage", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), nil)
		testUtil.MockDB.On("SetPrimaryProductImageForProduct", mock.Anything, mock.AnythingOfType("uint64"), mock.Anything).
			Return(buildTestTime(), nil)

		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithImages))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with failure to set primary image for product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{exampleWebhook}, nil).Once()

		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(images.ProductImageSet{})
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(&images.ProductImageLocations{}, nil)
		testUtil.MockDB.On("CreateProductImage", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), nil)
		testUtil.MockDB.On("SetPrimaryProductImageForProduct", mock.Anything, mock.AnythingOfType("uint64"), mock.Anything).
			Return(buildTestTime(), generateArbitraryError())

		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithNoPrimaryImages))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), generateArbitraryError())
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating option products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), generateArbitraryError())
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating bridge entries", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("when commit returns an error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving webhooks", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, buildTestTime(), buildTestTime(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetWebhooksByEventType", mock.Anything, ProductCreatedWebhookEvent).
			Return([]models.Webhook{}, generateArbitraryError()).Once()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.NoError(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
