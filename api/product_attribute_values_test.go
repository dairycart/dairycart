package api

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductAttributeValueCreationBody = `{"value": "something"}`
)

var exampleProductAttributeValue *ProductAttributeValue
var productAttributeValueHeaders []string
var productAttributeValueData []driver.Value

func init() {
	exampleProductAttributeValue = &ProductAttributeValue{
		ID:                 256,
		ProductAttributeID: 123, // == exampleProductAttribute.ID
		Value:              "something",
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

func setExpectationsForProductAttributeValueRetrieval(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).AddRow(productAttributeValueData...)
	query := formatQueryForSQLMock(buildProductAttributeValueRetrievalQuery(v.ID))
	mock.ExpectQuery(query).WithArgs(v.ID).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestRetrieveProductAttributeValueFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueRetrieval(mock, exampleProductAttributeValue, nil)

	_, err = retrieveProductAttributeValueFromDB(db, exampleProductAttributeValue.ID)
	assert.Nil(t, err)
	// TODO: fix this part of the test.
	// assert.Equal(t, exampleProductAttributeValue, actual, "expected and actual products should match")
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

func setExpectationsForProductAttributeValueCreation(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).AddRow(productAttributeValueData...)
	query, _ := buildProductAttributeValueCreationQuery(v)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(v.ProductAttributeID, v.Value).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductAttributeValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)

	actual, err := createProductAttributeValueInDB(db, exampleProductAttributeValue)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttributeValue, actual, "AttributeValue should be returned after successful creation")
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
	setExpectationsForProductAttributeExistence(mock, exampleProductAttribute.ID, true, nil)
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, nil)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
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
	setExpectationsForProductAttributeExistence(mock, exampleProductAttribute.ID, true, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeValueCreationHandlerWithInvalidValueBody(t *testing.T) {
	exampleInvalidProductAttributeValueCreationBody := `{"things": "stuff"}`
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeExistence(mock, exampleProductAttribute.ID, true, nil)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleInvalidProductAttributeValueCreationBody))
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
	setExpectationsForProductAttributeExistence(mock, exampleProductAttribute.ID, true, nil)
	setExpectationsForProductAttributeValueCreation(mock, exampleProductAttributeValue, arbitraryError)

	attributeValueEndpoint := buildRoute("product_attributes", "123", "value")
	req, err := http.NewRequest("POST", attributeValueEndpoint, strings.NewReader(exampleProductAttributeValueCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
