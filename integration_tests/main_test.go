package dairytest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify"

	jsonMinify "github.com/tdewolff/minify/json"
)

var jsonMinifier *minify.M

const (
	// we can't reliably predict what the `updated_at` or `archived_at` columns could possibly equal,
	// so we strip them out of the body becuase we're bad at programming. The (sort of) plus side to
	// this is that we ensure our timestamps have a particular format (because if they didn't, this
	// function, and as a consequence, the tests, would fail spectacularly).
	// Note that this pattern needs to be run as ungreedy because of the possiblity of prefix and or
	// suffixed commas
	timeFieldReplacementPattern = `(?U)(,?)"(created_at|updated_at|archived_at)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`

	existentID     = "1"
	nonexistentID  = "999999999"
	existentSKU    = "skateboard"
	nonexistentSKU = "nonexistent"

	exampleGarbageInput           = `{"testing_garbage_input": true}`
	expectedInternalErrorResponse = `{"status":500,"message":"Unexpected internal error occurred"}`
)

func init() {
	ensureThatDairycartIsAlive()
	jsonMinifier = minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
}

func loadExpectedResponse(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("expected_responses/%s/%s.json", folder, filename))
	assert.Nil(t, err)
	assert.NotEmpty(t, bodyBytes, "example response file requested is empty and should not be")
	return strings.TrimSpace(string(bodyBytes))
}

func loadExampleInput(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("example_inputs/%s/%s.json", folder, filename))
	assert.Nil(t, err)
	assert.NotEmpty(t, bodyBytes, "example input file requested is empty and should not be")
	return strings.TrimSpace(string(bodyBytes))
}

func replaceTimeStringsForTests(body string) string {
	re := regexp.MustCompile(timeFieldReplacementPattern)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func turnResponseBodyIntoString(t *testing.T, res *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	return strings.TrimSpace(string(bodyBytes))
}

func minifyJSON(t *testing.T, json string) string {
	minified, err := jsonMinifier.String("application/json", json)
	assert.Nil(t, err)
	return minified
}

func TestProductExistenceRouteForExistingProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product that exists should respond 200")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductExistenceRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body for nonexistent product should be empty")
}

func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_product_does_not_exist"))
	assert.Equal(t, expected, actual, "trying to retrieve a product that doesn't exist should respond 404")
}

func TestProductRetrievalRoute(t *testing.T) {
	resp, err := retrieveProduct(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	skateboardProductJSON := loadExpectedResponse(t, "products", "skateboard")
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, skateboardProductJSON)
	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
}

func TestProductListRouteWithDefaultFilter(t *testing.T) {
	resp, err := retrieveListOfProducts(nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "list_with_default_filter"))
	assert.Equal(t, expected, actual, "product list route should respond with a list of products")
}

func TestProductListRouteWithCustomFilter(t *testing.T) {
	customFilter := map[string]string{
		"page":  "2",
		"limit": "5",
	}
	resp, err := retrieveListOfProducts(customFilter)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "list_with_custom_filter"))
	assert.Equal(t, expected, actual, "product list route should respond with a customized list of products")
}

func TestProductUpdateRoute(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	minified := minifyJSON(t, loadExpectedResponse(t, "products", "skateboard"))
	expected := strings.Replace(minified, `"quantity":123`, `"quantity":666`, 1)
	assert.Equal(t, expected, actual, "product response upon update should reflect the updated fields")
}

func TestProductUpdateRouteWithCompletelyInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProduct(existentSKU, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product update route should respond with failure message when you try to update a product with invalid input")
}

func TestProductUpdateRouteWithInvalidSKU(t *testing.T) {
	t.Parallel()
	JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with an invalid sku should respond 400")
}

func TestProductUpdateRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(nonexistentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_product_does_not_exist"))
	assert.Equal(t, expected, actual, "trying to update a product that doesn't exist should respond 404")
}

func TestProductCreation(t *testing.T) {
	newProductJSON := loadExampleInput(t, "products", "new")
	resp, err := createProduct(newProductJSON)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode, "creating a product that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "created"))
	assert.Equal(t, expected, actual, "product creation route should respond with created product body")
}

func TestProductCreationWithAlreadyExistentSKU(t *testing.T) {
	t.Parallel()
	existentProductJSON := loadExpectedResponse(t, "products", "skateboard")
	resp, err := createProduct(existentProductJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_sku_already_exists"))
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
}

func TestProductCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProduct(exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a product with invalid input")
}

func TestProductAttributeListRetrievalWithDefaultFilter(t *testing.T) {
	resp, err := retrieveProductAttributes(nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "list_with_default_filter"))
	assert.Equal(t, expected, actual, "product attribute list route should respond with a list of product attributes and their values")
}

func TestProductAttributeListRetrievalWithCustomFilter(t *testing.T) {
	customFilter := map[string]string{
		"page":  "2",
		"limit": "1",
	}
	resp, err := retrieveProductAttributes(customFilter)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "list_with_custom_filter"))
	assert.Equal(t, expected, actual, "product attribute list route should respond with a list of product attributes and their values")
}

func TestProductAttributeCreation(t *testing.T) {
	newAttributeJSON := loadExampleInput(t, "product_attributes", "new")
	resp, err := createProductAttribute(newAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode, "creating a product attribute that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "created"))
	assert.Equal(t, expected, actual, "product attribute creation route should respond with created product attribute body")
}

func TestProductAttributeCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductAttribute(exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to create a new product attribute with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product attribute creation route should respond with failure message when you provide it invalid input")
}

func TestProductAttributeCreationWithAlreadyExistentName(t *testing.T) {
	newAttributeJSON := loadExampleInput(t, "product_attributes", "new")
	existingAttributeJSON := strings.Replace(newAttributeJSON, "example_value", "Material", 1)
	resp, err := createProductAttribute(existingAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product attribute that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "error_creating_name_already_exists"))
	assert.Equal(t, expected, actual, "product attribute creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductAttributeUpdate(t *testing.T) {
	updatedAttributeJSON := loadExampleInput(t, "product_attributes", "update")
	resp, err := updateProductAttribute(existentID, updatedAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "updated"))
	assert.Equal(t, expected, actual, "product attribute update response should reflect the updated fields")
}

func TestProductAttributeUpdateWithInvalidInput(t *testing.T) {
	resp, err := updateProductAttribute(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product attribute update route should respond with failure message when you provide it invalid input")
}

func TestProductAttributeUpdateForNonexistentAttribute(t *testing.T) {
	t.Parallel()
	updatedAttributeJSON := loadExampleInput(t, "product_attributes", "update")
	resp, err := updateProductAttribute(nonexistentID, updatedAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "error_attribute_does_not_exist"))
	assert.Equal(t, expected, actual, "product attribute update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductAttributeValueCreation(t *testing.T) {
	newAttributeValueJSON := loadExampleInput(t, "product_attribute_values", "new")
	resp, err := createProductAttributeValueForAttribute(existentID, newAttributeValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode, "creating a product attribute value that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "created"))
	assert.Equal(t, expected, actual, "product attribute value creation route should respond with created product attribute body")
}

func TestProductAttributeValueCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductAttributeValueForAttribute(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to create a new product attribute value with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product attribute value creation route should respond with failure message when you provide it invalid input")
}

func TestProductAttributeValueCreationWithAlreadyExistentValue(t *testing.T) {
	newAttributeJSON := loadExampleInput(t, "product_attribute_values", "new")
	existingAttributeJSON := strings.Replace(newAttributeJSON, "example_value", "Cotton", 1)
	resp, err := createProductAttributeValueForAttribute(existentID, existingAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product attribute value that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "error_value_already_exists"))
	assert.Equal(t, expected, actual, "product attribute value creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductAttributeValueUpdate(t *testing.T) {
	updatedAttributeValueJSON := loadExampleInput(t, "product_attribute_values", "update")
	resp, err := updateProductAttributeValueForAttribute(existentID, updatedAttributeValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "updated"))
	assert.Equal(t, expected, actual, "product attribute update response should reflect the updated fields")
}

func TestProductAttributeValueUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductAttributeValueForAttribute(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product attribute update route should respond with failure message when you provide it invalid input")
}

func TestProductAttributeValueUpdateForNonexistentAttribute(t *testing.T) {
	t.Parallel()
	updatedAttributeValueJSON := loadExampleInput(t, "product_attribute_values", "update")
	resp, err := updateProductAttributeValueForAttribute(nonexistentID, updatedAttributeValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attribute_values", "error_value_does_not_exist"))
	assert.Equal(t, expected, actual, "product attribute update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductAttributeValueUpdateForAlreadyExistentValue(t *testing.T) {
	// Say you have a product attribute called `color`, and it has three values (`red`, `green`, and `blue`).
	// Let's say you try to change `red` to `blue` for whatever reason. That will fail at the database level,
	// because the schema ensures a unique combination of value and attribute ID. Should I prevent users from
	// being able to do this? On the one hand, it adds yet another query to a route that should presumably never
	// experience that issue at all. On the other hand it does provide a convenient and clear explanation
	// for why a given problem occurred.
	duplicatedAttributeValueJSON := loadExampleInput(t, "product_attribute_values", "duplicate")
	resp, err := updateProductAttributeValueForAttribute("4", duplicatedAttributeValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode, "updating a product attribute value with an already existent value should respond 500")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	assert.Equal(t, expectedInternalErrorResponse, actual, "product attribute update route should respond with 404 message when you try to delete a product that doesn't exist")
}

////////////////////////////////////////////////////////
//                                                    //
//                Discount Route Tests                //
//                                                    //
////////////////////////////////////////////////////////

func replaceTimeStringsForDiscountTests(body string) string {
	re := regexp.MustCompile(`(?U)(,?)"(starts_on|expires_on)":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z"(,?)`)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func TestDiscountRetrievalForExistingDiscount(t *testing.T) {
	resp, err := getDiscountByID(existentID)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "a successfully retrieved discount should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	removedTimeFields := replaceTimeStringsForTests(body)
	actual := replaceTimeStringsForDiscountTests(removedTimeFields)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "retrieved"))
	assert.Equal(t, expected, actual, "discount route should return a serialized discount object")
}

func TestDiscountRetrievalForNonexistentDiscount(t *testing.T) {
	resp, err := getDiscountByID(nonexistentID)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "a request for a nonexistent discount should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "discounts", "error_discount_does_not_exist"))
	assert.Equal(t, expected, actual, "product attribute update route should respond with 404 message when you try to delete a product that doesn't exist")

}

// I'd like to keep these functions last if at all possible.

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := deleteProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "trying to delete a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "products", "error_product_does_not_exist"))
	assert.Equal(t, expected, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductDeletionRouteForNewlyCreatedProduct(t *testing.T) {
	resp, err := deleteProduct("new-product")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "trying to delete a product that exists should respond 200")

	actual := turnResponseBodyIntoString(t, resp)
	expected := "Successfully deleted product `new-product`"
	assert.Equal(t, expected, actual, "product deletion route should respond with affirmative message upon successful deletion")
}
