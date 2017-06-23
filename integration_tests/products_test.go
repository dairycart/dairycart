package dairytest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestProductOptionListRetrievalWithDefaultFilter(t *testing.T) {
	resp, err := retrieveProductOptions("1", nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "list_with_default_filter"))
	assert.Equal(t, expected, actual, "product option list route should respond with a list of product options and their values")
}

func TestProductOptionListRetrievalWithCustomFilter(t *testing.T) {
	customFilter := map[string]string{
		"page":  "2",
		"limit": "1",
	}
	resp, err := retrieveProductOptions("1", customFilter)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "list_with_custom_filter"))
	assert.Equal(t, expected, actual, "product option list route should respond with a list of product options and their values")
}

func TestProductOptionCreation(t *testing.T) {
	newOptionJSON := loadExampleInput(t, "product_options", "new")
	resp, err := createProductOptionForProgenitor(existentID, newOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode, "creating a product option that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "created"))
	assert.Equal(t, expected, actual, "product option creation route should respond with created product option body")
}

func TestProductOptionCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionForProgenitor(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to create a new product option with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product option creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionCreationWithAlreadyExistentName(t *testing.T) {
	newOptionJSON := loadExampleInput(t, "product_options", "new")
	existingOptionJSON := strings.Replace(newOptionJSON, "example_value", "color", 1)
	resp, err := createProductOptionForProgenitor(existentID, existingOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product option that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "error_creating_name_already_exists"))
	assert.Equal(t, expected, actual, "product option creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductOptionUpdate(t *testing.T) {
	updatedOptionJSON := loadExampleInput(t, "product_options", "update")
	resp, err := updateProductOption(existentID, updatedOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "updated"))
	assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
}

func TestProductOptionUpdateWithInvalidInput(t *testing.T) {
	resp, err := updateProductOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()
	updatedOptionJSON := loadExampleInput(t, "product_options", "update")
	resp, err := updateProductOption(nonexistentID, updatedOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_options", "error_option_does_not_exist"))
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductOptionValueCreation(t *testing.T) {
	newOptionValueJSON := loadExampleInput(t, "product_option_values", "new")
	resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode, "creating a product option value that doesn't exist should respond 201")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "created"))
	assert.Equal(t, expected, actual, "product option value creation route should respond with created product option body")
}

func TestProductOptionValueCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to create a new product option value with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product option value creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueCreationWithAlreadyExistentValue(t *testing.T) {
	newOptionJSON := loadExampleInput(t, "product_option_values", "new")
	existingOptionJSON := strings.Replace(newOptionJSON, "example_value", "blue", 1)
	resp, err := createProductOptionValueForOption(existentID, existingOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product option value that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "error_value_already_exists"))
	assert.Equal(t, expected, actual, "product option value creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductOptionValueUpdate(t *testing.T) {
	updatedOptionValueJSON := loadExampleInput(t, "product_option_values", "update")
	resp, err := updateProductOptionValueForOption(existentID, updatedOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "updated"))
	assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
}

func TestProductOptionValueUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "error_invalid_body"))
	assert.Equal(t, expected, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()
	updatedOptionValueJSON := loadExampleInput(t, "product_option_values", "update")
	resp, err := updateProductOptionValueForOption(nonexistentID, updatedOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_option_values", "error_value_does_not_exist"))
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductOptionValueUpdateForAlreadyExistentValue(t *testing.T) {
	// Say you have a product option called `color`, and it has three values (`red`, `green`, and `blue`).
	// Let's say you try to change `red` to `blue` for whatever reason. That will fail at the database level,
	// because the schema ensures a unique combination of value and option ID. Should I prevent users from
	// being able to do this? On the one hand, it adds yet another query to a route that should presumably never
	// experience that issue at all. On the other hand it does provide a convenient and clear explanation
	// for why a given problem occurred.
	duplicatedOptionValueJSON := loadExampleInput(t, "product_option_values", "duplicate")
	resp, err := updateProductOptionValueForOption("4", duplicatedOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode, "updating a product option value with an already existent value should respond 500")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	assert.Equal(t, expectedInternalErrorResponse, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
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
