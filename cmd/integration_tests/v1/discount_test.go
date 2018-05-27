package dairytest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dairycart/dairycart/models/v1"

	"github.com/stretchr/testify/assert"
)

func compareDiscounts(t *testing.T, expected, actual models.Discount) {
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name should be equal")
	assert.Equal(t, expected.DiscountType, actual.DiscountType, "expected and actual DiscountType should be equal")
	assert.Equal(t, expected.Amount, actual.Amount, "expected and actual Amount should be equal")
	assert.Equal(t, expected.RequiresCode, actual.RequiresCode, "expected and actual RequiresCode should be equal")
	assert.Equal(t, expected.Code, actual.Code, "expected and actual Code should be equal")
	assert.Equal(t, expected.LimitedUse, actual.LimitedUse, "expected and actual LimitedUse should be equal")
	assert.Equal(t, expected.NumberOfUses, actual.NumberOfUses, "expected and actual NumberOfUses should be equal")
	assert.Equal(t, expected.LoginRequired, actual.LoginRequired, "expected and actual LoginRequired should be equal")
}

func TestDiscountRetrievalRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		resp, err := getDiscountByID(existentID)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		expected := models.Discount{
			Name:          `10% off`,
			DiscountType:  "percentage",
			Amount:        10,
			RequiresCode:  false,
			LimitedUse:    false,
			LoginRequired: false,
		}
		var actual models.Discount

		unmarshalBody(_t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("for nonexistent discount", func(_t *testing.T) {
		resp, err := getDiscountByID(nonexistentID)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse

		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestDiscountListRoute(t *testing.T) {
	t.Parallel()

	t.Run("with standard filter", func(_t *testing.T) {
		resp, err := getListOfDiscounts(nil)
		assert.NoError(_t, err)
		assertStatusCode(_t, resp, http.StatusOK)

		expected := models.ListResponse{
			Limit: 25,
			Page:  1,
		}
		var actual models.ListResponse
		unmarshalBody(_t, resp, &actual)
		compareListResponses(t, expected, actual)
	})

	t.Run("with custom filter", func(_t *testing.T) {
		customFilter := map[string]string{
			"page":  "2",
			"limit": "5",
		}
		resp, err := getListOfDiscounts(customFilter)
		assert.NoError(_t, err)
		assertStatusCode(_t, resp, http.StatusOK)

		expected := models.ListResponse{
			Limit: 5,
			Page:  2,
		}
		var actual models.ListResponse
		unmarshalBody(_t, resp, &actual)
		compareListResponses(t, expected, actual)
	})
}

func TestDiscountCreationRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		expected := models.Discount{
			Name:          "test discount creation",
			DiscountType:  "percentage",
			Amount:        5,
			RequiresCode:  true,
			Code:          "discount code",
			LimitedUse:    true,
			NumberOfUses:  123,
			LoginRequired: true,
		}

		discountCreationJSON := createJSONBody(t, expected)
		resp, err := createDiscount(discountCreationJSON)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(_t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("with invalid input", func(_t *testing.T) {
		resp, err := createDiscount(exampleGarbageInput)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: expectedBadRequestResponse,
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestDiscountUpdateRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		exampleInput := models.Discount{
			Name:          "test discount update",
			DiscountType:  "percentage",
			Amount:        5,
			RequiresCode:  true,
			Code:          "discount code",
			LimitedUse:    true,
			NumberOfUses:  123,
			LoginRequired: true,
		}

		discountCreationJSON := createJSONBody(t, exampleInput)
		resp, err := createDiscount(discountCreationJSON)
		var expected models.Discount
		unmarshalBody(_t, resp, &expected)

		expected.Code = "this has changed now"
		updateJSON := createJSONBody(t, expected)
		resp, err = updateDiscount(expected.ID, updateJSON)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(_t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("with garbage input", func(_t *testing.T) {
		resp, err := updateDiscount(existentID, exampleGarbageInput)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: expectedBadRequestResponse,
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("for nonexistent discount", func(_t *testing.T) {
		exampleInput := createJSONBody(t, models.Discount{Name: "test nonexistent discount update"})

		resp, err := updateDiscount(nonexistentID, exampleInput)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestDiscountDeletionRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		exampleInput := models.Discount{
			Name:          "test discount update",
			DiscountType:  "percentage",
			Amount:        5,
			RequiresCode:  true,
			Code:          "discount code",
			LimitedUse:    true,
			NumberOfUses:  123,
			LoginRequired: true,
		}

		discountCreationJSON := createJSONBody(t, exampleInput)
		resp, err := createDiscount(discountCreationJSON)
		var expected models.Discount
		unmarshalBody(_t, resp, &expected)

		resp, err = deleteDiscount(expected.ID)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(_t, resp, &actual)
		assert.NotNil(t, actual.ArchivedOn)
	})

	t.Run("for nonexistent discount", func(_t *testing.T) {
		resp, err := deleteDiscount(nonexistentID)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}
