package api

import (
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
)

const ()

var exampleFilterStartTime time.Time
var exampleFilterEndTime time.Time
var defaultQueryFilter *QueryFilter

func init() {
	defaultQueryFilter = &QueryFilter{
		Page:  1,
		Limit: 25,
	}
}

type RawFilterParamsTest struct {
	input          string
	expected       *QueryFilter
	shouldFail     bool
	failureMessage string
}

func TestParseRawFilterParams(t *testing.T) {
	exampleUnixStartTime := int64(232747200)
	exampleUnixEndTime := int64(232747200 + 10000)

	exampleFilterStartTime := time.Unix(exampleUnixStartTime, 0)
	exampleFilterEndTime := time.Unix(exampleUnixEndTime, 0)

	testSuite := []RawFilterParamsTest{
		RawFilterParamsTest{
			input:          "https://test.com/example",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no query params should parse to the default query filter",
		},
		RawFilterParamsTest{
			input:          "https://test.com/example?page=1&limit=25",
			expected:       defaultQueryFilter,
			failureMessage: "URL with query params set to the defaults should parse to the default query filter",
		},
		RawFilterParamsTest{
			input: "https://test.com/example?page=1&limit=500000",
			expected: &QueryFilter{
				Page:  1,
				Limit: 50,
			},
			failureMessage: "URL with limit param set to high should default to 50",
		},
		RawFilterParamsTest{
			input: "https://test.com/example?page=2&limit=40",
			expected: &QueryFilter{
				Page:  2,
				Limit: 40,
			},
			failureMessage: "URL with non-default page and limit params should parse correctly",
		},
		RawFilterParamsTest{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v", exampleUnixStartTime),
			expected: &QueryFilter{
				Page:         1,
				Limit:        25,
				UpdatedAfter: exampleFilterStartTime,
			},
			failureMessage: "URL with specified updated_after field should have a non-nil time value set for UpdatedAfter",
		},
		RawFilterParamsTest{
			input: fmt.Sprintf("https://test.com/example?updated_before=%v", exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_before field should have a non-nil time value set for UpdatedBefore",
		},
		RawFilterParamsTest{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_after and updated_before fields should have a non-nil time value set for both UpdatedAfter and UpdatedBefore",
		},
		RawFilterParamsTest{
			input: fmt.Sprintf("https://test.com/example?page=2&limit=35&updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          2,
				Limit:         35,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with all relevant filters should have a completely custom QueryFilter value",
		},
		RawFilterParamsTest{
			input:          fmt.Sprintf("https://test.com/example?rage=2&dimit=35&upgraded_after=%v&agitated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
	}

	for _, test := range testSuite {
		earl, err := url.Parse(test.input)
		if err != nil {
			log.Fatal(err)
		}
		actual, err := parseRawFilterParams(earl.Query())
		if !test.shouldFail {
			assert.Nil(t, err)
		}
		assert.Equal(t, test.expected, actual, test.failureMessage)
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
