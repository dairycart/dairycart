package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProductRootFromProduct(t *testing.T) {
	exampleInput := &models.Product{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		Manufacturer:       "mfgr",
		Brand:              "brand",
		QuantityPerPackage: 666,
		Taxable:            true,
		Cost:               12.34,
		ProductWeight:      1,
		ProductHeight:      1,
		ProductWidth:       1,
		ProductLength:      1,
		PackageWeight:      1,
		PackageHeight:      1,
		PackageWidth:       1,
		PackageLength:      1,
		AvailableOn:        buildTestTime(),
	}
	expected := &models.ProductRoot{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKUPrefix:          "sku",
		Manufacturer:       "mfgr",
		Brand:              "brand",
		QuantityPerPackage: 666,
		Taxable:            true,
		Cost:               12.34,
		ProductWeight:      1,
		ProductHeight:      1,
		ProductWidth:       1,
		ProductLength:      1,
		PackageWeight:      1,
		PackageHeight:      1,
		PackageWidth:       1,
		PackageLength:      1,
		AvailableOn:        buildTestTime(),
	}
	actual := createProductRootFromProduct(exampleInput)

	assert.Equal(t, expected, actual, "expected output should match actual output")
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestSingleProductRootRetrievalHandler(t *testing.T) {
	exampleProductOption := models.ProductOption{
		ID:            123,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    buildTestTime(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := models.Product{
		ID:          2,
		CreatedOn:   buildTestTime(),
		Name:        "Skateboard",
		Description: "This is a skateboard. Please wear a helmet.",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, nil)
		testUtil.MockDB.On("GetProductOptionsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.ProductOption{exampleProductOption}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error querying database for product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving associated products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving product options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, nil)
		testUtil.MockDB.On("GetProductOptionsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.ProductOption{exampleProductOption}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, nil)
		testUtil.MockDB.On("GetProductOptionsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.ProductOption{exampleProductOption}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})
}

func TestProductRootListRetrievalHandler(t *testing.T) {
	exampleProductRoot := models.ProductRoot{
		ID:           2,
		CreatedOn:    buildTestTime(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := models.Product{
		ID:          2,
		CreatedOn:   buildTestTime(),
		SKU:         "skateboard",
		Name:        "Skateboard",
		UPC:         "1234567890",
		Quantity:    123,
		Price:       12.34,
		Cost:        5,
		Taxable:     true,
		Description: "This is a skateboard. Please wear a helmet.",
	}

	t.Run("optimal behavior", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRootCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetProductRootList", mock.Anything, mock.Anything).
			Return([]models.ProductRoot{exampleProductRoot}, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error getting row count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRootCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving product root list", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRootCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetProductRootList", mock.Anything, mock.Anything).
			Return([]models.ProductRoot{exampleProductRoot}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRootCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetProductRootList", mock.Anything, mock.Anything).
			Return([]models.ProductRoot{exampleProductRoot}, nil)
		testUtil.MockDB.On("GetProductsByProductRootID", mock.Anything, exampleProductRoot.ID).
			Return([]models.Product{exampleProduct}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductRootDeletionHandler(t *testing.T) {
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    buildTestTime(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	t.Run("optimal behavior", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error starting transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error archiving bridge entries", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error archiving product option values", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error archiving options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error archiving products", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error archiving product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error committing transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(exampleProductRoot, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductVariantBridgesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionValuesWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductOptionsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("ArchiveProductsWithProductRootID", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductRoot", mock.Anything, exampleProductRoot.ID).
			Return(buildTestTime(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
