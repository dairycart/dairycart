package api

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
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

var exampleProductOptionValue *ProductOptionValue
var exampleUpdatedProductOptionValue *ProductOptionValue
var productOptionValueHeaders []string
var productOptionValueData []driver.Value

func init() {
	exampleProductOptionValue = &ProductOptionValue{
		ID:              256,
		ProductOptionID: 123, // == exampleProductOption.ID
		Value:           "something",
		CreatedAt:       exampleTime,
	}
	exampleUpdatedProductOptionValue = &ProductOptionValue{
		ID:              256,
		ProductOptionID: 123, // == exampleProductOption.ID
		Value:           "something else",
		CreatedAt:       exampleTime,
	}
	productOptionValueHeaders = []string{"id", "product_option_id", "value", "created_at", "updated_at", "archived_at"}
	productOptionValueData = []driver.Value{
		exampleProductOptionValue.ID,
		exampleProductOptionValue.ProductOptionID,
		exampleProductOptionValue.Value,
		exampleTime,
		nil,
		nil,
	}
}

func setExpectationsForProductOptionValueExistence(mock sqlmock.Sqlmock, v *ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductOptionValueExistenceQuery(v.ID)
	stringID := strconv.Itoa(int(v.ID))
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueRetrievalByOptionID(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionValueHeaders).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(buildProductOptionValueRetrievalForOptionIDQuery(a.ID))
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueRetrieval(mock sqlmock.Sqlmock, v *ProductOptionValue, err error) {
	exampleRows := sqlmock.NewRows(productOptionValueHeaders).AddRow(productOptionValueData...)
	query := formatQueryForSQLMock(buildProductOptionValueRetrievalQuery(v.ID))
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
		AddRow([]driver.Value{v.ID, v.ProductOptionID, v.Value, exampleTime, nil, nil}...)
	query, args := buildProductOptionValueUpdateQuery(v)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionValueForOptionExistence(mock sqlmock.Sqlmock, a *ProductOption, v *ProductOptionValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query, args := buildProductOptionValueExistenceForOptionIDQuery(a.ID, v.Value)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestValidateProductOptionValueCreationInput(t *testing.T) {
	t.Parallel()
	expected := &ProductOptionValue{Value: "something"}
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleProductOptionValueCreationBody))
	actual, err := validateProductOptionValueCreationInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "ProductOptionUpdateInput should match expectation")
}

func TestValidateProductOptionValueCreationInputWithCompletelyInvalidInput(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateProductOptionValueCreationInput(req)
	assert.NotNil(t, err)
}

func TestValidateProductOptionValueCreationInputWithGarbageInput(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleGarbageInput))
	_, err := validateProductOptionValueCreationInput(req)
	assert.NotNil(t, err)
}

func TestRetrieveProductOptionValueFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionValueRetrieval(mock, exampleProductOptionValue, nil)

	actual, err := retrieveProductOptionValueFromDB(db, exampleProductOptionValue.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue, actual, "expected and actual product option values should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductOptionValueFromDBThatDoesNotExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionValueRetrieval(mock, exampleProductOptionValue, sql.ErrNoRows)

	_, err = retrieveProductOptionValueFromDB(db, exampleProductOptionValue.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductOptionValue(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(mock, exampleProductOptionValue, nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionValueInDB(tx, exampleProductOptionValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOptionValue.ID, actual, "OptionValue should be returned after successful creation")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestUpdateProductOptionValueInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionValueUpdate(mock, exampleProductOptionValue, nil)

	err = updateProductOptionValueInDB(db, exampleProductOptionValue)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionValueCreationHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(mock, exampleProductOptionValue, nil)
	mock.ExpectCommit()

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToBegin(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, nil)
	mock.ExpectBegin().WillReturnError(arbitraryError)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWhenTransactionFailsToCommit(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(mock, exampleProductOptionValue, nil)
	mock.ExpectCommit().WillReturnError(arbitraryError)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWithNonexistentProductOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, arbitraryError)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWhenValueAlreadyExistsForOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, true, nil)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, sql.ErrNoRows)
	mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(mock, exampleProductOptionValue, nil)
	mock.ExpectCommit()

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWhenValueExistenceCheckReturnsError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, arbitraryError)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueCreationHandlerWithRowCreationError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionValueForOptionExistence(mock, exampleProductOption, exampleProductOptionValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductOptionValueCreation(mock, exampleProductOptionValue, arbitraryError)
	mock.ExpectRollback()

	optionValueEndpoint := buildRoute("product_options", "123", "value")
	req, err := http.NewRequest("POST", optionValueEndpoint, strings.NewReader(exampleProductOptionValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueUpdateHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(mock, exampleUpdatedProductOptionValue, nil)

	productOptionValueEndpoint := buildRoute("product_option_values", optionValueIDString)
	req, err := http.NewRequest("PUT", productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueUpdateHandlerWhereOptionValueDoesNotExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(mock, exampleProductOptionValue, false, nil)

	productOptionValueEndpoint := buildRoute("product_option_values", optionValueIDString)
	req, err := http.NewRequest("PUT", productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueUpdateHandlerWhereInputIsInvalid(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(mock, exampleProductOptionValue, true, nil)

	productOptionValueEndpoint := buildRoute("product_option_values", optionValueIDString)
	req, err := http.NewRequest("PUT", productOptionValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredRetrievingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(mock, exampleProductOptionValue, arbitraryError)

	productOptionValueEndpoint := buildRoute("product_option_values", optionValueIDString)
	req, err := http.NewRequest("PUT", productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionValueUpdateHandlerWhereErrorEncounteredUpdatingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionValueIDString := strconv.Itoa(int(exampleProductOptionValue.ID))

	setExpectationsForProductOptionValueExistence(mock, exampleProductOptionValue, true, nil)
	setExpectationsForProductOptionValueRetrieval(mock, exampleProductOptionValue, nil)
	setExpectationsForProductOptionValueUpdate(mock, exampleUpdatedProductOptionValue, arbitraryError)

	productOptionValueEndpoint := buildRoute("product_option_values", optionValueIDString)
	req, err := http.NewRequest("PUT", productOptionValueEndpoint, strings.NewReader(exampleProductOptionValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
