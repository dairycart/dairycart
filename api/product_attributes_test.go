package api

import (
	"database/sql/driver"
	"strconv"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var exampleProductAttribute *ProductAttribute
var productAttributeHeaders []string
var productAttributeData []driver.Value

func init() {
	exampleProductAttribute = &ProductAttribute{
		ID:   123,
		Name: "attribute",
	}
	productAttributeHeaders = []string{"id", "name", "product_progenitor_id", "created_at", "updated_at", "archived_at"}
	productAttributeData = []driver.Value{1, 2, "Attribute", exampleTime, nil, nil}
}

func setExpectationsForProductAttributeExistence(mock sqlmock.Sqlmock, id int64, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductAttributeexistenceQuery(id)
	stringID := strconv.Itoa(int(id))
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}
