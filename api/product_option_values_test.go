package main

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductOptionValueCreationBody = `{"value": "something"}`
	exampleProductOptionValueUpdateBody   = `{"value": "something else"}`
)

var (
	exampleProductOptionValue        *ProductOptionValue
	exampleUpdatedProductOptionValue *ProductOptionValue
	productOptionValueHeaders        []string
	productOptionValueData           []driver.Value
)

func init() {
	exampleProductOptionValue = &ProductOptionValue{
		ID:              256,
		ProductOptionID: 123, // == exampleProductOption.ID
		Value:           "something",
		CreatedOn:       generateExampleTimeForTests(),
	}
	exampleUpdatedProductOptionValue = &ProductOptionValue{
		ID:              256,
		ProductOptionID: 123, // == exampleProductOption.ID
		Value:           "something else",
		CreatedOn:       generateExampleTimeForTests(),
	}
	productOptionValueHeaders = []string{"id", "product_option_id", "value", "created_on", "updated_on", "archived_on"}
	productOptionValueData = []driver.Value{
		exampleProductOptionValue.ID,
		exampleProductOptionValue.ProductOptionID,
		exampleProductOptionValue.Value,
		generateExampleTimeForTests(),
		nil,
		nil,
	}
}

func setExpectationsForProductOptionValueExistence(mock sqlmock.Sqlmock, v *ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionValueExistenceQuery)
	stringID := strconv.Itoa(int(v.ID))
	mock.ExpectQuery(query).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueRetrievalByOptionID(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionValueHeaders).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(productOptionValueRetrievalForOptionIDQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueRetrieval(mock sqlmock.Sqlmock, v *ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows(productOptionValueHeaders).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(productOptionValueRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(v.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueCreation(mock sqlmock.Sqlmock, v *ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProductOptionValue.ID)
	query, _ := buildProductOptionValueCreationQuery(v)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(v.ProductOptionID, v.Value).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueUpdate(mock sqlmock.Sqlmock, v *ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows(productOptionHeaders).
		AddRow([]driver.Value{v.ID, v.ProductOptionID, v.Value, generateExampleTimeForTests(), nil, nil}...)
	query, args := buildProductOptionValueUpdateQuery(v)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueForOptionExistence(mock sqlmock.Sqlmock, a *ProductOption, v *ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionValueExistenceForOptionIDQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID, v.Value).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestRetrieveProductOptionValueFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)

	actual, err := retrieveProductOptionValueFromDB(testUtil.DB, exampleProductOptionValue.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue, actual, "expected and actual product option values should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductOptionValueFromDBThatDoesNotExist(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, sql.ErrNoRows)

	_, err := retrieveProductOptionValueFromDB(testUtil.DB, exampleProductOptionValue.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionValueInDB(tx, exampleProductOptionValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue.ID, actual, "OptionValue should be returned after successful creation")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateProductOptionValueInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleProductOptionValue, nil)

	err := updateProductOptionValueInDB(testUtil.DB, exampleProductOptionValue)
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

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToBegin(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(arbitraryError)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToCommit(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(arbitraryError)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithNonexistentProductOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, arbitraryError)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueAlreadyExistsForOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, true, nil)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, sql.ErrNoRows)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, nil)
	testUtil.Mock.ExpectCommit()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, arbitraryError)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueCreationHandlerWithRowCreationError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(testUtil.Mock, exampleProductOption, exampleProductOptionValue, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(testUtil.Mock, exampleProductOptionValue, arbitraryError)
	testUtil.Mock.ExpectRollback()

	optionValueEndpoint := buildRoute("v1", "product_options", "123", "value")
	req, err := http.NewRequest(http.MethodPost, optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleUpdatedProductOptionValue, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereOptionValueDoesNotExist(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, false, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereInputIsInvalid(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredRetrievingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, arbitraryError)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredUpdatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(testUtil.Mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(testUtil.Mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(testUtil.Mock, exampleUpdatedProductOptionValue, arbitraryError)

	productOptionValueEndpoint := buildRoute("v1", "product_option_values", optionValueIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
