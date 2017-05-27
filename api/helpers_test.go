package api

import (
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
)

var baseQueryFilter *QueryFilter

func init() {
	baseQueryFilter = &QueryFilter{
		Page:  1,
		Limit: 25,
	}
}

type RawFilterParamsTest struct {
	input      string
	expected   *QueryFilter
	shouldFail bool
}

func TestParseRawFilterParams(t *testing.T) {
	testSuite := []RawFilterParamsTest{
		RawFilterParamsTest{
			input:    "https://test.com/example",
			expected: baseQueryFilter,
		},
		RawFilterParamsTest{
			input:    "https://test.com/example?page=1&limit=25",
			expected: baseQueryFilter,
		},
		RawFilterParamsTest{
			input: "https://test.com/example?page=2&limit=40",
			expected: &QueryFilter{
				Page:  2,
				Limit: 40,
			},
		},
	}

	for i, test := range testSuite {
		earl, err := url.Parse(test.input)
		if err != nil {
			log.Fatal(err)
		}
		actual, err := parseRawFilterParams(earl.Query())
		assert.Nil(t, err)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("parseRawFilterParams test #%d (input: %s) should not fail", i, test.input))
	}

}

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
