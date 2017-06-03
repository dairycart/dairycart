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
	exampleProductAttributeValueCreationBody = `{"value": "something"}`
	exampleProductAttributeValueUpdateBody   = `{"value": "something else"}`
)

var exampleProductAttributeValue *ProductAttributeValue
var exampleUpdatedProductAttributeValue *ProductAttributeValue
var productAttributeValueHeaders []string
var productAttributeValueData []driver.Value

func init() {
	exampleProductAttributeValue = &ProductAttributeValue{
		ID:                 256,
		ProductAttributeID: 123, // == exampleProductAttribute.ID
		Value:              "something",
		CreatedAt:          exampleTime,
	}
	exampleUpdatedProductAttributeValue = &ProductAttributeValue{
		ID:                 256,
		ProductAttributeID: 123, // == exampleProductAttribute.ID
		Value:              "something else",
		CreatedAt:          exampleTime,
	}
	productAttributeValueHeaders = []string{"id", "product_attribute_id", "value", "created_at", "updated_at", "archived_at"}
	productAttributeValueData = []driver.Value{
		exampleProductAttributeValue.ID,
		exampleProductAttributeValue.ProductAttributeID,
		exampleProductAttributeValue.Value,
		exampleTime,
		nil,
		nil,
	}
}

func setExpectationsForProductAttributeValueExistence(mock sqlmock.Sqlmock, v *ProductAttributeValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductAttributeValueExistenceQuery(v.ID)
	stringID := strconv.Itoa(int(v.ID))
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeValueRetrievalByAttributeID(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).AddRow(productAttributeValueData...)
	query := formatQueryForSQLMock(buildProductAttributeValueRetrievalForAttributeIDQuery(a.ID))
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeValueRetrieval(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).AddRow(productAttributeValueData...)
	query := formatQueryForSQLMock(buildProductAttributeValueRetrievalQuery(v.ID))
	mock.ExpectQuery(query).
		WithArgs(v.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeValueCreation(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProductAttributeValue.ID)
	query, _ := buildProductAttributeValueCreationQuery(v)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(v.ProductAttributeID, v.Value).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeValueUpdate(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeHeaders).
		AddRow([]driver.Value{v.ID, v.ProductAttributeID, v.Value, exampleTime, nil, nil}...)
	query, args := buildProductAttributeValueUpdateQuery(v)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeValueForAttributeExistence(mock sqlmock.Sqlmock, a *ProductAttribute, v *ProductAttributeValue, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query, args := buildProductAttributeValueExistenceForAttributeIDQuery(a.ID, v.Value)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestValidateProductAttributeValueCreationInput(t *testing.T) {
	t.Parallel()
	expected := &ProductAttributeValue{Value: "something"}
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleProductAttributeValueCreationBody))
	actual, err := validateProductAttributeValueCreationInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "ProductAttributeUpdateInput should match expectation")
}

func TestValidateProductAttributeValueCreationInputWithCompletelyInvalidInput(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateProductAttributeValueCreationInput(req)
	assert.NotNil(t, err)
}

func TestValidateProductAttributeValueCreationInputWithGarbageInput(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleGarbageInput))
	_, err := validateProductAttributeValueCreationInput(req)
	assert.NotNil(t, err)
}

func TestRetrieveProductAttributeValueFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)

	actual, err := retrieveProductAttributeValueFromDB(db, exampleProductAttributeValue.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttributeValue, actual, "expected and actual product attribute values should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductAttributeValueFromDBThatDoesNotExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, sql.ErrNoRows)

	_, err = retrieveProductAttributeValueFromDB(db, exampleProductAttributeValue.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeValue(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectBegin()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductAttributeValueInDB(tx, exampleProductAttributeValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttributeValue.ID, actual, "AttributeValue should be returned after successful creation")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestUpdateProductAttributeValueInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueUpdate(mock, exampleProductAttributeValue, nil)

	err = updateProductAttributeValueInDB(db, exampleProductAttributeValue)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

//////////////////////////////////////////////////////////////////////
//                                                                  //
//    HTTP                                                          //
//     Handler                                     _.-~`  `~-.      //
//        Tests        _.--~~~---,.__          _.,;; .   -=(@'`\    //
//                  .-`              ``~~~~--~~` ';;;       ____)   //
//               _.'            '.              ';;;;;    '`_.'     //
//            .-~;`               `\           ' ';;;;;__.~`        //
//          .' .'          `'.     |           /  /;''              //
//           \/      .---'''``)   /'-._____.--'\  \                 //
//          _/|    (`        /  /`              `\ \__              //
//   ',    `/- \   \      __/  (_                /-\-\-`            //
//     `;'-..___)   |     `/-\-\-`                                  //
//       `-.       .'                                               //
//          `~~~~``                                                 //
//////////////////////////////////////////////////////////////////////

func TestProductAttributeValueCreationHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)
	mock.ExpectCommit()

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWhenTransactionFailsToBegin(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, nil)
	mock.ExpectBegin().WillReturnError(arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWhenTransactionFailsToCommit(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)
	mock.ExpectCommit().WillReturnError(arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithNonexistentProductAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWhenValueAlreadyExistsForAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, true, nil)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWhenValueExistenceCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, sql.ErrNoRows)
	mock.ExpectBegin()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)
	mock.ExpectCommit()

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWhenValueExistenceCheckReturnsError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithRowCreationError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueForAttributeExistence(mock, exampleProductAttribute, exampleProductAttributeValue, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, arbitraryError)
	mock.ExpectRollback()

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)
	setExpectationsForProductAttributeValueUpdate(mock, exampleUpdatedProductAttributeValue, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereAttributeValueDoesNotExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, false, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereInputIsInvalid(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereErrorEncounteredRetrievingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, arbitraryError)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereErrorEncounteredUpdatingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)
	setExpectationsForProductAttributeValueUpdate(mock, exampleUpdatedProductAttributeValue, arbitraryError)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
