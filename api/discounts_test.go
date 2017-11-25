package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	// "encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleDiscountID = 1
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

func setExpectationsForDiscountCreation(mock sqlmock.Sqlmock, d *models.Discount, err error) {
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

func setExpectationsForDiscountUpdate(mock sqlmock.Sqlmock, d *models.Discount, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	query, rawArgs := buildDiscountUpdateQuery(d)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestRetrieveDiscountFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	discountIDString := strconv.Itoa(int(exampleDiscountID))
	setExpectationsForDiscountRetrievalByID(testUtil.Mock, discountIDString, nil)

	expectedDiscount := &models.Discount{
		ID:           1,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    models.NullTime{NullTime: pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
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
	exampleDiscount := &models.Discount{
		ID:           1,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    models.NullTime{NullTime: pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
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
	exampleDiscount := &models.Discount{
		ID:           1,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    models.NullTime{NullTime: pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
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

func TestDiscountRetrievalHandler(t *testing.T) {
	exampleDiscount := &models.Discount{
		ID:           1,
		Name:         "example",
		DiscountType: "percentage",
		Amount:       12.34,
		Code:         "example",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestDiscountListHandler(t *testing.T) {
	exampleDiscount := models.Discount{
		ID:           1,
		Name:         "example",
		DiscountType: "percentage",
		Amount:       12.34,
		Code:         "example",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscountCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetDiscountList", mock.Anything, mock.Anything).
			Return([]models.Discount{exampleDiscount}, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error retrieving discount count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscountCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving discount list", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscountCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetDiscountList", mock.Anything, mock.Anything).
			Return([]models.Discount{}, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestDiscountCreationHandler(t *testing.T) {
	exampleDiscountCreationInput := `
		{
			"name": "Test",
			"discount_type": "flat_amount",
			"amount": 12.34,
			"starts_on": "2016-12-01T12:00:00+05:00",
			"requires_code": true,
			"code": "TEST"
		}
	`

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("CreateDiscount", mock.Anything, mock.Anything).
			Return(uint64(1), generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("CreateDiscount", mock.Anything, mock.Anything).
			Return(uint64(1), generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestDiscountDeletionHandler(t *testing.T) {
	exampleDiscount := &models.Discount{
		ID:           1,
		Name:         "example",
		DiscountType: "percentage",
		Amount:       12.34,
		Code:         "example",
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("DeleteDiscount", mock.Anything, exampleDiscount.ID).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("DeleteDiscount", mock.Anything, exampleDiscount.ID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestDiscountUpdateHandler(t *testing.T) {
	exampleDiscount := &models.Discount{
		ID:           1,
		Name:         "example",
		DiscountType: "percentage",
		Amount:       12.34,
		Code:         "example",
	}

	exampleDiscountUpdateInput := `
		{
			"name": "New Name",
			"requires_code": true,
			"code": "TEST"
		}
	`

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("UpdateDiscount", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("UpdateDiscount", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
