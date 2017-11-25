package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"
	// "time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	badSKUUpdateJSON = `{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`
)

func TestCreateProductsInDBFromOptionRows(t *testing.T) {
	exampleID := uint64(1)
	exampleProductRoot := &models.ProductRoot{
		Options: []models.ProductOption{
			{
				ID:            exampleID,
				Name:          "name",
				ProductRootID: exampleID,
				CreatedOn:     generateExampleTimeForTests(),
				Values:        []models.ProductOptionValue{{Value: "one"}, {Value: "two"}, {Value: "three"}},
			},
		},
	}
	exampleProduct := &models.Product{}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		tx, err := testUtil.PlainDB.Begin()
		assert.Nil(t, err)

		actual, err := createProductsInDBFromOptionRows(testUtil.MockDB, tx, exampleProductRoot, exampleProduct)
		assert.Nil(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("with error creating products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), generateExampleTimeForTests(), generateArbitraryError())

		tx, err := testUtil.PlainDB.Begin()
		assert.Nil(t, err)

		_, err = createProductsInDBFromOptionRows(testUtil.MockDB, tx, exampleProductRoot, exampleProduct)
		assert.NotNil(t, err)
	})

}

func TestProductExistenceHandler(t *testing.T) {
	exampleSKU := "example"
	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
	})
	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
	t.Run("with error performing check", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestProductRetrievalHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
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
		AvailableOn:   generateExampleTimeForTests(),
	}
	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with DB error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestProductListHandler(t *testing.T) {
	exampleProduct := models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
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
		AvailableOn:   generateExampleTimeForTests(),
	}
	exampleLength := uint64(3)

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).
			Return(exampleLength, nil)
		testUtil.MockDB.On("GetProductList", mock.Anything, mock.Anything).
			Return([]models.Product{exampleProduct, exampleProduct, exampleProduct}, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error retrieving count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).Return(exampleLength, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with database error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductCount", mock.Anything, mock.Anything).Return(exampleLength, nil)
		testUtil.MockDB.On("GetProductList", mock.Anything, mock.Anything).
			Return([]models.Product{exampleProduct, exampleProduct, exampleProduct}, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductUpdateHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
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
		AvailableOn:   generateExampleTimeForTests(),
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

	t.Run("normal operation", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
		ensureExpectationsWereMet(t, testUtil.Mock)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with database error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with input validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/product/example", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with SKU validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/product/skateboard", strings.NewReader(badSKUUpdateJSON))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with database error updating product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).Return(generateExampleTimeForTests(), generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductDeletionHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:        2,
		CreatedOn: generateExampleTimeForTests(),
		SKU:       exampleSKU,
		Name:      "Skateboard",
	}

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(generateExampleTimeForTests(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, sql.ErrNoRows).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error beginning transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

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
			Return(generateExampleTimeForTests(), generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

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
			Return(generateExampleTimeForTests(), nil).
			Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(generateExampleTimeForTests(), generateArbitraryError()).
			Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).
			Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).
			Return(generateExampleTimeForTests(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).
			Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductCreationHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
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
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &models.ProductOption{
		ID:            exampleProductOption.ID,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []models.ProductOptionValue{
			{
				ID:              128, // == exampleProductOptionValue.ID,
				CreatedOn:       generateExampleTimeForTests(),
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				ID:              256, // == exampleProductOptionValue.ID,
				CreatedOn:       generateExampleTimeForTests(),
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				ID:              512, // == exampleProductOptionValue.ID,
				CreatedOn:       generateExampleTimeForTests(),
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

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("without options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.Mock.ExpectCommit()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with error validating input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with invalid product sku", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(badSKUUpdateJSON))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with already existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("when transaction fails to begin", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("where product creation fails", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), generateArbitraryError())
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating option products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), generateArbitraryError())
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating bridge entries", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("when commit returns an error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootWithSKUPrefixExists", mock.Anything, exampleProduct.SKU).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductRoot", mock.Anything, mock.Anything).
			Return(exampleRoot.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProduct", mock.Anything, mock.Anything).
			Return(exampleProduct.ID, generateExampleTimeForTests(), generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateMultipleProductVariantBridgesForProductID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(expectedCreatedProductOption.Values[0].ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
