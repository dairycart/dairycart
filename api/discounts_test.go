package api

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"strconv"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
)

var discountHeaders []string
var exampleDiscountData []driver.Value
var exampleDiscount *Discount

func init() {
	exampleDiscount = &Discount{
		ID:        1,
		Name:      "Example Discount",
		Type:      "flat_discount",
		ProductID: 10, // == exampleProduct.ID
		Amount:    12.34,
		StartsOn:  exampleTime,
		ExpiresOn: exampleTime.Add(30 * (24 * time.Hour)),
		CreatedAt: exampleTime,
	}

	discountHeaders = []string{"id", "name", "type", "amount", "product_id", "starts_on", "expires_on", "created_at", "updated_at", "archived_at"}
	exampleDiscountData = []driver.Value{exampleDiscount.ID, exampleDiscount.Name, exampleDiscount.Type, exampleDiscount.Amount, exampleDiscount.ProductID, exampleDiscount.StartsOn, exampleDiscount.ExpiresOn, exampleDiscount.CreatedAt, nil, nil}
}

func setExpectationsForDiscountRetrievalByID(mock sqlmock.Sqlmock, id string, err error) {
	exampleRows := sqlmock.NewRows(discountHeaders).AddRow(exampleDiscountData...)
	discountQuery := buildDiscountRetrievalQuery(exampleSKU)
	query := formatQueryForSQLMock(discountQuery)
	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestDiscountGenerateScanArgs(t *testing.T) {
	// this test will go away soon, but for now I can't stand having less than 100% test coverage
	t.Parallel()
	d := &Discount{}
	actual := d.generateScanArgs()
	assert.Equal(t, 10, len(actual), "there should be eight scan args for discounts")
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
