package main

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

var (
	exampleProductOption         *ProductOption
	exampleUpdatedProductOption  *ProductOption
	expectedCreatedProductOption *ProductOption
	exampleProductOptionInput    *ProductOptionCreationInput
	productOptionHeaders         []string
)

func init() {
	exampleProductOption = &ProductOption{
		ID:        123,
		Name:      "something",
		ProductID: 2, // == exampleProduct.ID
		CreatedOn: generateExampleTimeForTests(),
	}
	exampleUpdatedProductOption = &ProductOption{
		ID:        exampleProductOption.ID,
		Name:      "something else",
		ProductID: exampleProductOption.ProductID,
	}
	productOptionHeaders = []string{"id", "name", "product_id", "created_on", "updated_on", "archived_on"}

	expectedCreatedProductOption = &ProductOption{
		ID:        exampleProductOption.ID,
		Name:      "something",
		ProductID: exampleProductOption.ProductID,
		Values: []ProductOptionValue{
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

func setExpectationsForProductOptionExistenceByName(mock sqlmock.Sqlmock, a *ProductOption, productID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionExistenceQueryForProductByName)
	mock.ExpectQuery(query).
		WithArgs(a.Name, productID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionRetrievalQuery(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductID, generateExampleTimeForTests(), nil, nil}...)
	query := formatQueryForSQLMock(productOptionRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionListQueryWithCount(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"count", "id", "name", "product_id", "created_on", "updated_on", "archived_on"}).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductID, generateExampleTimeForTests(), nil, nil}...)
	query, _ := buildProductOptionListQuery(exampleProduct.ID, defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionCreation(mock sqlmock.Sqlmock, a *ProductOption, productID uint64, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProductOption.ID)
	query, args := buildProductOptionCreationQuery(a, productID)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionUpdate(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows(productOptionHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductID, generateExampleTimeForTests(), nil, nil}...)
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
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)

	actual, err := retrieveProductOptionFromDB(testUtil.DB, exampleProductOption.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption, actual, "expected and actual product options should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductOptionFromDBWithNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, sql.ErrNoRows)

	_, err := retrieveProductOptionFromDB(testUtil.DB, exampleProductOption.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, exampleProductOption, exampleProduct.ID, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionInDB(tx, exampleProductOption, exampleProduct.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption, actual, "Creating a product option should return the created product option")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, exampleProduct.ID)
	assert.Nil(t, err)
	assert.Equal(t, expectedCreatedProductOption, actual, "output from product option creation should match expectation")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, arbitraryError)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOptionValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], arbitraryError)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateProductOptionInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionUpdate(testUtil.Mock, expectedCreatedProductOption, nil)

	err := updateProductOptionInDB(testUtil.DB, exampleProductOption)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionListHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionListQueryWithCount(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProduct.ID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")

	expected := &ProductOptionsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductOptionsResponse{}
	bodyString := testUtil.Response.Body.String()
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product option pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product option limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product option counts should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionListHandlerWithErrorsRetrievingValues(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionListQueryWithCount(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProduct.ID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionListHandlerWithDBErrors(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductOptionListQueryWithCount(testUtil.Mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProduct.ID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerFailureToSetupTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerFailureToCommitTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit().WillReturnError(arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, sql.ErrNoRows)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithNonExistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, false, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithInvalidOptionCreationInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameExists(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, true, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithProblemsCreatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	productIDString := strconv.Itoa(int(exampleProductOption.ProductID))
	setExpectationsForProductExistenceByID(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProduct.ID, arbitraryError)
	testUtil.Mock.ExpectRollback()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(testUtil.Mock, exampleUpdatedProductOption, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPut, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithNonexistentOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, false, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPut, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPut, productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithErrorRetrievingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPut, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithErrorUpdatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(testUtil.Mock, exampleUpdatedProductOption, arbitraryError)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPut, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
