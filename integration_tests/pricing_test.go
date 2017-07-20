package dairytest

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func replaceTimeStringsForDiscountTests(body string) string {
	re := regexp.MustCompile(`(?U)(,?)"(starts_on|expires_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func createDiscountBody(code string) string {
	output := fmt.Sprintf(`
		{
			"name": "Test",
			"type": "flat_amount",
			"amount": 12.34,
			"starts_on": "2016-12-01T12:00:00+05:00",
			"requires_code": true,
			"code": "%s"
		}
	`, code)
	return output
}

func TestDiscountRetrievalForExistingDiscount(t *testing.T) {
	// FIXME: parallelize
	resp, err := getDiscountByID(existentID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "a successfully retrieved discount should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	removedTimeFields := replaceTimeStringsForTests(body)
	actual := replaceTimeStringsForDiscountTests(removedTimeFields)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "retrieved"))
	assert.Equal(t, expected, actual, "discount route should return a serialized discount object")
}

func TestDiscountRetrievalForNonexistentDiscount(t *testing.T) {
	t.Parallel()
	resp, err := getDiscountByID(nonexistentID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "a request for a nonexistent discount should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "error_discount_does_not_exist"))
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")

}

func TestDiscountListRetrievalWithDefaultFilter(t *testing.T) {
	// FIXME: parallelize
	resp, err := getListOfDiscounts(nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	respBody := turnResponseBodyIntoString(t, resp)
	removedTimeFields := replaceTimeStringsForTests(respBody)
	actual := replaceTimeStringsForDiscountTests(removedTimeFields)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "list_with_default_filter"))
	assert.Equal(t, expected, actual, "product list route should respond with a list of products")
}

func TestDiscountListRouteWithCustomFilter(t *testing.T) {
	// FIXME: parallelize
	customFilter := map[string]string{
		"page":  "2",
		"limit": "2",
	}
	resp, err := getListOfDiscounts(customFilter)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	respBody := turnResponseBodyIntoString(t, resp)
	removedTimeFields := replaceTimeStringsForTests(respBody)
	actual := replaceTimeStringsForDiscountTests(removedTimeFields)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "list_with_custom_filter"))
	assert.Equal(t, expected, actual, "product list route should respond with a customized list of products")
}

func TestDiscountCreation(t *testing.T) {
	t.Parallel()

	var createdDiscountID uint64
	testDiscountCode := "TEST"

	testCreateDiscount := func(t *testing.T) {
		newDiscountJSON := createDiscountBody(testDiscountCode)
		resp, err := createDiscount(newDiscountJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a discount that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdDiscountID = retrieveIDFromResponseBody(body, t)

		actual := replaceTimeStringsForTests(body)
		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"id": %d,
				"name": "Test",
				"type": "flat_amount",
				"amount": 12.34,
				"starts_on": "2016-12-01T12:00:00Z",
				"expires_on": "",
				"requires_code": true,
				"code": "%s",
				"limited_use": false,
				"login_required": false
			}
		`, createdDiscountID, testDiscountCode))

		assert.Equal(t, expected, actual, "discount creation route should respond with created product body")
	}

	testDeleteDiscount := func(t *testing.T) {
		resp, err := deleteDiscount(strconv.Itoa(int(createdDiscountID)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a discount that exists should respond 200")
	}

	subtests := []subtest{
		subtest{
			Message: "create discount",
			Test:    testCreateDiscount,
		},
		subtest{
			Message: "delete created discount",
			Test:    testDeleteDiscount,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestDiscountDeletion(t *testing.T) {
	t.Parallel()

	var createdDiscountID uint64
	testDiscountCode := "deletion"

	testCreateDiscount := func(t *testing.T) {
		newDiscountJSON := createDiscountBody(testDiscountCode)
		resp, err := createDiscount(newDiscountJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a discount that doesn't exist should respond 201")
		body := turnResponseBodyIntoString(t, resp)
		createdDiscountID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteDiscount := func(t *testing.T) {
		resp, err := deleteDiscount(strconv.Itoa(int(createdDiscountID)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a discount that exists should respond 200")

		actual := turnResponseBodyIntoString(t, resp)
		expected := fmt.Sprintf("Successfully archived discount `%d`", createdDiscountID)
		assert.Equal(t, expected, actual, "discount deletion route should respond with affirmative message upon successful deletion")
	}

	subtests := []subtest{
		subtest{
			Message: "create discount",
			Test:    testCreateDiscount,
		},
		subtest{
			Message: "delete created discount",
			Test:    testDeleteDiscount,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestDiscountCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createDiscount(exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a discount that doesn't exist should respond 400")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	assert.Equal(t, expectedBadRequestResponse, actual, "discount creation route should respond with created product body")
}

func TestDiscountUpdate(t *testing.T) {
	// FIXME: parallelize
	updatedDiscountJSON := loadExampleInput(t, "discounts", "update")
	resp, err := updateDiscount(existentID, updatedDiscountJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "updated"))
	assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
}

func TestDiscountUpdateInvalidDiscount(t *testing.T) {
	t.Parallel()
	updatedDiscountJSON := loadExampleInput(t, "discounts", "update")
	resp, err := updateDiscount(nonexistentID, updatedDiscountJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "successfully updating a product should respond 404")
}

func TestDiscountUpdateWithInvalidBody(t *testing.T) {
	t.Parallel()
	resp, err := updateDiscount(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "successfully updating a product should respond 400")
}
