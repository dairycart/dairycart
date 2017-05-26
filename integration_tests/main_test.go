package dairytest

import (
	"fmt"
	"io/ioutil"
	"log"
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
	// we can't reliably predict what the `updated_at` or `archived_at` columns could possibly equal, so we strip them out of the body becuase we're bad at programming.
	timeFieldReplacementPatterns = `,"(created_at|updated_at|archived_at)":({"Time":)?"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z"(,"Valid":(true|false))?(})?`
	applicationJSON              = "application/json"

	existentSKU            = "skateboard"
	nonexistentSKU         = "nonexistent"
	expected404SKUResponse = "The product you were looking for (sku `nonexistent`) does not exist"
)

func init() {
	err := ensureThatDairycartIsAlive()
	if err != nil {
		log.Fatalf("dairycart isn't up: %v", err)
	}
	jsonMinifier = minify.New()
	jsonMinifier.AddFunc(applicationJSON, jsonMinify.Minify)
}

func readTestFile(t *testing.T, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("test_files/%s.json", filename))
	assert.Nil(t, err)
	return strings.TrimSpace(string(bodyBytes))
}

func replaceTimeStringsForTests(body string) string {
	re := regexp.MustCompile(timeFieldReplacementPatterns)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func turnResponseBodyIntoString(res *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	return strings.TrimSpace(string(bodyBytes)), err
}

func minifyExampleJSON(t *testing.T, json string) string {
	minified, err := jsonMinifier.String(applicationJSON, json)
	assert.Nil(t, err)
	return minified
}

func TestProductExistenceRouteForExistingProduct(t *testing.T) {
	resp, err := checkProductExistence(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product that exists should respond 200")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)
	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductExistenceRouteForNonexistentProduct(t *testing.T) {
	resp, err := checkProductExistence(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)
	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
	resp, err := retrieveProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Equal(t, expected404SKUResponse, actual, "trying to retrieve a product that doesn't exist should respond 404")
}

func TestProductRetrievalRoute(t *testing.T) {
	resp, err := retrieveProduct(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product should respond 200")

	body, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	skateboardProductJSON := readTestFile(t, existentSKU)
	actual := replaceTimeStringsForTests(body)
	expected := minifyExampleJSON(t, skateboardProductJSON)
	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
}

func TestProductListRoute(t *testing.T) {
	resp, err := retrieveListOfProducts()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")

	respBody, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)
	actual := replaceTimeStringsForTests(respBody)

	expected := minifyExampleJSON(t, readTestFile(t, "product_list_response"))
	assert.Equal(t, expected, actual, "product list route should respond with a list of products")
}

func TestProductUpdateRoute(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	body, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	actual := replaceTimeStringsForTests(body)
	minified := minifyExampleJSON(t, readTestFile(t, existentSKU))
	expected := strings.Replace(minified, `"quantity":123`, `"quantity":666`, 1)
	assert.Equal(t, expected, actual, "product response upon update should reflect the updated fields")
}

func TestProductUpdateRouteWithCompletelyInvalidInput(t *testing.T) {
	JSONBody := `{"testing": true}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with invalid input should respond 400")
}

func TestProductUpdateRouteWithInvalidSKU(t *testing.T) {
	JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with an invalid sku should respond 400")
}

func TestProductUpdateRouteForNonexistentProduct(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(nonexistentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Equal(t, expected404SKUResponse, actual, "trying to update a product that doesn't exist should respond 404")
}

func TestProductCreation(t *testing.T) {
	newProductJSON := readTestFile(t, "example_new_product")
	resp, err := createProduct(newProductJSON)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "creating a product that doesn't exist should respond 200")

	respBody, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)
	actual := replaceTimeStringsForTests(respBody)

	expected := minifyExampleJSON(t, readTestFile(t, "created_product_response"))
	assert.Equal(t, expected, actual, "product creation route should respond with created product body")
}

func TestProductCreationWithAlreadyExistentSKU(t *testing.T) {
	newProductJSON := readTestFile(t, "example_new_product")
	bodyToUse := strings.Replace(newProductJSON, `"sku": "new-product"`, `"sku": "skateboard"`, 1)
	resp, err := createProduct(bodyToUse)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product that already exists should respond 400")

	actual, err := turnResponseBodyIntoString(resp)
	expected := "product with sku `skateboard` already exists"
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
}

func TestProductDeletionRouteForNewlyCreatedProduct(t *testing.T) {
	resp, err := deleteProduct("new-product")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "trying to delete a product that exists should respond 200")

	actual, err := turnResponseBodyIntoString(resp)
	expected := "Successfully deleted product `new-product`"
	assert.Equal(t, expected, actual, "product deletion route should respond with affirmative message upon successful deletion")
}

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	resp, err := deleteProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "trying to delete a product that doesn't exist should respond 404")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Equal(t, expected404SKUResponse, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
}
