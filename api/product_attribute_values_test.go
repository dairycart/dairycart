package api

import (
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
		ProductAttributeID: 123, // exampleProductAttribute.ID
		Value:              "something",
	}
	productAttributeValueHeaders = []string{"id", "product_attribute_id", "value", "created_at", "updated_at", "archived_at"}
	productAttributeValueData = []driver.Value{1, 2, "Value", exampleTime, nil, nil}
}

func setExpectationsForProductAttributeValueCreation(mock sqlmock.Sqlmock, v *ProductAttributeValue, err error) {
	exampleRows := sqlmock.NewRows(productAttributeValueHeaders).
		AddRow(productAttributeValueData...)

	query, _ := buildProductAttributeValueCreationQuery(v)
	mock.ExpectQuery(formatConstantQueryForSQLMock(query)).
		WithArgs(
			v.ProductAttributeID,
			v.Value,
		).
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
