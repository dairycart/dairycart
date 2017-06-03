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
	//
	// Note that this pattern needs to be run as ungreedy because of the possiblity of prefix and or
	// suffixed commas
	timeFieldReplacementPattern = `(?U)(,?)"(created_at|updated_at|archived_at)":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z"(,?)`
	applicationJSON             = "application/json"

	existentProgenitorID   = "1"
	existentSKU            = "skateboard"
	nonexistentSKU         = "nonexistent"
	exampleGarbageInput    = `{"testing_garbage_input": true}`
	expected404SKUResponse = "The product you were looking for (sku `nonexistent`) does not exist"
)

func init() {
	ensureThatDairycartIsAlive()
	jsonMinifier = minify.New()
	jsonMinifier.AddFunc(applicationJSON, jsonMinify.Minify)
}

func loadExpectedResponse(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("expected_responses/%s/%s.json", folder, filename))
	assert.Nil(t, err)
	return strings.TrimSpace(string(bodyBytes))
}

func loadExampleInput(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("example_inputs/%s/%s.json", folder, filename))
	assert.Nil(t, err)
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
	minified, err := jsonMinifier.String(applicationJSON, json)
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
	assert.Equal(t, expected404SKUResponse, actual, "trying to retrieve a product that doesn't exist should respond 404")
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
	expected := "Invalid input provided for product body"
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
	assert.Equal(t, expected404SKUResponse, actual, "trying to update a product that doesn't exist should respond 404")
}

func TestProductCreation(t *testing.T) {
	newProductJSON := loadExampleInput(t, "products", "new")
	resp, err := createProduct(newProductJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "creating a product that doesn't exist should respond 200")

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
	expected := "product with sku `skateboard` already exists"
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
}

func TestProductCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProduct(exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := "Invalid input provided for product body"
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a product with invalid input")
}

func TestProductAttributeListRetrievalWithDefaultFilter(t *testing.T) {
	resp, err := retrieveProductAttributes("1", nil)
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
	resp, err := retrieveProductAttributes("1", customFilter)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(body)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "list_with_custom_filter"))
	assert.Equal(t, expected, actual, "product attribute list route should respond with a list of product attributes and their values")
}

func TestProductAttributeCreation(t *testing.T) {
	newAttributeJSON := loadExampleInput(t, "product_attributes", "new")
	resp, err := createProductAttributeForProgenitor(existentProgenitorID, newAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "creating a product attribute that doesn't exist should respond 200")

	respBody := turnResponseBodyIntoString(t, resp)
	actual := replaceTimeStringsForTests(respBody)
	expected := minifyJSON(t, loadExpectedResponse(t, "product_attributes", "created"))
	assert.Equal(t, expected, actual, "product attribute creation route should respond with created product attribute body")
}

func TestProductAttributeCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductAttributeForProgenitor(existentProgenitorID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to create a new product attribute with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := "Invalid input provided for product attribute body"
	assert.Equal(t, expected, actual, "product attribute creation route should respond with failure message when you provide it invalid input")
}

func TestProductAttributeCreationWithAlreadyExistentName(t *testing.T) {
	newAttributeJSON := loadExampleInput(t, "product_attributes", "new")
	existingAttributeJSON := strings.Replace(newAttributeJSON, "example_value", "color", 1)
	resp, err := createProductAttributeForProgenitor(existentProgenitorID, existingAttributeJSON)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product attribute that doesn't exist should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	expected := "product attribute with the name `color` already exists"
	assert.Equal(t, expected, actual, "product attribute creation route should respond with failure message when you try to create a value that already exists")
}

// I'd like to keep these functions last if at all possible.

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := deleteProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "trying to delete a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expected404SKUResponse, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductDeletionRouteForNewlyCreatedProduct(t *testing.T) {
	resp, err := deleteProduct("new-product")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "trying to delete a product that exists should respond 200")

	actual := turnResponseBodyIntoString(t, resp)
	expected := "Successfully deleted product `new-product`"
	assert.Equal(t, expected, actual, "product deletion route should respond with affirmative message upon successful deletion")
}
