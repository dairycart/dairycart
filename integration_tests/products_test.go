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

func replaceTimeStringsForProductTests(body string) string {
	re := regexp.MustCompile(`(?U)(,?)"(available_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func createProductCreationBody(sku string, upc string) string {
	upcPart := ""
	if upc != "" {
		upcPart = fmt.Sprintf(`
		"upc": "%s",`, upc)
	}
	bodyTemplate := `
		{
			"name": "New Product",
			"subtitle": "this is a product",
			"description": "this product is neat or maybe its not who really knows for sure?",
			"sku": "%s",%s
			"manufacturer": "Manufacturer",
			"brand": "Brand",
			"quantity": 123,
			"taxable": false,
			"price": 12.34,
			"on_sale": true,
			"sale_price": 10.00,
			"cost": 5,
			"product_weight": 9,
			"product_height": 9,
			"product_width": 9,
			"product_length": 9,
			"package_weight": 9,
			"package_height": 9,
			"package_width": 9,
			"package_length": 9,
			"quantity_per_package": 3
		}
	`
	return fmt.Sprintf(bodyTemplate, sku, upcPart)
}

func createProductOptionCreationBody(name string) string {
	output := fmt.Sprintf(`
		{
			"name": "%s",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`, name)
	return output
}

func TestProductExistenceRouteForExistingProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a product that exists should respond 200")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductExistenceRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body for nonexistent product should be empty")
}

func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku `+"`nonexistent`"+`) does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "trying to retrieve a product that doesn't exist should respond 404")
}

func TestProductRetrievalRoute(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProduct(existentSKU)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"sku": "t-shirt",
			"upc": "",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity": 666,
			"taxable": true,
			"price": 20,
			"on_sale": false,
			"sale_price": 0,
			"cost": 10,
			"product_weight": 1,
			"product_height": 5,
			"product_width": 5,
			"product_length": 5,
			"package_weight": 1,
			"package_height": 5,
			"package_width": 5,
			"package_length": 5,
			"quantity_per_package": 1
		}
	`)
	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
}

func TestProductListRouteWithDefaultFilter(t *testing.T) {
	resp, err := retrieveListOfProducts(nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"count": 11,
			"limit": 25,
			"page": 1,
			"data": [
				{
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"sku": "t-shirt",
					"upc": "",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"taxable": true,
					"price": 20,
					"on_sale": false,
					"sale_price": 0,
					"cost": 10,
					"product_weight": 1,
					"product_height": 5,
					"product_width": 5,
					"product_length": 5,
					"package_weight": 1,
					"package_height": 5,
					"package_width": 5,
					"package_length": 5,
					"quantity_per_package": 1
				}, {
					"name": "Sleeping People - Sleeping People",
					"subtitle": "A solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "sleeping-people",
					"upc": "656605908410",
					"manufacturer": "Record Company",
					"brand": "Sleeping People",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Jaga Jazzist - One Armed Bandit",
					"subtitle": "A solid jazz album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "one-armed-bandit",
					"upc": "5021392578187",
					"manufacturer": "Record Company",
					"brand": "Jaga Jazzist",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Cloudkicker - Let Yourself Be Huge",
					"subtitle": "A solid instrumental album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "let-yourself-be-huge",
					"upc": "",
					"manufacturer": "Record Company",
					"brand": "Cloudkicker",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Animals As Leaders - The Joy Of Motion",
					"subtitle": "A solid prog metal album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "the-joy-of-motion",
					"upc": "817424013895",
					"manufacturer": "Record Company",
					"brand": "Animals As Leaders",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Mort Garson - Mother Earth's Plantasia",
					"subtitle": "A solid synth album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "mother-earths-plantasia",
					"upc": "5291103812552",
					"manufacturer": "Record Company",
					"brand": "Mort Garson",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Camel - The Snow Goose",
					"subtitle": "A solid prog rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "the-snow-goose",
					"upc": "600753356661",
					"manufacturer": "Record Company",
					"brand": "Camel",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Piglet - Lava Land",
					"subtitle": "Another solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "lava-land",
					"upc": "",
					"manufacturer": "Record Company",
					"brand": "Piglet",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Tera Melos - Untitled",
					"subtitle": "Yet another solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "untitled",
					"upc": "634457550513",
					"manufacturer": "Record Company",
					"brand": "Tera Melos",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Frank Zappa - Jazz From Hell",
					"subtitle": "A solid Zappa album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "jazz-from-hell",
					"upc": "013347420516",
					"manufacturer": "Record Company",
					"brand": "Frank Zappa",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "CHON - Newborn Sun",
					"subtitle": "Yet another solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "newborn-sun",
					"upc": "794558090315",
					"manufacturer": "Record Company",
					"brand": "CHON",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}
			]
		}
	`)
	assert.Equal(t, expected, actual, "product list route should respond with a list of products")
}

func TestProductListRouteWithCustomFilter(t *testing.T) {
	customFilter := map[string]string{
		"page":  "2",
		"limit": "5",
	}
	resp, err := retrieveListOfProducts(customFilter)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"count": 11,
			"limit": 5,
			"page": 2,
			"data": [
				{
					"name": "Mort Garson - Mother Earth's Plantasia",
					"subtitle": "A solid synth album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "mother-earths-plantasia",
					"upc": "5291103812552",
					"manufacturer": "Record Company",
					"brand": "Mort Garson",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Camel - The Snow Goose",
					"subtitle": "A solid prog rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "the-snow-goose",
					"upc": "600753356661",
					"manufacturer": "Record Company",
					"brand": "Camel",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Piglet - Lava Land",
					"subtitle": "Another solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "lava-land",
					"upc": "",
					"manufacturer": "Record Company",
					"brand": "Piglet",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Tera Melos - Untitled",
					"subtitle": "Yet another solid math rock album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "untitled",
					"upc": "634457550513",
					"manufacturer": "Record Company",
					"brand": "Tera Melos",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}, {
					"name": "Frank Zappa - Jazz From Hell",
					"subtitle": "A solid Zappa album",
					"description": "Arbitrary description can go here because real product descriptions are technically copywritten.",
					"sku": "jazz-from-hell",
					"upc": "013347420516",
					"manufacturer": "Record Company",
					"brand": "Frank Zappa",
					"quantity": 123,
					"taxable": true,
					"price": 12.34,
					"on_sale": false,
					"sale_price": 0,
					"cost": 5,
					"product_weight": 1,
					"product_height": 12,
					"product_width": 12,
					"product_length": 0.5,
					"package_weight": 1,
					"package_height": 12,
					"package_width": 12,
					"package_length": 0.5,
					"quantity_per_package": 1
				}
			]
		}
	`)
	assert.Equal(t, expected, actual, "product list route should respond with a customized list of products")
}

func TestProductUpdateRoute(t *testing.T) {
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"sku": "t-shirt",
			"upc": "",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity": 666,
			"taxable": true,
			"price": 20,
			"on_sale": false,
			"sale_price": 0,
			"cost": 10,
			"product_weight": 1,
			"product_height": 5,
			"product_width": 5,
			"product_length": 5,
			"package_weight": 1,
			"package_height": 5,
			"package_width": 5,
			"package_length": 5,
			"quantity_per_package": 1
		}
	`)
	assert.Equal(t, expected, actual, "product response upon update should reflect the updated fields")
}

func TestProductUpdateRouteWithCompletelyInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProduct(existentSKU, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "trying to update a product with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product update route should respond with failure message when you try to update a product with invalid input")
}

func TestProductUpdateRouteWithInvalidSKU(t *testing.T) {
	t.Parallel()
	JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "trying to update a product with an invalid sku should respond 400")
}

func TestProductUpdateRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(nonexistentSKU, JSONBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "requesting a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku `+"`nonexistent`"+`) does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "trying to update a product that doesn't exist should respond 404")
}

func TestProductCreation(t *testing.T) {
	t.Parallel()

	testSKU := "test-product-creation"
	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "0123456789")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a product that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)

		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"name": "New Product",
				"subtitle": "this is a product",
				"description": "this product is neat or maybe its not who really knows for sure?",
				"sku": "%s",
				"upc": "0123456789",
				"manufacturer": "Manufacturer",
				"brand": "Brand",
				"quantity": 123,
				"taxable": false,
				"price": 12.34,
				"on_sale": true,
				"sale_price": 10,
				"cost": 5,
				"product_weight": 9,
				"product_height": 9,
				"product_width": 9,
				"product_length": 9,
				"package_weight": 9,
				"package_height": 9,
				"package_width": 9,
				"package_length": 9,
				"quantity_per_package": 3
			}
		`, testSKU))
		assert.Equal(t, expected, actual, "product creation route should respond with created product body")
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProduct(testSKU)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a product that exists should respond 200")
	}

	t.Run("create product", testProductCreation)
	t.Run("delete created product", testDeleteProduct)
}

func TestProductDeletion(t *testing.T) {
	t.Parallel()

	testSKU := "test-product-deletion"
	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a product that doesn't exist should respond 201")
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProduct(testSKU)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a product that exists should respond 200")

		actual := turnResponseBodyIntoString(t, resp)
		expected := fmt.Sprintf("Successfully deleted product `%s`", testSKU)
		assert.Equal(t, expected, actual, "product deletion route should respond with affirmative message upon successful deletion")
	}

	t.Run("create product", testProductCreation)
	t.Run("delete created product", testDeleteProduct)
}

func TestProductCreationWithAlreadyExistentSKU(t *testing.T) {
	t.Parallel()
	existentProductJSON := `
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"sku": "t-shirt",
			"upc": "",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity": 666,
			"taxable": true,
			"price": 20,
			"on_sale": false,
			"sale_price": 0,
			"cost": 10,
			"product_weight": 1,
			"product_height": 5,
			"product_width": 5,
			"product_length": 5,
			"package_weight": 1,
			"package_height": 5,
			"package_width": 5,
			"package_length": 5,
			"quantity_per_package": 1
		}
	`
	resp, err := createProduct(existentProductJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a product that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding string
	expected := minifyJSON(t, `
		{
			"status": 400,
			"message": "product with sku `+"`t-shirt`"+` already exists"
		}
	`)
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
}

func TestProductCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProduct(exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a product that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product creation route should respond with failure message when you try to create a product with invalid input")
}

func TestProductOptionListRetrievalWithDefaultFilter(t *testing.T) {
	resp, err := retrieveProductOptions("1", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"count": 2,
			"limit": 25,
			"page": 1,
			"data": [{
				"name": "color",
				"product_id": 1,
				"values": [{
					"product_option_id": 1,
					"value": "red"
				}, {
					"product_option_id": 1,
					"value": "green"
				}, {
					"product_option_id": 1,
					"value": "blue"
				}]
			}, {
				"name": "size",
				"product_id": 1,
				"values": [{
					"product_option_id": 2,
					"value": "small"
				}, {
					"product_option_id": 2,
					"value": "medium"
				}, {
					"product_option_id": 2,
					"value": "large"
				}]
			}]
		}
	`)
	assert.Equal(t, expected, actual, "product option list route should respond with a list of product options and their values")
}

func TestProductOptionListRetrievalWithCustomFilter(t *testing.T) {
	customFilter := map[string]string{
		"page":  "2",
		"limit": "1",
	}
	resp, err := retrieveProductOptions("1", customFilter)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "requesting a list of products should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"count": 2,
			"limit": 1,
			"page": 2,
			"data": [{
				"name": "size",
				"product_id": 1,
				"values": [{
					"product_option_id": 2,
					"value": "small"
				}, {
					"product_option_id": 2,
					"value": "medium"
				}, {
					"product_option_id": 2,
					"value": "large"
				}]
			}]
		}
	`)
	assert.Equal(t, expected, actual, "product option list route should respond with a list of product options and their values")
}

func TestProductOptionCreation(t *testing.T) {
	t.Parallel()

	testOptionName := "example_option_to_create"
	var createdOptionID uint64
	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a product option that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
		actual := cleanAPIResponseBody(body)

		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"name": "%s",
				"product_id": 1,
				"values": [
					{
						"product_option_id": %d,
						"value": "one"
					},{
						"product_option_id": %d,
						"value": "two"
					},{
						"product_option_id": %d,
						"value": "three"
					}
				]
			}
		`, testOptionName, createdOptionID, createdOptionID, createdOptionID))
		assert.Equal(t, expected, actual, "product option creation route should respond with created product option body")
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a product that exists should respond 200")
	}

	t.Run("create product option", testProductOptionCreation)
	t.Run("delete created product option", testDeleteProductOption)
}

func TestProductOptionDeletion(t *testing.T) {
	t.Parallel()

	testOptionName := "example_option_to_delete"
	var createdOptionID uint64
	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a product option that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a product that exists should respond 200")

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
	}

	t.Run("create product option", testProductOptionCreation)
	t.Run("delete created product option", testDeleteProductOption)
}

func TestProductOptionCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionForProduct(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "trying to create a new product option with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionCreationWithAlreadyExistentName(t *testing.T) {
	// /* TODO: */
	// t.Parallel()
	newOptionJSON := loadExampleInput(t, "product_options", "new")
	existingOptionJSON := strings.Replace(newOptionJSON, "example_value", "color", 1)
	resp, err := createProductOptionForProduct(existentID, existingOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a product option that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 400,
			"message": "product option with the name `+"`color`"+` already exists"
		}
	`)
	assert.Equal(t, expected, actual, "product option creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductOptionUpdate(t *testing.T) {
	updatedOptionJSON := loadExampleInput(t, "product_options", "update")
	resp, err := updateProductOption(existentID, updatedOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"name": "not_the_same",
			"product_id": 1,
			"values": null
		}
	`)
	assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
}

func TestProductOptionUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()
	updatedOptionJSON := loadExampleInput(t, "product_options", "update")
	resp, err := updateProductOption(nonexistentID, updatedOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product option you were looking for (id `+"`999999999`"+`) does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductOptionValueCreation(t *testing.T) {
	newOptionValueJSON := loadExampleInput(t, "product_option_values", "new")
	resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a product option value that doesn't exist should respond 201")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"product_option_id": 1,
			"value": "example_value"
		}
	`)
	assert.Equal(t, expected, actual, "product option value creation route should respond with created product option body")
}

func TestProductOptionValueCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "trying to create a new product option value with invalid input should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option value creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueCreationWithAlreadyExistentValue(t *testing.T) {
	t.Parallel()
	newOptionJSON := loadExampleInput(t, "product_option_values", "new")
	existingOptionJSON := strings.Replace(newOptionJSON, "example_value", "blue", 1)
	resp, err := createProductOptionValueForOption(existentID, existingOptionJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a product option value that already exists should respond 400")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 400,
			"message": "product option value `+"`blue`"+` already exists for option ID 1"
		}
	`)
	assert.Equal(t, expected, actual, "product option value creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductOptionValueUpdate(t *testing.T) {
	updatedOptionValueJSON := loadExampleInput(t, "product_option_values", "update")
	resp, err := updateProductOptionValueForOption(existentID, updatedOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "successfully updating a product should respond 200")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"product_option_id": 1,
			"value": "not_the_same_value"
		}
	`)
	assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
}

func TestProductOptionValueUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "successfully updating a product should respond 400")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()
	updatedOptionValueJSON := loadExampleInput(t, "product_option_values", "update")
	resp, err := updateProductOptionValueForOption(nonexistentID, updatedOptionValueJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "successfully updating a product should respond 404")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product option value you were looking for (id `+"`999999999`"+`) does not exist"
		}
	`)
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
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "updating a product option value with an already existent value should respond 500")

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedInternalErrorResponse, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := deleteProduct(nonexistentSKU)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "trying to delete a product that doesn't exist should respond 404")

	actual := turnResponseBodyIntoString(t, resp)
	// TODO: stop adding strings
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku `+"`nonexistent`"+`) does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
}
