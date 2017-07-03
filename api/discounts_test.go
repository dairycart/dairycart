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

var (
	discountHeaders              []string
	exampleDiscountData          []driver.Value
	discountHeadersWithCount     []string
	exampleDiscountDataWithCount []driver.Value
	exampleDiscount              *Discount
)

func init() {
	exampleDiscount = &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_discount",
		Amount:    12.34,
		StartsOn:  generateExampleTimeForTests(),
		ExpiresOn: NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

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
		nil,
	}

	discountHeaders = strings.Split(strings.TrimSpace(discountsTableColumns), ",\n\t\t")
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
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, nil)

	expectedDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_discount",
		Amount:    12.34,
		StartsOn:  generateExampleTimeForTests(),
		ExpiresOn: NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	actual, err := retrieveDiscountFromDB(testUtil.DB, discountIDString)
	assert.Nil(t, err)
	assert.Equal(t, *expectedDiscount, actual, "expected and actual discounts should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, sql.ErrNoRows)

	_, err := retrieveDiscountFromDB(testUtil.DB, discountIDString)
	assert.Equal(t, sql.ErrNoRows, err, "retrieveDiscountFromDB should return errors it receives")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, arbitraryError)

	_, err := retrieveDiscountFromDB(testUtil.DB, discountIDString)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateDiscountInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCreation(testUtil.Mock, exampleDiscount, nil)

	actualDiscount, err := createDiscountInDB(testUtil.DB, exampleDiscount)
	assert.Nil(t, err)
	assert.Equal(t, exampleDiscount, actualDiscount, "createProductInDB should return the created Discount")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID)

	err := archiveDiscount(testUtil.DB, exampleDiscountID)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateDiscountInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountUpdate(testUtil.Mock, exampleDiscount, nil)

	err := updateDiscountInDatabase(testUtil.DB, exampleDiscount)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestDiscountRetrievalHandlerWithExistingDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountRetrievalHandlerWithNoRowsFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountRetrievalHandlerWithDBError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, defaultQueryFilter, nil)
	setExpectationsForDiscountListQuery(testUtil.Mock, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")

	expected := &DiscountsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &DiscountsResponse{}
	err = json.NewDecoder(strings.NewReader(testUtil.Response.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	assert.Equal(t, uint64(len(actual.Data)), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromCountQuery(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, defaultQueryFilter, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromListQuery(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, defaultQueryFilter, nil)
	setExpectationsForDiscountListQuery(testUtil.Mock, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountCreationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	dummyTime, _ := time.Parse("2006-01-02T15:04:05-07:00", exampleDiscountStartTime)
	exampleCreatedDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "Test",
		Type:         "flat_amount",
		Amount:       12.34,
		StartsOn:     dummyTime,
		RequiresCode: true,
		Code:         "TEST",
	}

	setExpectationsForDiscountCreation(testUtil.Mock, exampleCreatedDiscount, nil)
	req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")

	actual := &Discount{}
	err = json.NewDecoder(strings.NewReader(testUtil.Response.Body.String())).Decode(actual)
	assert.Nil(t, err)

	expected := exampleDiscount
	expected.UpdatedOn.Valid = true
	expected.ArchivedOn.Valid = true
	assert.Equal(t, expected, actual, "discount creation endpoint should return created discount")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
}

func TestDiscountCreationHandlerWithDatabaseErrorUponCreation(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCreation(testUtil.Mock, exampleDiscount, arbitraryError)
	req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	//ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountExistence(testUtil.Mock, exampleDiscountID, true, nil)
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID)

	req, err := http.NewRequest(http.MethodDelete, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountDeletionHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountExistence(testUtil.Mock, exampleDiscountID, false, nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	updateInput := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(testUtil.Mock, updateInput, nil)

	req, err := http.NewRequest(http.MethodPut, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscount.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodPut, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerWithErrorValidatingInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPut, "/v1/discount/1", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerWithErrorRetrievingDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	updateInput := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, arbitraryError)

	req, err := http.NewRequest(http.MethodPut, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerWithErrorUpdatingDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	updateInput := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "New Name",
		Type:         "flat_discount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(testUtil.Mock, updateInput, arbitraryError)

	req, err := http.NewRequest(http.MethodPut, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
