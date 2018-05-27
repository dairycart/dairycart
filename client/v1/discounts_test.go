package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
)

func buildNotFoundDiscountResponse(id uint64) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The discount you were looking for (id '%d') does not exist"
		}
	`, id)
}

func TestGetDiscountByID(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponse := loadExampleResponse(t, "discount")

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/discount/%d", existentID):    generateGetHandler(t, exampleResponse, http.StatusOK),
		fmt.Sprintf("/v1/discount/%d", nonexistentID): generateGetHandler(t, buildNotFoundDiscountResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		expected := &models.Discount{
			ID:            1,
			Name:          "10 percent off",
			DiscountType:  "percentage",
			Amount:        10,
			StartsOn:      buildTestTime(t),
			ExpiresOn:     buildTestDairytime(t),
			RequiresCode:  false,
			LimitedUse:    false,
			NumberOfUses:  0,
			LoginRequired: false,
			CreatedOn:     buildTestTime(t),
		}
		actual, err := c.GetDiscountByID(existentID)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected product doesn't match actual product")
	})

	t.Run("nonexistent product", func(*testing.T) {
		_, err := c.GetDiscountByID(nonexistentID)
		assert.NotNil(t, err)
	})
}

func TestGetDiscounts(t *testing.T) {
	exampleGoodResponse := loadExampleResponse(t, "discounts")

	t.Run("normal usage", func(*testing.T) {
		expected := []models.Discount{
			{
				ID:           1,
				Name:         "10 percent off",
				DiscountType: "percentage",
				Amount:       10,
				StartsOn:     buildTestTime(t),
				ExpiresOn:    buildTestDairytime(t),
				CreatedOn:    buildTestTime(t),
			},
			{
				ID:           2,
				Name:         "50 percent off",
				DiscountType: "percentage",
				Amount:       50,
				StartsOn:     buildTestTime(t),
				ExpiresOn:    buildTestDairytime(t),
				CreatedOn:    buildTestTime(t),
			},
			{
				ID:           3,
				Name:         "New customer special",
				DiscountType: "flat_amount",
				Amount:       10,
				StartsOn:     buildTestTime(t),
				ExpiresOn:    buildTestDairytime(t),
				CreatedOn:    buildTestTime(t),
			},
		}

		handlers := map[string]http.HandlerFunc{
			"/v1/discounts": generateGetHandler(t, exampleGoodResponse, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		actual, err := c.GetDiscounts(nil)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected discount list doesn't match actual product")
	})

	t.Run("with bad server response", func(*testing.T) {
		handlers := map[string]http.HandlerFunc{
			"/v1/discounts": generateGetHandler(t, exampleBadJSON, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		_, err := c.GetDiscounts(nil)
		assert.NotNil(t, err, "GetDiscounts should return an error when it receives nonsense")
	})
}

func TestCreateDiscount(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "discount")
	expectedBody := `
		{
			"name": "example_discount"
		}
	`
	exampleInput := models.DiscountCreationInput{
		Name: "example_discount",
	}

	t.Run("normal operation", func(*testing.T) {
		handlers := map[string]http.HandlerFunc{"/v1/discount": generatePostHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated)}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := &models.Discount{
			ID:            1,
			Name:          "10 percent off",
			DiscountType:  "percentage",
			Amount:        10,
			StartsOn:      buildTestTime(t),
			ExpiresOn:     buildTestDairytime(t),
			RequiresCode:  false,
			LimitedUse:    false,
			NumberOfUses:  0,
			LoginRequired: false,
			CreatedOn:     buildTestTime(t),
		}

		actual, err := c.CreateDiscount(exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("with bad response from server", func(*testing.T) {
		handlers := map[string]http.HandlerFunc{"/v1/discount": generatePostHandler(t, expectedBody, exampleBadJSON, http.StatusNotFound)}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.CreateDiscount(exampleInput)
		assert.NotNil(t, err)
	})
}

func TestUpdateDiscount(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "updated_discount")
	expectedBody := `
		{
			"name": "update_discount"
		}
	`
	exampleInput := models.DiscountUpdateInput{
		Name: "update_discount",
	}

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/discount/%d", existentID):    generatePatchHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		fmt.Sprintf("/v1/discount/%d", nonexistentID): generatePatchHandler(t, expectedBody, buildNotFoundDiscountResponse(nonexistentID), http.StatusNotFound),
	}
	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := &models.Discount{
			ID:            1,
			Name:          "update_discount",
			DiscountType:  "percentage",
			Amount:        10,
			StartsOn:      buildTestTime(t),
			ExpiresOn:     buildTestDairytime(t),
			RequiresCode:  false,
			LimitedUse:    false,
			NumberOfUses:  0,
			LoginRequired: false,
			CreatedOn:     buildTestTime(t),
			UpdatedOn:     buildTestDairytime(t),
		}

		actual, err := c.UpdateDiscount(existentID, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent discount", func(*testing.T) {
		_, err := c.UpdateDiscount(nonexistentID, exampleInput)
		assert.NotNil(t, err)
	})
}

func TestDeleteDiscount(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "deleted_discount")

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/discount/%d", existentID):    generateDeleteHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/discount/%d", nonexistentID): generateDeleteHandler(t, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent discount", func(*testing.T) {
		err := c.DeleteDiscount(existentID)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent discount", func(*testing.T) {
		err := c.DeleteDiscount(nonexistentID)
		assert.NotNil(t, err)
	})
}
