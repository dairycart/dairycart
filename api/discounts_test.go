package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var discountHeaders []string
var exampleDiscountData []driver.Value
var discountHeadersWithCount []string
var exampleDiscountDataWithCount []driver.Value
var exampleDiscount *Discount

func init() {
	exampleDiscount = &Discount{
		ID:        1,
		Name:      "Example Discount",
		Type:      "flat_discount",
		Amount:    12.34,
		StartsOn:  exampleTime,
		ExpiresOn: NullTime{pq.NullTime{Time: exampleTime.Add(30 * (24 * time.Hour)), Valid: true}},
		CreatedOn: exampleTime,
	}

	discountHeaders = []string{"id", "name", "type", "amount", "starts_on", "expires_on", "requires_code", "code", "limited_use", "number_of_uses", "login_required", "created_on", "updated_on", "archived_on"}
	exampleDiscountData = []driver.Value{
		exampleDiscount.ID,
		exampleDiscount.Name,
		exampleDiscount.Type,
		exampleDiscount.Amount,
		exampleDiscount.StartsOn,
		exampleDiscount.ExpiresOn.Time,
		exampleDiscount.RequiresCode,
		exampleDiscount.Code,
		exampleDiscount.LimitedUse,
		exampleDiscount.NumberOfUses,
		exampleDiscount.LoginRequired,
		exampleDiscount.CreatedOn,
		nil,
		nil}

	discountHeadersWithCount = append([]string{"count"}, discountHeaders...)
	exampleDiscountDataWithCount = append([]driver.Value{3}, exampleDiscountData...)
}

func setExpectationsForDiscountRetrievalByID(mock sqlmock.Sqlmock, id string, err error) {
	exampleRows := sqlmock.NewRows(discountHeaders).AddRow(exampleDiscountData...)
	query := formatQueryForSQLMock(discountRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountListQuery(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(discountHeadersWithCount).
		AddRow(exampleDiscountDataWithCount...).
		AddRow(exampleDiscountDataWithCount...).
		AddRow(exampleDiscountDataWithCount...)

	discountListRetrievalQuery, _ := buildDiscountListQuery(defaultQueryFilter)
	query := formatQueryForSQLMock(discountListRetrievalQuery)
	mock.ExpectQuery(query).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestDiscountTypeIsValidWithValidInput(t *testing.T) {
	t.Parallel()
	d := &Discount{Type: "flat_discount"}
	assert.False(t, d.discountTypeIsValid())
}

func TestDiscountTypeIsValidWithInvalidInput(t *testing.T) {
	t.Parallel()
	d := &Discount{Type: "this is nonsense"}
	assert.False(t, d.discountTypeIsValid())
}

func TestRetrieveDiscountFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, nil)

	actual, err := retrieveDiscountFromDB(db, discountIDString)
	assert.Nil(t, err)
	assert.Equal(t, exampleDiscount, actual, "expected and actual discounts should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, sql.ErrNoRows)

	_, err = retrieveDiscountFromDB(db, discountIDString)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, arbitraryError)

	_, err = retrieveDiscountFromDB(db, discountIDString)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountsFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForDiscountListQuery(mock, nil)

	discounts, count, err := retrieveListOfDiscountsFromDB(db, defaultQueryFilter)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(discounts), "there should be 3 discounts")
	assert.Equal(t, uint64(3), count, "there should be 3 discounts")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountsFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForDiscountListQuery(mock, sql.ErrNoRows)

	_, count, err := retrieveListOfDiscountsFromDB(db, defaultQueryFilter)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), count, "count returned should be zero when error is encountered")
	ensureExpectationsWereMet(t, mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestDiscountRetrievalHandlerWithExistingDiscount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, nil)

	req, err := http.NewRequest("GET", "/v1/discount/1", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountRetrievalHandlerWithNoRowsFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/v1/discount/1", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountRetrievalHandlerWithDBError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, arbitraryError)

	req, err := http.NewRequest("GET", "/v1/discount/1", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountListHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForDiscountListQuery(mock, nil)

	req, err := http.NewRequest("GET", "/v1/discounts", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")

	expected := &DiscountsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &DiscountsResponse{}
	err = json.NewDecoder(strings.NewReader(res.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	assert.Equal(t, uint64(len(actual.Data)), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountListHandlerWithDBError(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForDiscountListQuery(mock, arbitraryError)

	req, err := http.NewRequest("GET", "/v1/discounts", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
