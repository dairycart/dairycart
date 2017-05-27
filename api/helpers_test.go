package api

import (
	"errors"
	"net/http/httptest"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
)

func TestRespondThatRowDoesNotExist(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	respondThatRowDoesNotExist(req, w, "item", "something")

	assert.Equal(t, "The item you were looking for (identified by `something`) does not exist\n", w.Body.String(), "response should indicate the row was not found")
	assert.Equal(t, 404, w.Code, "status code should be 404")
}

func TestNotifyOfInvalidRequestBody(t *testing.T) {
	w := httptest.NewRecorder()
	notifyOfInvalidRequestBody(w, errors.New("test"))

	assert.Equal(t, "test\n", w.Body.String(), "response should indicate the request body was invalid")
	assert.Equal(t, 400, w.Code, "status code should be 404")
}

func TestNotifyOfInternalIssue(t *testing.T) {
	w := httptest.NewRecorder()

	notifyOfInternalIssue(w, errors.New("test"), "do a thing")

	assert.Equal(t, "Unexpected internal error\n", w.Body.String(), "response should indicate their was an internal error")
	assert.Equal(t, 500, w.Code, "status code should be 404")
}

func TestRowExistsInDBForExistingRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	exampleRows := sqlmock.NewRows([]string{""}).AddRow("true")
	skuExistenceQuery := buildProductExistenceQuery(exampleSKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(exampleSKU).
		WillReturnRows(exampleRows)

	exists, err := rowExistsInDB(db, "products", "sku", exampleSKU)

	assert.Nil(t, err)
	assert.True(t, exists)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestRowExistsInDBForNonexistentRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	exampleRows := sqlmock.NewRows([]string{""}).AddRow("false")
	skuExistenceQuery := buildProductExistenceQuery(exampleSKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(exampleSKU).
		WillReturnRows(exampleRows)

	exists, err := rowExistsInDB(db, "products", "sku", exampleSKU)

	assert.Nil(t, err)
	assert.False(t, exists)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
