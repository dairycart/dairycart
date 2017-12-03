package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionValueCreationHandler(t *testing.T) {
	exampleProductOptionValueCreationBody := `{"value": "something"}`
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOptionValue.ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(false, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error checking product option existence", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with pre-existing value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error checking for value existence", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOptionValue.ID, generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error committing transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOptionValue.ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

}

func TestProductOptionValueUpdateHandler(t *testing.T) {
	exampleProductOptionValueUpdateBody := `{"value": "something else"}`
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("UpdateProductOptionValue", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("UpdateProductOptionValue", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductOptionValueDeletionHandler(t *testing.T) {
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("DeleteProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("DeleteProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
