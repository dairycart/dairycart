package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discount/1", nil)
		assert.NoError(t, err)

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
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error retrieving discount count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscountCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving discount list", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscountCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetDiscountList", mock.Anything, mock.Anything).
			Return([]models.Discount{}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/discounts", nil)
		assert.NoError(t, err)

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
			Return(uint64(1), buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("CreateDiscount", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/discount", strings.NewReader(exampleDiscountCreationInput))
		assert.NoError(t, err)

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
			Return(buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("DeleteDiscount", mock.Anything, exampleDiscount.ID).
			Return(buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/discount/%d", exampleDiscount.ID), nil)
		assert.NoError(t, err)

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
			Return(buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating discount", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetDiscount", mock.Anything, exampleDiscount.ID).
			Return(exampleDiscount, nil)
		testUtil.MockDB.On("UpdateDiscount", mock.Anything, mock.Anything).
			Return(buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/discount/1", strings.NewReader(exampleDiscountUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
