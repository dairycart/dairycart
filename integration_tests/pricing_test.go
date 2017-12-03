package dairytest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

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

	t.Run("normal usage", func(*testing.T) {
		resp, err := getDiscountByID(existentID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		expected := models.Discount{
			Name:          `10% off`,
			DiscountType:  "percentage",
			Amount:        10,
			RequiresCode:  false,
			LimitedUse:    false,
			LoginRequired: false,
		}
		var actual models.Discount

		unmarshalBody(t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("for nonexistent discount", func(*testing.T) {
		resp, err := getDiscountByID(nonexistentID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse

		unmarshalBody(t, resp, &actual)
		assert.Equal(t, expected, actual)
	})
}

func TestDiscountListRoute(t *testing.T) {
	t.Parallel()

	t.Run("with standard filter", func(*testing.T) {
		resp, err := getListOfDiscounts(nil)
		assert.NoError(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		expected := models.ListResponse{
			Limit: 25,
			Page:  1,
		}
		var actual models.ListResponse
		unmarshalBody(t, resp, &actual)
		compareListResponses(t, expected, actual)
	})

	t.Run("with custom filter", func(*testing.T) {
		customFilter := map[string]string{
			"page":  "2",
			"limit": "5",
		}
		resp, err := getListOfDiscounts(customFilter)
		assert.NoError(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		expected := models.ListResponse{
			Limit: 5,
			Page:  2,
		}
		var actual models.ListResponse
		unmarshalBody(t, resp, &actual)
		compareListResponses(t, expected, actual)
	})
}

func TestDiscountCreationRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(*testing.T) {
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
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("with invalid input", func(*testing.T) {
		resp, err := createDiscount(exampleGarbageInput)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: expectedBadRequestResponse,
		}
		var actual models.ErrorResponse
		unmarshalBody(t, resp, &actual)
		assert.Equal(t, expected, actual)
	})
}

func TestDiscountUpdateRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(*testing.T) {
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
		unmarshalBody(t, resp, &expected)

		expected.Code = "this has changed now"
		updateJSON := createJSONBody(t, expected)
		resp, err = updateDiscount(expected.ID, updateJSON)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(t, resp, &actual)
		compareDiscounts(t, expected, actual)
	})

	t.Run("with garbage input", func(*testing.T) {
		resp, err := updateDiscount(existentID, exampleGarbageInput)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: expectedBadRequestResponse,
		}
		var actual models.ErrorResponse
		unmarshalBody(t, resp, &actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent discount", func(*testing.T) {
		exampleInput := createJSONBody(t, models.Discount{Name: "test nonexistent discount update"})

		resp, err := updateDiscount(nonexistentID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse
		unmarshalBody(t, resp, &actual)
		assert.Equal(t, expected, actual)
	})
}

func TestDiscountDeletionRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(*testing.T) {
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
		unmarshalBody(t, resp, &expected)

		resp, err = deleteDiscount(expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actual models.Discount
		unmarshalBody(t, resp, &actual)
		assert.True(t, actual.ArchivedOn.Valid)
	})

	t.Run("for nonexistent discount", func(*testing.T) {
		resp, err := deleteDiscount(nonexistentID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The discount you were looking for (id '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse
		unmarshalBody(t, resp, &actual)
		assert.Equal(t, expected, actual)
	})
}
