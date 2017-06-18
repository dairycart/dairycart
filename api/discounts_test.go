package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	exampleDiscountStartTime     = "2016-12-01T12:00:00+05:00"
	exampleDiscountCreationInput = `
	{
		"name": "Test",
		"type": "flat_amount",
		"amount": 12.34,
		"starts_on": "2016-12-01T12:00:00+05:00",
		"requires_code": true,
		"code": "TEST"
	}`

	exampleDiscountUpdateInput = `
	{
		"name": "New Name",
		"requires_code": true,
		"code": "TEST"
	}`
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

func setExpectationsForDiscountCountQuery(mock sqlmock.Sqlmock, queryFilter *QueryFilter, err error) {
	exampleRows := sqlmock.NewRows([]string{"count"}).AddRow(3)

	discountListRetrievalQuery := buildCountQuery("discounts", queryFilter)
	query := formatQueryForSQLMock(discountListRetrievalQuery)
	mock.ExpectQuery(query).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountListQueryWithCount(mock sqlmock.Sqlmock, err error) {
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

func setExpectationsForDiscountListQuery(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(discountHeaders).
		AddRow(exampleDiscountData...).
		AddRow(exampleDiscountData...).
		AddRow(exampleDiscountData...)

	discountListRetrievalQuery, _ := buildDiscountListQuery(defaultQueryFilter)
	query := formatQueryForSQLMock(discountListRetrievalQuery)
	mock.ExpectQuery(query).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountCreation(mock sqlmock.Sqlmock, d *Discount, err error) {
	exampleRows := sqlmock.NewRows(discountHeaders).AddRow(exampleDiscountData...)
	discountCreationQuery, args := buildDiscountCreationQuery(d)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(discountCreationQuery)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountDeletion(mock sqlmock.Sqlmock, discountID string) {
	mock.ExpectExec(formatQueryForSQLMock(discountDeletionQuery)).
		WithArgs(discountID).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func setExpectationsForDiscountExistence(mock sqlmock.Sqlmock, discountID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(discountExistenceQuery)).
		WithArgs(discountID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountUpdate(mock sqlmock.Sqlmock, d *Discount, err error) {
	exampleRows := sqlmock.NewRows(discountHeaders).AddRow(exampleDiscountData...)
	query, rawArgs := buildDiscountUpdateQuery(d)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(args...).
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
	db, mock := setupDBForTest(t)
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, nil)

	actual, err := retrieveDiscountFromDB(db, discountIDString)
	assert.Nil(t, err)
	assert.Equal(t, *exampleDiscount, actual, "expected and actual discounts should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, sql.ErrNoRows)

	_, err := retrieveDiscountFromDB(db, discountIDString)
	assert.Equal(t, sql.ErrNoRows, err, "retrieveDiscountFromDB should return errors it receives")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, discountIDString, arbitraryError)

	_, err := retrieveDiscountFromDB(db, discountIDString)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestValidateDiscountCreationInput(t *testing.T) {
	t.Parallel()
	dummyTime, _ := time.Parse("2006-01-02T15:04:05-07:00", exampleDiscountStartTime)
	expected := &Discount{
		Name:         "Test",
		Type:         "flat_amount",
		Amount:       12.34,
		StartsOn:     dummyTime,
		RequiresCode: true,
		Code:         "TEST",
	}

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleDiscountCreationInput))
	actual, err := validateDiscountCreationInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "valid discount creation input should parse into a proper discount creation struct")
}

func TestValidateDiscountCreationInputWithNoInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateDiscountCreationInput(req)

	assert.NotNil(t, err)
}

func TestValidateDiscountCreationInputWithInvalidInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleGarbageInput))
	_, err := validateDiscountCreationInput(req)

	assert.NotNil(t, err)
}

func TestCreateDiscountInDB(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	setExpectationsForDiscountCreation(mock, exampleDiscount, nil)

	actualDiscount, err := createDiscountInDB(db, exampleDiscount)
	assert.Nil(t, err)
	assert.Equal(t, exampleDiscount, actualDiscount, "createProductInDB should return the created Discount")

	ensureExpectationsWereMet(t, mock)
}

func TestArchiveDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountDeletion(mock, exampleDiscountID)

	err := archiveDiscount(db, exampleDiscountID)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestValidateDiscountUpdateInput(t *testing.T) {
	t.Parallel()

	expected := &Discount{
		Name:         "New Name",
		RequiresCode: true,
		Code:         "TEST",
	}

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleDiscountUpdateInput))
	actual, err := validateDiscountUpdateInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "valid discount creation input should parse into a proper discount creation struct")
}

func TestValidateDiscountUpdateInputWithNoInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateDiscountUpdateInput(req)

	assert.NotNil(t, err)
}

func TestValidateDiscountUpdateInputWithInvalidInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleGarbageInput))
	_, err := validateDiscountUpdateInput(req)

	assert.NotNil(t, err)
}

func TestUpdateDiscountInDB(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	setExpectationsForDiscountUpdate(mock, exampleDiscount, nil)

	err := updateDiscountInDatabase(db, exampleDiscount)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestDiscountRetrievalHandlerWithExistingDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
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
	db, mock := setupDBForTest(t)
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
	db, mock := setupDBForTest(t)
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
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForDiscountCountQuery(mock, defaultQueryFilter, nil)
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

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromCountQuery(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForDiscountCountQuery(mock, defaultQueryFilter, arbitraryError)

	req, err := http.NewRequest("GET", "/v1/discounts", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromListQuery(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForDiscountCountQuery(mock, defaultQueryFilter, nil)
	setExpectationsForDiscountListQuery(mock, arbitraryError)

	req, err := http.NewRequest("GET", "/v1/discounts", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountCreationHandler(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	dummyTime, _ := time.Parse("2006-01-02T15:04:05-07:00", exampleDiscountStartTime)
	exampleCreatedDiscount := &Discount{
		ID:           1,
		Name:         "Test",
		Type:         "flat_amount",
		Amount:       12.34,
		StartsOn:     dummyTime,
		RequiresCode: true,
		Code:         "TEST",
		CreatedOn:    exampleTime,
	}

	setExpectationsForDiscountCreation(mock, exampleCreatedDiscount, nil)
	req, err := http.NewRequest("POST", "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 201, res.Code, "status code should be 201")

	actual := &Discount{}
	err = json.NewDecoder(strings.NewReader(res.Body.String())).Decode(actual)
	assert.Nil(t, err)

	expected := exampleDiscount
	expected.UpdatedOn.Valid = true
	expected.ArchivedOn.Valid = true
	assert.Equal(t, expected, actual, "discount creation endpoint should return created discount")

	ensureExpectationsWereMet(t, mock)
}

func TestDiscountCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	db, _ := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	req, err := http.NewRequest("POST", "/v1/discount", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "status code should be 400")
}

func TestDiscountCreationHandlerWithDatabaseErrorUponCreation(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForDiscountCreation(mock, exampleDiscount, arbitraryError)
	req, err := http.NewRequest("POST", "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
}

func TestDiscountDeletionHandler(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountExistence(mock, exampleDiscountID, true, nil)
	setExpectationsForDiscountDeletion(mock, exampleDiscountID)

	req, err := http.NewRequest("DELETE", "/v1/discount/1", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountDeletionHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountExistence(mock, exampleDiscountID, false, nil)

	req, err := http.NewRequest("DELETE", "/v1/discount/1", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountUpdateHandler(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	updateInput := &Discount{
		ID:           1,
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
		CreatedOn:    exampleTime,
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(mock, updateInput, nil)

	req, err := http.NewRequest("PUT", "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountUpdateHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)
	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(mock, exampleDiscountID, sql.ErrNoRows)

	req, err := http.NewRequest("PUT", "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountUpdateHandlerWithErrorValidatingInput(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	req, err := http.NewRequest("PUT", "/v1/discount/1", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountUpdateHandlerWithErrorRetrievingDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	updateInput := &Discount{
		ID:           1,
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
		CreatedOn:    exampleTime,
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(mock, exampleDiscountID, arbitraryError)

	req, err := http.NewRequest("PUT", "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestDiscountUpdateHandlerWithErrorUpdatingDiscount(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	updateInput := &Discount{
		ID:           1,
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
		CreatedOn:    exampleTime,
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(mock, updateInput, arbitraryError)

	req, err := http.NewRequest("PUT", "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
