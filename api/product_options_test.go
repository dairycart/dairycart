package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductOptionCreationBody = `
		{
			"name": "something",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`
	exampleProductOptionUpdateBody = `
		{
			"name": "something else"
		}
	`
)

var exampleProductOption *ProductOption
var exampleUpdatedProductOption *ProductOption
var expectedCreatedProductOption *ProductOption
var exampleProductOptionInput *ProductOptionCreationInput
var productOptionHeaders []string

func init() {
	exampleProductOption = &ProductOption{
		ID:                  123,
		Name:                "something",
		ProductProgenitorID: 2, // == exampleProgenitor.ID
		CreatedOn:           exampleTime,
	}
	exampleUpdatedProductOption = &ProductOption{
		ID:                  exampleProductOption.ID,
		Name:                "something else",
		ProductProgenitorID: exampleProductOption.ProductProgenitorID,
	}
	productOptionHeaders = []string{"id", "name", "product_progenitor_id", "created_on", "updated_on", "archived_on"}

	expectedCreatedProductOption = &ProductOption{
		ID:                  exampleProductOption.ID,
		Name:                "something",
		ProductProgenitorID: exampleProductOption.ProductProgenitorID,
		Values: []*ProductOptionValue{
			{
				ID:              256, // == exampleProductOptionValue.ID,
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				ID:              256, // == exampleProductOptionValue.ID,
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				ID:              256, // == exampleProductOptionValue.ID,
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

	exampleProductOptionInput = &ProductOptionCreationInput{
		Name:   "something",
		Values: []string{"one", "two", "three"},
	}
}

func setExpectationsForProductOptionExistenceByID(mock sqlmock.Sqlmock, a *ProductOption, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionExistenceQuery)
	stringID := strconv.Itoa(int(a.ID))
	mock.ExpectQuery(query).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionExistenceByName(mock sqlmock.Sqlmock, a *ProductOption, progenitorID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionExistenceQueryForProductByName)
	mock.ExpectQuery(query).
		WithArgs(a.Name, progenitorID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionRetrievalQuery(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query := formatQueryForSQLMock(productOptionRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionListQueryWithCount(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"count", "id", "name", "product_progenitor_id", "created_on", "updated_on", "archived_on"}).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query, _ := buildProductOptionListQuery(strconv.Itoa(int(exampleProgenitor.ID)), defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionCreation(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProductOption.ID)
	query, args := buildProductOptionCreationQuery(a)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionUpdate(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query, args := buildProductOptionUpdateQuery(a)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestValidateProductOptionUpdateInput(t *testing.T) {
	t.Parallel()
	expected := &ProductOptionUpdateInput{Name: "something else"}
	exampleInput := strings.NewReader(exampleProductOptionUpdateBody)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual, err := validateProductOptionUpdateInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "ProductOptionUpdateInput should match expectation")
}

func TestValidateProductOptionUpdateInputWithInvalidInput(t *testing.T) {
	t.Parallel()
	exampleInput := strings.NewReader(exampleGarbageInput)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateProductOptionUpdateInput(req)

	assert.NotNil(t, err)
}

func TestRetrieveProductOptionFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionRetrievalQuery(mock, exampleProductOption, nil)

	actual, err := retrieveProductOptionFromDB(db, exampleProductOption.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption, actual, "expected and actual product options should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductOptionFromDBWithNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionRetrievalQuery(mock, exampleProductOption, sql.ErrNoRows)

	_, err = retrieveProductOptionFromDB(db, exampleProductOption.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductOptionInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, exampleProductOption, nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionInDB(tx, exampleProductOption)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption, actual, "Creating a product option should return the created product option")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductOptionAndValuesInDBFromInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[2], nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductProgenitorID)
	assert.Nil(t, err)
	assert.Equal(t, expectedCreatedProductOption, actual, "output from product option creation should match expectation")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, arbitraryError)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductProgenitorID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOptionValue(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[0], arbitraryError)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductProgenitorID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestUpdateProductOptionInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductOptionUpdate(mock, expectedCreatedProductOption, nil)

	err = updateProductOptionInDB(db, exampleProductOption)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionListHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionListQueryWithCount(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, nil)

	productOptionEndpoint := buildRoute("product_options", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productOptionEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")

	expected := &ProductOptionsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductOptionsResponse{}
	bodyString := res.Body.String()
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product option pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product option limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product option counts should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionListHandlerWithErrorsRetrievingValues(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionListQueryWithCount(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("product_options", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productOptionEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionListHandlerWithDBErrors(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductOptionListQueryWithCount(mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("product_options", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productOptionEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[2], nil)
	mock.ExpectCommit()

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerFailureToSetupTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, false, nil)
	mock.ExpectBegin().WillReturnError(arbitraryError)

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerFailureToCommitTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[2], nil)
	mock.ExpectCommit().WillReturnError(arbitraryError)

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, false, sql.ErrNoRows)
	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(mock, expectedCreatedProductOption.Values[2], nil)
	mock.ExpectCommit()

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerWithNonExistentProgenitor(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, false)

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerWithInvalidOptionCreationInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameExists(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, true, nil)

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionCreationHandlerWithProblemsCreatingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductOption.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductOptionExistenceByName(mock, expectedCreatedProductOption, progenitorIDString, false, nil)

	mock.ExpectBegin()
	setExpectationsForProductOptionCreation(mock, expectedCreatedProductOption, arbitraryError)
	mock.ExpectRollback()

	productOptionEndpoint := buildRoute("product_options", progenitorIDString)
	req, err := http.NewRequest("POST", productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionUpdateHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(mock, exampleUpdatedProductOption, nil)

	productOptionEndpoint := buildRoute("product_options", optionIDString)
	req, err := http.NewRequest("PUT", productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionUpdateHandlerWithNonexistentOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, false, nil)

	productOptionEndpoint := buildRoute("product_options", optionIDString)
	req, err := http.NewRequest("PUT", productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionUpdateHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)

	productOptionEndpoint := buildRoute("product_options", optionIDString)
	req, err := http.NewRequest("PUT", productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionUpdateHandlerWithErrorRetrievingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("product_options", optionIDString)
	req, err := http.NewRequest("PUT", productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductOptionUpdateHandlerWithErrorUpdatingOption(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(mock, exampleUpdatedProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("product_options", optionIDString)
	req, err := http.NewRequest("PUT", productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
