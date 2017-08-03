package dairytest

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"

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
			"quantity_per_package": 3,
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
			"package_length": 9
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

func createProductOptionBody(name string) string {
	output := fmt.Sprintf(`
		{
			"name": "%s"
		}
	`, name)
	return output
}

func createProductOptionValueBody(value string) string {
	output := fmt.Sprintf(`
		{
			"value": "%s"
		}
	`, value)
	return output
}

func TestProductExistenceRouteForExistingProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(existentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body should be empty")
}

func TestProductExistenceRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := checkProductExistence(nonexistentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, "", actual, "product existence body for nonexistent product should be empty")
}

func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProduct(nonexistentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku 'nonexistent') does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "trying to retrieve a product that doesn't exist should respond 404")
}

func TestProductRetrievalRoute(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProduct(existentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, fmt.Sprintf(`
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"option_summary": "Size: Small, Color: Red",
			"sku": "%s",
			"upc": "",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity": 666,
			"quantity_per_package": 1,
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
			"package_length": 5
		}
	`, existentSKU))
	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
}

func TestProductListRouteWithDefaultFilter(t *testing.T) {
	t.Parallel()
	resp, err := retrieveListOfProducts(nil)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)
	lr := parseResponseIntoStruct(t, body)
	assert.True(t, len(lr.Data) <= int(lr.Limit), "product list route should not return more data than the limit")
	assert.Equal(t, uint8(25), lr.Limit, "product list route should respond with the default limit when a ilmit is not specified")
	assert.Equal(t, uint64(1), lr.Page, "product list route should respond with the first page when a page is not specified")
}

func TestProductListRouteWithCustomFilter(t *testing.T) {
	t.Parallel()
	customFilter := map[string]string{
		"page":  "2",
		"limit": "5",
	}
	resp, err := retrieveListOfProducts(customFilter)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)
	lr := parseResponseIntoStruct(t, body)
	assert.Equal(t, uint8(5), lr.Limit, "product list route should respond with the specified limit")
	assert.Equal(t, uint64(2), lr.Page, "product list route should respond with the specified page")
}

func TestProductUpdateRoute(t *testing.T) {
	t.Parallel()
	testSKU := "test-product-updating"
	var productRootID uint64

	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)
		body := turnResponseBodyIntoString(t, resp)
		productRootID = retrieveIDFromResponseBody(body, t)
	}

	testUpdateProduct := func(t *testing.T) {
		JSONBody := `{"quantity":666}`
		resp, err := updateProduct(testSKU, JSONBody)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		body := turnResponseBodyIntoString(t, resp)
		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, `
			{
				"name": "New Product",
				"subtitle": "this is a product",
				"description": "this product is neat or maybe its not who really knows for sure?",
				"option_summary": "",
				"sku": "test-product-updating",
				"upc": "",
				"manufacturer": "Manufacturer",
				"brand": "Brand",
				"quantity": 666,
				"quantity_per_package": 3,
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
				"package_length": 9
			}
		`)
		assert.Equal(t, expected, actual, "product response upon update should reflect the updated fields")
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProduct(testSKU)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product",
			Test:    testProductCreation,
		},
		{
			Message: "update product",
			Test:    testUpdateProduct,
		},
		{
			Message: "delete created product",
			Test:    testDeleteProduct,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductUpdateRouteWithCompletelyInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProduct(existentSKU, exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product update route should respond with failure message when you try to update a product with invalid input")
}

func TestProductUpdateRouteWithInvalidSKU(t *testing.T) {
	t.Parallel()
	JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
	resp, err := updateProduct(existentSKU, JSONBody)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)
}

func TestProductUpdateRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	JSONBody := `{"quantity":666}`
	resp, err := updateProduct(nonexistentSKU, JSONBody)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku 'nonexistent') does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "trying to update a product that doesn't exist should respond 404")
}

func TestProductCreationRoute(t *testing.T) {
	t.Parallel()

	var productRootID uint64
	testSKU := "test-product-creation"
	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "0123456789")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		productRootID = retrieveIDFromResponseBody(body, t)

		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"name": "New Product",
				"subtitle": "this is a product",
				"description": "this product is neat or maybe its not who really knows for sure?",
				"sku_prefix": "%s",
				"manufacturer": "Manufacturer",
				"brand": "Brand",
				"quantity_per_package": 3,
				"taxable": false,
				"cost": 5,
				"product_weight": 9,
				"product_height": 9,
				"product_width": 9,
				"product_length": 9,
				"package_weight": 9,
				"package_height": 9,
				"package_width": 9,
				"package_length": 9,
				"options": [],
				"products": [{
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "",
					"sku": "%s",
					"upc": "0123456789",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"package_length": 9
				}]
			}
		`, testSKU, testSKU))
		assert.Equal(t, expected, actual, "product creation route should respond with created product body")
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product",
			Test:    testProductCreation,
		},
		{
			Message: "delete created product",
			Test:    testDeleteProduct,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductCreationRouteWithOptions(t *testing.T) {
	t.Parallel()

	var productRootID uint64
	testSKU := "test-product-creation-with-options"
	testProductCreation := func(t *testing.T) {
		newProductJSON := fmt.Sprintf(`
			{
				"name": "New Product",
				"subtitle": "this is a product",
				"description": "this product is neat or maybe its not who really knows for sure?",
				"sku": "%s",
				"manufacturer": "Manufacturer",
				"brand": "Brand",
				"quantity": 123,
				"quantity_per_package": 3,
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
				"options": [
					{
						"name": "color",
						"values": [
							"red",
							"green",
							"blue"
						]
					},
					{
						"name": "size",
						"values": [
							"small",
							"medium",
							"large"
						]
					}
				]
			}
		`, testSKU)
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		productRootID = retrieveIDFromResponseBody(body, t)
		actual := cleanAPIResponseBody(body)

		tmpl, err := template.New("resp").Parse(`
			{
				"name": "New Product",
				"subtitle": "this is a product",
				"description": "this product is neat or maybe its not who really knows for sure?",
				"sku_prefix": "{{.SKU}}",
				"manufacturer": "Manufacturer",
				"brand": "Brand",
				"quantity_per_package": 3,
				"taxable": false,
				"cost": 5,
				"product_weight": 9,
				"product_height": 9,
				"product_width": 9,
				"product_length": 9,
				"package_weight": 9,
				"package_height": 9,
				"package_width": 9,
				"package_length": 9,
				"options": [{
					"name": "color",
					"values": [{
						"value": "red"
					}, {
						"value": "green"
					}, {
						"value": "blue"
					}]
				}, {
					"name": "size",
					"values": [{
						"value": "small"
					}, {
						"value": "medium"
					}, {
						"value": "large"
					}]
				}],
				"products": [{
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: red, size: small",
					"sku": "{{.SKU}}_red_small",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "red"
					}, {
						"value": "small"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: red, size: medium",
					"sku": "{{.SKU}}_red_medium",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "red"
					}, {
						"value": "medium"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: red, size: large",
					"sku": "{{.SKU}}_red_large",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "red"
					}, {
						"value": "large"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: green, size: small",
					"sku": "{{.SKU}}_green_small",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "green"
					}, {
						"value": "small"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: green, size: medium",
					"sku": "{{.SKU}}_green_medium",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "green"
					}, {
						"value": "medium"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: green, size: large",
					"sku": "{{.SKU}}_green_large",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "green"
					}, {
						"value": "large"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: blue, size: small",
					"sku": "{{.SKU}}_blue_small",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "blue"
					}, {
						"value": "small"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: blue, size: medium",
					"sku": "{{.SKU}}_blue_medium",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "blue"
					}, {
						"value": "medium"
					}]
				}, {
					"name": "New Product",
					"subtitle": "this is a product",
					"description": "this product is neat or maybe its not who really knows for sure?",
					"option_summary": "color: blue, size: large",
					"sku": "{{.SKU}}_blue_large",
					"upc": "",
					"manufacturer": "Manufacturer",
					"brand": "Brand",
					"quantity": 123,
					"quantity_per_package": 3,
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
					"applicable_options": [{
						"value": "blue"
					}, {
						"value": "large"
					}]
				}]
			}
		`)
		assert.Nil(t, err)

		b := new(bytes.Buffer)
		x := struct {
			SKU string
		}{
			SKU: testSKU,
		}
		err = tmpl.Execute(b, x)
		assert.Nil(t, err)

		expected := minifyJSON(t, b.String())
		assert.Equal(t, expected, actual, "product creation route should respond with created product body")
	}

	testDeleteProductRoot := func(t *testing.T) {
		resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product",
			Test:    testProductCreation,
		},
		{
			Message: "delete created product",
			Test:    testDeleteProductRoot,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductDeletion(t *testing.T) {
	t.Parallel()

	testSKU := "test-product-deletion"
	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProduct(testSKU)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Empty(t, actual)
	}

	subtests := []subtest{
		{
			Message: "create product",
			Test:    testProductCreation,
		},
		{
			Message: "delete created product",
			Test:    testDeleteProduct,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
	t.Parallel()
	resp, err := deleteProduct(nonexistentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product you were looking for (sku 'nonexistent') does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductRootListRetrievalRouteWithDefaultPagination(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProductRoots(nil)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)
	body := turnResponseBodyIntoString(t, resp)

	lr := parseResponseIntoStruct(t, body)
	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
	assert.Equal(t, uint8(25), lr.Limit, "product option list route should respond with the default limit when a ilmit is not specified")
	assert.Equal(t, uint64(1), lr.Page, "product option list route should respond with the first page when a page is not specified")
}

func TestProductRootListRetrievalRouteWithCustomPagination(t *testing.T) {
	t.Parallel()
	customFilter := map[string]string{
		"page":  "2",
		"limit": "1",
	}
	resp, err := retrieveProductRoots(customFilter)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)
	body := turnResponseBodyIntoString(t, resp)

	lr := parseResponseIntoStruct(t, body)
	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
	assert.Equal(t, uint8(1), lr.Limit, "product option list route should respond with the default limit when a ilmit is not specified")
	assert.Equal(t, uint64(2), lr.Page, "product option list route should respond with the first page when a page is not specified")
}

func TestProductRootRetrievalRoute(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProductRoot(existentID)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"sku_prefix": "t-shirt",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity_per_package": 1,
			"taxable": true,
			"cost": 20,
			"product_weight": 1,
			"product_height": 5,
			"product_width": 5,
			"product_length": 5,
			"package_weight": 1,
			"package_height": 5,
			"package_width": 5,
			"package_length": 5,
			"options": [{
				"name": "color",
				"values": [{
					"value": "red"
				}, {
					"value": "green"
				}, {
					"value": "blue"
				}]
			}, {
				"name": "size",
				"values": [{
					"value": "small"
				}, {
					"value": "medium"
				}, {
					"value": "large"
				}]
			}],
			"products": [{
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Small, Color: Red",
				"sku": "t-shirt-small-red",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Medium, Color: Red",
				"sku": "t-shirt-medium-red",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Large, Color: Red",
				"sku": "t-shirt-large-red",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Small, Color: Blue",
				"sku": "t-shirt-small-blue",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Medium, Color: Blue",
				"sku": "t-shirt-medium-blue",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Large, Color: Blue",
				"sku": "t-shirt-large-blue",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Small, Color: Green",
				"sku": "t-shirt-small-green",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Medium, Color: Green",
				"sku": "t-shirt-medium-green",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}, {
				"name": "Your Favorite Band's T-Shirt",
				"subtitle": "A t-shirt you can wear",
				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				"option_summary": "Size: Large, Color: Green",
				"sku": "t-shirt-large-green",
				"upc": "",
				"manufacturer": "Record Company",
				"brand": "Your Favorite Band",
				"quantity": 666,
				"quantity_per_package": 1,
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
				"package_length": 5
			}]
		}
	`)
	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
}

func TestProductRootRetrievalRouteForNonexistentRoot(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProductRoot(nonexistentID)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := `{"status":404,"message":"The product_root you were looking for (identified by '999999999') does not exist"}`
	assert.Equal(t, expected, actual, "expected and actual bodies should match")
}

func TestProductRootDeletionRoute(t *testing.T) {
	t.Parallel()

	var productRootID uint64
	testSKU := "test-product-root-deletion"
	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)

		body := turnResponseBodyIntoString(t, resp)
		productRootID = retrieveIDFromResponseBody(body, t)
		assertStatusCode(t, resp, http.StatusCreated)
	}

	testDeleteProductRoot := func(t *testing.T) {
		resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
		assert.Nil(t, err)

		body := turnResponseBodyIntoString(t, resp)

		assert.Empty(t, body)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product",
			Test:    testProductCreation,
		},
		{
			Message: "delete created product root",
			Test:    testDeleteProductRoot,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductRootDeletionRouteForNonexistentRoot(t *testing.T) {
	t.Parallel()
	resp, err := deleteProductRoot(nonexistentID)
	assert.Nil(t, err)

	actual := turnResponseBodyIntoString(t, resp)
	expected := `{"status":404,"message":"The product_root you were looking for (identified by '999999999') does not exist"}`
	assert.Equal(t, expected, actual, "expected and actual bodies should match")
	assertStatusCode(t, resp, http.StatusNotFound)
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
			"quantity_per_package": 1,
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
			"package_length": 5
		}
	`
	resp, err := createProduct(existentProductJSON)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, `
		{
			"status": 400,
			"message": "product with sku 't-shirt' already exists"
		}
	`)
	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
}

func TestProductCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProduct(exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product creation route should respond with failure message when you try to create a product with invalid input")
}

func TestProductOptionListRetrievalWithDefaultFilter(t *testing.T) {
	t.Parallel()
	resp, err := retrieveProductOptions("1", nil)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)

	lr := parseResponseIntoStruct(t, body)
	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
	assert.Equal(t, uint8(25), lr.Limit, "product option list route should respond with the default limit when a ilmit is not specified")
	assert.Equal(t, uint64(1), lr.Page, "product option list route should respond with the first page when a page is not specified")
}

func TestProductOptionListRetrievalWithCustomFilter(t *testing.T) {
	t.Parallel()
	customFilter := map[string]string{
		"page":  "2",
		"limit": "1",
	}
	resp, err := retrieveProductOptions("1", customFilter)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	body := turnResponseBodyIntoString(t, resp)

	lr := parseResponseIntoStruct(t, body)
	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
	assert.Equal(t, uint8(1), lr.Limit, "product option list route should respond with the default limit when a ilmit is not specified")
	assert.Equal(t, uint64(2), lr.Page, "product option list route should respond with the first page when a page is not specified")
}

func TestProductOptionCreation(t *testing.T) {
	t.Parallel()

	testOptionName := "example_option_to_create"
	var createdOptionID uint64
	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
		actual := cleanAPIResponseBody(body)

		expected := minifyJSON(t, `
			{
				"name": "example_option_to_create",
				"values": [
					{
						"value": "one"
					},
					{
						"value": "two"
					},
					{
						"value": "three"
					}
				]
			}
		`)
		assert.Equal(t, expected, actual, "product option creation route should respond with created product option body")
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product option",
			Test:    testProductOptionCreation,
		},
		{
			Message: "delete created product option",
			Test:    testDeleteProductOption,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionDeletion(t *testing.T) {
	t.Parallel()

	testOptionName := "example_option_to_delete"
	var createdOptionID uint64
	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
	}

	subtests := []subtest{
		{
			Message: "create product option",
			Test:    testProductOptionCreation,
		},
		{
			Message: "delete created product option",
			Test:    testDeleteProductOption,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionDeletionForNonexistentOption(t *testing.T) {
	t.Parallel()

	resp, err := deleteProductOption(nonexistentID)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := `{"status":404,"message":"The product option you were looking for (id '999999999') does not exist"}`
	assert.Equal(t, expected, actual, "product option deletion route should respond with affirmative message upon successful deletion")
}

func TestProductOptionCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionForProduct(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionCreationWithAlreadyExistentName(t *testing.T) {
	t.Parallel()
	testOptionName := "already-existent-option"
	var createdOptionID uint64
	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
	}

	testDuplicateProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusBadRequest)
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
	}

	subtests := []subtest{
		{
			Message: "create product option",
			Test:    testProductOptionCreation,
		},
		{
			Message: "create product option again",
			Test:    testDuplicateProductOptionCreation,
		},
		{
			Message: "delete created product option",
			Test:    testDeleteProductOption,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionUpdate(t *testing.T) {
	t.Parallel()
	testSKU := "testing_product_options"
	testOptionName := "example_option_to_update"
	var createdOptionID uint64
	var createdRootID uint64

	updatedOptionName := "not_the_same"

	testProductCreation := func(t *testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)
		body := turnResponseBodyIntoString(t, resp)
		createdRootID = retrieveIDFromResponseBody(body, t)
	}

	testProductOptionCreation := func(t *testing.T) {
		newOptionJSON := createProductOptionCreationBody(testOptionName)
		resp, err := createProductOptionForProduct(strconv.Itoa(int(createdRootID)), newOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionID = retrieveIDFromResponseBody(body, t)
	}

	testUpdateProductOption := func(t *testing.T) {
		updatedOptionJSON := createProductOptionBody(updatedOptionName)
		resp, err := updateProductOption(strconv.Itoa(int(createdOptionID)), updatedOptionJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		body := turnResponseBodyIntoString(t, resp)
		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, `
			{
				"name": "not_the_same",
				"values": [{
					"value": "one"
				}, {
					"value": "two"
				}, {
					"value": "three"
				}]
			}
		`)
		assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
	}

	testDeleteProductOption := func(t *testing.T) {
		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
	}

	testDeleteProduct := func(t *testing.T) {
		resp, err := deleteProductRoot(strconv.Itoa(int(createdRootID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product to add option to",
			Test:    testProductCreation,
		},
		{
			Message: "create product option",
			Test:    testProductOptionCreation,
		},
		{
			Message: "update product option",
			Test:    testUpdateProductOption,
		},
		{
			Message: "delete created product option",
			Test:    testDeleteProductOption,
		},
		{
			Message: "delete created product root",
			Test:    testDeleteProduct,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()
	updatedOptionName := "nonexistent-not-the-same"
	updatedOptionJSON := createProductOptionBody(updatedOptionName)
	resp, err := updateProductOption(nonexistentID, updatedOptionJSON)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product option you were looking for (id '999999999') does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductOptionValueCreation(t *testing.T) {
	t.Parallel()

	var createdOptionValueID uint64
	testValue := "test-value-creation"
	testCreateProductOptionValue := func(t *testing.T) {
		newOptionValueJSON := createProductOptionValueBody(testValue)
		resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionValueID = retrieveIDFromResponseBody(body, t)
		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"value": "%s"
			}
		`, testValue))
		assert.Equal(t, expected, actual, "product option value creation route should respond with created product option body")
	}

	testDeleteProductOptionValue := func(t *testing.T) {
		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product option value",
			Test:    testCreateProductOptionValue,
		},
		{
			Message: "delete created product option value",
			Test:    testDeleteProductOptionValue,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionValueUpdate(t *testing.T) {
	t.Parallel()

	var createdOptionValueID uint64
	testValue := "test-value-update"
	testUpdatedValue := "not_the_same_value"
	testCreateProductOptionValue := func(t *testing.T) {
		newOptionValueJSON := createProductOptionValueBody(testValue)
		resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		body := turnResponseBodyIntoString(t, resp)
		createdOptionValueID = retrieveIDFromResponseBody(body, t)
	}

	testUpdateProductOptionValue := func(t *testing.T) {
		updatedOptionValueJSON := createProductOptionValueBody(testUpdatedValue)
		resp, err := updateProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)), updatedOptionValueJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		body := turnResponseBodyIntoString(t, resp)
		actual := cleanAPIResponseBody(body)
		expected := minifyJSON(t, fmt.Sprintf(`
			{
				"value": "%s"
			}
		`, testUpdatedValue))
		assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
	}

	testDeleteProductOptionValue := func(t *testing.T) {
		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "create product option value",
			Test:    testCreateProductOptionValue,
		},
		{
			Message: "update created product option value",
			Test:    testUpdateProductOptionValue,
		},
		{
			Message: "delete created product option value",
			Test:    testDeleteProductOptionValue,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionValueDeletion(t *testing.T) {
	t.Parallel()

	var createdOptionValueID uint64
	testValue := "test-value-deletion"
	testCreateProductOptionValue := func(t *testing.T) {
		newOptionValueJSON := createProductOptionValueBody(testValue)
		resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)
		body := turnResponseBodyIntoString(t, resp)
		createdOptionValueID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteProductOptionValue := func(t *testing.T) {
		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
	}

	subtests := []subtest{
		{
			Message: "create product option value",
			Test:    testCreateProductOptionValue,
		},
		{
			Message: "delete created product option value",
			Test:    testDeleteProductOptionValue,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestProductOptionValueDeletionForNonexistentOptionValue(t *testing.T) {
	t.Parallel()

	resp, err := deleteProductOptionValueForOption(nonexistentID)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	actual := turnResponseBodyIntoString(t, resp)
	expected := `{"status":404,"message":"The product option value you were looking for (id '999999999') does not exist"}`
	assert.Equal(t, expected, actual, "product option deletion route should respond with affirmative message upon successful deletion")
}

func TestProductOptionValueCreationWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := createProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option value creation route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueCreationWithAlreadyExistentValue(t *testing.T) {
	t.Parallel()

	alreadyExistentValue := "blue"
	existingOptionJSON := createProductOptionValueBody(alreadyExistentValue)
	resp, err := createProductOptionValueForOption(existentID, existingOptionJSON)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	actual := turnResponseBodyIntoString(t, resp)
	expected := minifyJSON(t, `
		{
			"status": 400,
			"message": "product option value 'blue' already exists for option ID 1"
		}
	`)
	assert.Equal(t, expected, actual, "product option value creation route should respond with failure message when you try to create a value that already exists")
}

func TestProductOptionValueUpdateWithInvalidInput(t *testing.T) {
	t.Parallel()
	resp, err := updateProductOptionValueForOption(existentID, exampleGarbageInput)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusBadRequest)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
}

func TestProductOptionValueUpdateForNonexistentOption(t *testing.T) {
	t.Parallel()

	obligatoryValue := "whatever"
	updatedOptionValueJSON := createProductOptionValueBody(obligatoryValue)
	resp, err := updateProductOptionValueForOption(nonexistentID, updatedOptionValueJSON)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusNotFound)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	expected := minifyJSON(t, `
		{
			"status": 404,
			"message": "The product option value you were looking for (id '999999999') does not exist"
		}
	`)
	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}

func TestProductOptionValueUpdateForAlreadyExistentValue(t *testing.T) {
	// Say you have a product option called `color`, and it has three values (`red`, `green`, and `blue`).
	// Let's say you try to change `red` to `blue` for whatever reason. That will fail at the database level,
	// because the schema ensures a unique combination of value and option ID and archived date.
	t.Parallel()

	duplicatedOptionValueJSON := createProductOptionValueBody("medium")
	resp, err := updateProductOptionValueForOption("4", duplicatedOptionValueJSON)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusInternalServerError)

	body := turnResponseBodyIntoString(t, resp)
	actual := cleanAPIResponseBody(body)
	assert.Equal(t, expectedInternalErrorResponse, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
}
