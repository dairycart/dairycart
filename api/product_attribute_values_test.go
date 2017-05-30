package api

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

func setExpectationsForProductAttributeValueRetrieval(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).AddRow(productAttributeValueData...)
	query := formatQueryForSQLMock(buildProductAttributeValueRetrievalQuery(v.ID))
	mock.ExpectQuery(query).WithArgs(v.ID).WillReturnRows(exampleRows).WillReturnError(err)
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

func TestRetrieveProductAttributeValueFromDB(t *testing.T) {
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
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, sql.ErrNoRows)

	_, err = retrieveProductAttributeValueFromDB(db, exampleProductAttributeValue.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)

	actual, err := createProductAttributeValueInDB(db, exampleProductAttributeValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttributeValue.ID, actual, "AttributeValue should be returned after successful creation")
	ensureExpectationsWereMet(t, mock)
}

func TestUpdateProductAttributeValueInDB(t *testing.T) {
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
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)

	attributeValueEndpoint := buildRoute("product_attribute_values", "123")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithNonexistentProductAttribute(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attribute_values", "123")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)

	attributeValueEndpoint := buildRoute("product_attribute_values", "123")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 400, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithRowCreationError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attribute_values", "123")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)
	setExpectationsForProductAttributeValueUpdate(mock, exampleUpdatedProductAttributeValue, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereAttributeDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, false, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereAttributeValueDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, false, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereInputIsInvalid(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 400, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereErrorEncounteredRetrievingAttribute(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, arbitraryError)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueUpdateHandlerWhereErrorEncounteredUpdatingAttribute(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttributeValue.ProductAttributeID))
	attributeValueIDString := strconv.Itoa(int(exampleProductAttributeValue.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeValueExistence(mock, exampleProductAttributeValue, true, nil)
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)
	setExpectationsForProductAttributeValueUpdate(mock, exampleUpdatedProductAttributeValue, arbitraryError)

	productAttributeValueEndpoint := buildRoute("product_attribute_values", attributeIDString, attributeValueIDString)
	req, err := http.NewRequest("PUT", productAttributeValueEndpoint, strings.NewReader(exampleProductAttributeValueUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
