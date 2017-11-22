package main

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductOptionValueCreationBody = `{"value": "something"}`
	exampleProductOptionValueUpdateBody   = `{"value": "something else"}`
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

func setExpectationsForProductOptionValueRetrievalByOptionID(mock sqlmock.Sqlmock, a *models.ProductOption, err error) {
	productOptionValueData := []driver.Value{
		256,
		123,
		"something",
		generateExampleTimeForTests(),
		nil,
		nil,
	}
	exampleRows := sqlmock.NewRows([]string{"id", "product_option_id", "value", "created_on", "updated_on", "archived_on"}).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(productOptionValueRetrievalForOptionIDQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID).
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

func setExpectationsForMultipleProductOptionValuesCreation(mock sqlmock.Sqlmock, vs []models.ProductOptionValue, err error, errorIndex int) {
	for i, v := range vs {
		var errToUse error = nil
		if i == errorIndex && err != nil {
			errToUse = err
		}
		setExpectationsForProductOptionValueCreation(mock, &v, errToUse)
	}
}

func setExpectationsForProductValueBridgeEntryCreation(mock sqlmock.Sqlmock, productID uint64, optionValueIDs []uint64, err error) {
	query, _ := buildProductVariantBridgeCreationQuery(productID, optionValueIDs)
	mock.ExpectExec(formatQueryForSQLMock(query)).
		// I can't think of a sane way to expect a given set of arguments, so we'll just have to count the queries for now
		//WithArgs(queryArgs...).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueCreation(mock sqlmock.Sqlmock, v *models.ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(v.ID, generateExampleTimeForTests())
	query, _ := buildProductOptionValueCreationQuery(v)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(v.ProductOptionID, v.Value).
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
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToBegin(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToCommit(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithNonexistentProductOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, generateArbitraryError())

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueAlreadyExistsForOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, true, nil)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, sql.ErrNoRows)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, generateArbitraryError())

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithRowCreationError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}
	exampleUpdatedProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something else",
	}

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleUpdatedProductOptionValue, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereOptionValueDoesNotExist(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, false, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereInputIsInvalid(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredRetrievingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, generateArbitraryError())

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredUpdatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something",
	}
	exampleUpdatedProductOptionValue := &models.ProductOptionValue{
		ID:              256,
		CreatedOn:       generateExampleTimeForTests(),
		ProductOptionID: 123,
		Value:           "something else",
	}

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleUpdatedProductOptionValue, generateArbitraryError())

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "product_option_values", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForProductOptionValueExistence(testUtil.Mock, &models.ProductOptionValue{ID: exampleID}, true, nil)
	setExpectationsForProductOptionValueDeletion(testUtil.Mock, exampleID, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueDeletionHandlerWithNonexistentOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "product_option_values", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForProductOptionValueExistence(testUtil.Mock, &models.ProductOptionValue{ID: exampleID}, false, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueDeletionHandlerWithErrorDeletingValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "product_option_values", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForProductOptionValueExistence(testUtil.Mock, &models.ProductOptionValue{ID: exampleID}, true, nil)
	setExpectationsForProductOptionValueDeletion(testUtil.Mock, exampleID, generateArbitraryError())
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}
