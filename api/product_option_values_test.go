// +build migrated

package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setExpectationsForProductOptionValueExistence(mock sqlmock.Sqlmock, v *models.ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionValueExistenceQuery)
	stringID := strconv.Itoa(int(v.ID))
	mock.ExpectQuery(query).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueRetrieval(mock sqlmock.Sqlmock, v *models.ProductOptionValue, err error) {
	productOptionValueData := []driver.Value{
		256,
		123,
		"something",
		generateExampleTimeForTests(),
		nil,
		nil,
	}
	exampleRows := sqlmock.NewRows([]string{"id", "product_option_id", "value", "created_on", "updated_on", "archived_on"}).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(productOptionValueRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(v.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueUpdate(mock sqlmock.Sqlmock, v *models.ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	query, args := buildProductOptionValueUpdateQuery(v)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueForOptionExistence(mock sqlmock.Sqlmock, a *models.ProductOption, v *models.ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionValueExistenceForOptionIDQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID, v.Value).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueDeletion(mock sqlmock.Sqlmock, id uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productOptionValueDeletionQuery)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestRetrieveProductOptionValueFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)

	actual, err := retrieveProductOptionValueFromDB(testUtil.DB, exampleProductOptionValue.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue, actual, "expected and actual product option values should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductOptionValueFromDBThatDoesNotExist(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, sql.ErrNoRows)

	_, err := retrieveProductOptionValueFromDB(testUtil.DB, exampleProductOptionValue.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	newID, createdOn, err := createProductOptionValueInDB(tx, exampleProductOptionValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue.ID, newID, "OptionValue ID should be returned after successful creation")
	assert.Equal(t, generateExampleTimeForTests(), createdOn, "OptionValue CreatedOn should be returned after successful creation")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateProductOptionValueInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleProductOptionValue, nil)

	updatedOn, err := updateProductOptionValueInDB(testUtil.DB, exampleProductOptionValue)
	assert.Nil(t, err)
	assert.Equal(t, generateExampleTimeForTests(), updatedOn, "updateProductOptionValueInDB should return the correct updated time")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveProductOptionValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionValueDeletion(testUtil.Mock, 1, nil)

	err := archiveProductOptionValue(testUtil.DB, 1)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(false, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error checking product option existence", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with pre-existing value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error checking for value existence", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductOptionExists", mock.Anything, exampleProductOption.ID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionValueForOptionIDExists", mock.Anything, exampleProductOption.ID, exampleProductOptionValue.Value).
			Return(false, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product_options/%d/value", exampleProductOption.ID), strings.NewReader(exampleProductOptionValueCreationBody))
		assert.Nil(t, err)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("UpdateProductOptionValue", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), strings.NewReader(exampleProductOptionValueUpdateBody))
		assert.Nil(t, err)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting option value", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(exampleProductOptionValue, nil)
		testUtil.MockDB.On("DeleteProductOptionValue", mock.Anything, exampleProductOptionValue.ID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_option_values/%d", exampleProductOptionValue.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
