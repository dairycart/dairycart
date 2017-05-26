package dairytest

import (
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
	skateboardProductJSON        = `
		{
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": false,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"id": 10,
			"product_progenitor_id": 2,
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"on_sale": false,
			"price": 12.34,
			"sale_price": ""
		}
	`
	newProductJSON = `
		{
			"description": "This is a new product.",
			"taxable": true,
			"product_weight": 9,
			"product_height": 9,
			"product_width": 9,
			"product_length": 9,
			"package_weight": 9,
			"package_height": 9,
			"package_width": 9,
			"package_length": 9,
			"sku": "new-product",
			"name": "New Product",
			"upc": "0123456789",
			"quantity": 123,
			"on_sale": false,
			"price": 12.34,
			"sale_price": null
		}
	`
)

func init() {
	err := ensureThatDairycartIsAlive()
	if err != nil {
		log.Fatalf("dairycart isn't up: %v", err)
	}
	jsonMinifier = minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
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
	minified, err := jsonMinifier.String("application/json", json)
	assert.Nil(t, err)
	return minified
}

func TestProductExistenceRouteForExistingProduct(t *testing.T) {
	resp, err := checkProductExistence("skateboard")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product that exists should respond 200")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductExistenceRouteForNonexistentProduct(t *testing.T) {
	resp, err := checkProductExistence("nonexistent")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
	resp, err := retrieveProduct("nonexistent")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")
}

func TestProductRetrievalRoute(t *testing.T) {
	resp, err := retrieveProduct("skateboard")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a product should respond 200")

	body, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	actual := replaceTimeStringsForTests(body)
	expected := minifyExampleJSON(t, skateboardProductJSON)
	assert.Equal(t, expected, actual, "product response should contain a complete product")
}

func TestProductListRoute(t *testing.T) {
	resp, err := retrieveListOfProducts()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "requesting a list of products should respond 200")
}

func TestProductUpdateRoute(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct("skateboard", JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "successfully updating a product should respond 200")

	actual, err := turnResponseBodyIntoString(resp)
	assert.Nil(t, err)

	expected := `"Product updated"`
	assert.Equal(t, expected, actual, "product response should reflect the updated fields")
}

func TestProductUpdateRouteWithCompletelyInvalidInput(t *testing.T) {
	JSONBody := `{"testing": true}`
	resp, err := updateProduct("skateboard", JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with invalid input should respond 400")
}

func TestProductUpdateRouteWithInvalidSKU(t *testing.T) {
	JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
	resp, err := updateProduct("skateboard", JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "trying to update a product with invalid input should respond 400")
}

func TestProductUpdateRouteForNonexistentProduct(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct("nonexistent", JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "requesting a product that doesn't exist should respond 404")
}

// // FIXME: broken as shit
// func TestProductCreation(t *testing.T) {
// 	resp, err := createProduct(newProductJSON)
// 	assert.Nil(t, err)

// 	body, err := turnResponseBodyIntoString(resp)
// 	assert.Nil(t, err)

// 	log.Printf(`

// 		received body:
// 		 	%v
// 	`, body)

// 	assert.Equal(t, 200, resp.StatusCode, "creating a product that doesn't exist should respond 200")
// }

func TestProductCreationWithAlreadyExistentSKU(t *testing.T) {
	bodyToUse := strings.Replace(newProductJSON, `"sku": "new-product"`, `"sku": "skateboard"`, 1)
	resp, err := createProduct(bodyToUse)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode, "creating a product that already exist should respond 400")
}

// // uncomment after implementing (and testing) product creation
// func TestProductDeletionRouteForNewlyCreatedProduct(t *testing.T) {
// 	resp, err := deleteProduct("new_product")
// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, resp.StatusCode, "trying to delete a product that exists should respond 200")
// }

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	resp, err := deleteProduct("nonexistent")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "trying to delete a product that doesn't exist should respond 404")
}
