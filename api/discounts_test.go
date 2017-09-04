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
	exampleDiscountID            = 1
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

func setExpectationsForDiscountRetrievalByID(mock sqlmock.Sqlmock, id string, err error) {
	exampleDiscountReadData := []driver.Value{
		1,
		"Example Discount",
		"flat_amount",
		12.34,
		generateExampleTimeForTests(),
		generateExampleTimeForTests().Add(30 * (24 * time.Hour)),
		false,
		"",
		false,
		0,
		false,
		generateExampleTimeForTests(),
		nil,
		nil,
	}

	exampleRows := sqlmock.NewRows(strings.Split(strings.TrimSpace(discountsTableColumns), ",\n\t\t")).AddRow(exampleDiscountReadData...)
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

func setExpectationsForDiscountListQuery(mock sqlmock.Sqlmock, err error) {
	exampleDiscountReadData := []driver.Value{
		1,
		"Example Discount",
		"flat_amount",
		12.34,
		generateExampleTimeForTests(),
		generateExampleTimeForTests().Add(30 * (24 * time.Hour)),
		false,
		"",
		false,
		0,
		false,
		generateExampleTimeForTests(),
		nil,
		nil,
	}

	exampleRows := sqlmock.NewRows(strings.Split(strings.TrimSpace(discountsTableColumns), ",\n\t\t")).
		AddRow(exampleDiscountReadData...).
		AddRow(exampleDiscountReadData...).
		AddRow(exampleDiscountReadData...)

	discountListRetrievalQuery, _ := buildDiscountListQuery(genereateDefaultQueryFilter())
	query := formatQueryForSQLMock(discountListRetrievalQuery)
	mock.ExpectQuery(query).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountCreation(mock sqlmock.Sqlmock, d *Discount, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(d.ID, d.CreatedOn)
	discountCreationQuery, args := buildDiscountCreationQuery(d)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(discountCreationQuery)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountDeletion(mock sqlmock.Sqlmock, discountID string, err error) {
	mock.ExpectExec(formatQueryForSQLMock(discountDeletionQuery)).
		WithArgs(discountID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForDiscountExistence(mock sqlmock.Sqlmock, discountID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(discountExistenceQuery)).
		WithArgs(discountID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForDiscountUpdate(mock sqlmock.Sqlmock, d *Discount, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	query, rawArgs := buildDiscountUpdateQuery(d)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestDiscountTypeIsValidWithValidInput(t *testing.T) {
	t.Parallel()
	d := &Discount{Type: "invalid_type"}
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

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, nil)

	expectedDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_amount",
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

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, sql.ErrNoRows)

	_, err := retrieveDiscountFromDB(testUtil.DB, discountIDString)
	assert.Equal(t, sql.ErrNoRows, err, "retrieveDiscountFromDB should return errors it receives")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveDiscountFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, generateArbitraryError())

	_, err := retrieveDiscountFromDB(testUtil.DB, discountIDString)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateDiscountInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_amount",
		Amount:    12.34,
		StartsOn:  generateExampleTimeForTests(),
		ExpiresOn: NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	setExpectationsForDiscountCreation(testUtil.Mock, exampleDiscount, nil)

	discountID, createdOn, err := createDiscountInDB(testUtil.DB, exampleDiscount)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), discountID, "createDiscountInDB should return the created discount's ID")
	assert.Equal(t, generateExampleTimeForTests(), createdOn, "createDiscountInDB should return the created discount's creation time")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID, nil)

	err := archiveDiscount(testUtil.DB, exampleDiscountID)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveDiscountReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID, generateArbitraryError())

	err := archiveDiscount(testUtil.DB, exampleDiscountID)
	assert.Equal(t, err, generateArbitraryError(), "archiveDiscount should return an error when it encounters one")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateDiscountInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_amount",
		Amount:    12.34,
		StartsOn:  generateExampleTimeForTests(),
		ExpiresOn: NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	setExpectationsForDiscountUpdate(testUtil.Mock, exampleDiscount, nil)

	updatedTime, err := updateDiscountInDatabase(testUtil.DB, exampleDiscount)
	assert.Nil(t, err)

	assert.Equal(t, updatedTime, generateExampleTimeForTests(), "updateDiscountInDatabase should return the appropriate time")
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

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountRetrievalHandlerWithNoRowsFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountRetrievalHandlerWithDBError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, genereateDefaultQueryFilter(), nil)
	setExpectationsForDiscountListQuery(testUtil.Mock, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)

	expected := &ListResponse{
		Page:  1,
		Limit: 25,
		Count: 3,
	}

	actual := &ListResponse{}
	err = json.NewDecoder(strings.NewReader(testUtil.Response.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromCountQuery(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, genereateDefaultQueryFilter(), generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountListHandlerWithDBErrorWithErrorReturnedFromListQuery(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForDiscountCountQuery(testUtil.Mock, genereateDefaultQueryFilter(), nil)
	setExpectationsForDiscountListQuery(testUtil.Mock, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
	assertStatusCode(t, testUtil, http.StatusCreated)

	actual := &Discount{}
	bodyStr := testUtil.Response.Body.String()
	err = json.NewDecoder(strings.NewReader(bodyStr)).Decode(actual)
	assert.Nil(t, err)

	actual.UpdatedOn.Valid = false
	actual.ArchivedOn.Valid = false
	actual.ExpiresOn.Valid = false
	assert.Equal(t, exampleCreatedDiscount, actual, "discount creation endpoint should return created discount")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusBadRequest)
}

func TestDiscountCreationHandlerWithDatabaseErrorUponCreation(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:      "Example Discount",
		Type:      "flat_amount",
		Amount:    12.34,
		StartsOn:  generateExampleTimeForTests(),
		ExpiresOn: NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	setExpectationsForDiscountCreation(testUtil.Mock, exampleDiscount, generateArbitraryError())
	req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	//ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountExistence(testUtil.Mock, exampleDiscountID, true, nil)
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID, nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountDeletionHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountExistence(testUtil.Mock, exampleDiscountID, false, nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountDeletionHandlerWithErrorDeletingDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountExistence(testUtil.Mock, exampleDiscountID, true, nil)
	setExpectationsForDiscountDeletion(testUtil.Mock, exampleDiscountID, generateArbitraryError())

	req, err := http.NewRequest(http.MethodDelete, "/v1/discount/1", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		Type:         "flat_amount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(testUtil.Mock, updateInput, nil)

	req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleDiscountID := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDiscountUpdateHandlerWithErrorValidatingInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusBadRequest)
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
		Type:         "flat_amount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, generateArbitraryError())

	req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		Type:         "flat_amount",
		Amount:       12.34,
		RequiresCode: true,
		Code:         "TEST",
	}

	exampleDiscountID := strconv.Itoa(int(updateInput.ID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, exampleDiscountID, nil)
	setExpectationsForDiscountUpdate(testUtil.Mock, updateInput, generateArbitraryError())

	req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}
