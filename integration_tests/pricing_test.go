package dairytest

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func replaceTimeStringsForDiscountTests(body string) string {
	re := regexp.MustCompile(`(?U)(,?)"(starts_on|expires_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func TestDiscountRetrievalForExistingDiscount(t *testing.T) {
	// /* TODO: */
	// t.Parallel()
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
	newDiscountJSON := loadExampleInput(t, "discounts", "new")
	resp, err := createDiscount(newDiscountJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a discount that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "created"))
	assert.Equal(t, expected, actual, "discount creation route should respond with created product body")
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
