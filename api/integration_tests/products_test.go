package dairytest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	// "text/template"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func replaceTimeStringsForProductTests(body string) string {
	re := regexp.MustCompile(`(?U)(,?)"(available_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func logBodyAndResetResponse(t *testing.T, resp *http.Response) {
	t.Helper()
	respStr := turnResponseBodyIntoString(t, resp)
	log.Printf(`

		%s

	`, respStr)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(respStr)))
}

func createTestProduct(t *testing.T, sku string, upc string) uint64 {
	t.Helper()
	newProductJSON := createProductCreationBody(sku, upc)
	resp, err := createProduct(newProductJSON)
	require.Nil(t, err)
	assertStatusCode(t, resp, http.StatusCreated)

	var p models.Product
	unmarshalBody(resp, &p)
	return p.ID
}

func deleteTestProduct(t *testing.T, sku string) {
	resp, err := deleteProduct(sku)
	require.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)
}

func deleteTestProductRoot(t *testing.T, productRootID uint64) {
	resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
	require.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)
}

func compareListResponses(t *testing.T, expected models.ListResponse, actual models.ListResponse) {
	t.Helper()
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual limit don't match")
	assert.Equal(t, expected.Page, actual.Page, "expected and actual page don't match")
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

func TestProductExistenceRoute(t *testing.T) {
	// t.Parallel()

	t.Run("for existing product", func(*testing.T) {
		resp, err := checkProductExistence(existentSKU)
		require.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product existence body should be empty")
	})

	t.Run("for nonexistent product", func(*testing.T) {
		resp, err := checkProductExistence(nonexistentSKU)
		require.Nil(t, err)
		assertStatusCode(t, resp, http.StatusNotFound)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, "", actual, "product existence body for nonexistent product should be empty")
	})
}

// TODO: maybe these functions should just set the values that we don't care about equality for rather than check for the equality of each field
// for instance, we don't really worry about IDs, so make this function set the expected.ID to actual.ID and then use assert to check equality
func compareProducts(t *testing.T, expected models.Product, actual models.Product) {
	t.Helper()
	// assert.Equal(t, expected.ProductRootID, actual.ProductRootID, "expected and actual ProductRootID values don't match")
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name values don't match")
	assert.Equal(t, expected.Subtitle, actual.Subtitle, "expected and actual Subtitle values don't match")
	assert.Equal(t, expected.Description, actual.Description, "expected and actual Description values don't match")
	assert.Equal(t, expected.OptionSummary, actual.OptionSummary, "expected and actual OptionSummary values don't match")
	assert.Equal(t, expected.SKU, actual.SKU, "expected and actual SKU values don't match")
	assert.Equal(t, expected.UPC, actual.UPC, "expected and actual UPC values don't match")
	assert.Equal(t, expected.Manufacturer, actual.Manufacturer, "expected and actual Manufacturer values don't match")
	assert.Equal(t, expected.Brand, actual.Brand, "expected and actual Brand values don't match")
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected and actual Quantity values don't match")
	assert.Equal(t, expected.Taxable, actual.Taxable, "expected and actual Taxable values don't match")
	assert.Equal(t, expected.Price, actual.Price, "expected and actual Price values don't match")
	assert.Equal(t, expected.OnSale, actual.OnSale, "expected and actual OnSale values don't match")
	assert.Equal(t, expected.SalePrice, actual.SalePrice, "expected and actual SalePrice values don't match")
	assert.Equal(t, expected.Cost, actual.Cost, "expected and actual Cost values don't match")
	assert.Equal(t, expected.ProductWeight, actual.ProductWeight, "expected and actual ProductWeight values don't match")
	assert.Equal(t, expected.ProductHeight, actual.ProductHeight, "expected and actual ProductHeight values don't match")
	assert.Equal(t, expected.ProductWidth, actual.ProductWidth, "expected and actual ProductWidth values don't match")
	assert.Equal(t, expected.ProductLength, actual.ProductLength, "expected and actual ProductLength values don't match")
	assert.Equal(t, expected.PackageWeight, actual.PackageWeight, "expected and actual PackageWeight values don't match")
	assert.Equal(t, expected.PackageHeight, actual.PackageHeight, "expected and actual PackageHeight values don't match")
	assert.Equal(t, expected.PackageWidth, actual.PackageWidth, "expected and actual PackageWidth values don't match")
	assert.Equal(t, expected.PackageLength, actual.PackageLength, "expected and actual PackageLength values don't match")
	assert.Equal(t, expected.QuantityPerPackage, actual.QuantityPerPackage, "expected and actual QuantityPerPackage values don't match")
	assert.Equal(t, expected.ApplicableOptionValues, actual.ApplicableOptionValues, "expected and actual ApplicableOptionValues values don't match")
}

func compareProductRoots(t *testing.T, expected, actual models.ProductRoot) {
	t.Helper()
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name should match")
	assert.Equal(t, expected.Subtitle, actual.Subtitle, "expected and actual Subtitle should match")
	assert.Equal(t, expected.Description, actual.Description, "expected and actual Description should match")
	assert.Equal(t, expected.SKUPrefix, actual.SKUPrefix, "expected and actual SKUPrefix should match")
	assert.Equal(t, expected.Manufacturer, actual.Manufacturer, "expected and actual Manufacturer should match")
	assert.Equal(t, expected.Brand, actual.Brand, "expected and actual Brand should match")
	assert.Equal(t, expected.Taxable, actual.Taxable, "expected and actual Taxable should match")
	assert.Equal(t, expected.Cost, actual.Cost, "expected and actual Cost should match")
	assert.Equal(t, expected.ProductWeight, actual.ProductWeight, "expected and actual ProductWeight should match")
	assert.Equal(t, expected.ProductHeight, actual.ProductHeight, "expected and actual ProductHeight should match")
	assert.Equal(t, expected.ProductWidth, actual.ProductWidth, "expected and actual ProductWidth should match")
	assert.Equal(t, expected.ProductLength, actual.ProductLength, "expected and actual ProductLength should match")
	assert.Equal(t, expected.PackageWeight, actual.PackageWeight, "expected and actual PackageWeight should match")
	assert.Equal(t, expected.PackageHeight, actual.PackageHeight, "expected and actual PackageHeight should match")
	assert.Equal(t, expected.PackageWidth, actual.PackageWidth, "expected and actual PackageWidth should match")
	assert.Equal(t, expected.PackageLength, actual.PackageLength, "expected and actual PackageLength should match")
	assert.Equal(t, expected.QuantityPerPackage, actual.QuantityPerPackage, "expected and actual QuantityPerPackage should match")

	for i := range expected.Options {
		assert.Equal(t, expected.Options[i], actual.Options[i], "expected and actual option #%d should match", i)
	}

	for i := range expected.Products {
		compareProducts(t, expected.Products[i], actual.Products[i])
	}
}

func TestProductRetrievalRoute(t *testing.T) {
	// // t.Parallel()
	resp, err := retrieveProduct(existentSKU)
	assert.Nil(t, err)
	assertStatusCode(t, resp, http.StatusOK)

	expected := models.Product{
		ProductRootID:          1,
		Name:                   "Your Favorite Band's T-Shirt",
		Subtitle:               "A t-shirt you can wear",
		Description:            "Wear this if you'd like. Or don't, I'm not in charge of your actions",
		OptionSummary:          "Size: Small, Color: Red",
		SKU:                    existentSKU,
		Manufacturer:           "Record Company",
		Brand:                  "Your Favorite Band",
		Quantity:               666,
		Taxable:                true,
		Price:                  20,
		OnSale:                 false,
		SalePrice:              0,
		Cost:                   10,
		ProductWeight:          1,
		ProductHeight:          5,
		ProductWidth:           5,
		ProductLength:          5,
		PackageWeight:          1,
		PackageHeight:          5,
		PackageWidth:           5,
		PackageLength:          5,
		QuantityPerPackage:     1,
		ApplicableOptionValues: nil,
	}

	var actual models.Product
	unmarshalBody(resp, &actual)
	compareProducts(t, expected, actual)
}

// func TestProductRetrievalRouteForNonexistentProduct(t *testing.T) {
// 	t.Skip()
// 	resp, err := retrieveProduct(nonexistentSKU)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 404,
// 			"message": "The product you were looking for (sku 'nonexistent') does not exist"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "trying to retrieve a product that doesn't exist should respond 404")
// }

func TestProductListRouteFilters(t *testing.T) {
	// // t.Parallel()

	t.Run("with standard filter", func(*testing.T) {
		resp, err := retrieveListOfProducts(nil)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		expected := models.ListResponse{
			Limit: 25,
			Page:  1,
		}
		var actual models.ListResponse
		unmarshalBody(resp, &actual)
		compareListResponses(t, expected, actual)
	})

	// FIXME
	// t.Run("with nonstandard filter", func(*testing.T) {
	// 	customFilter := map[string]string{
	// 		"page":  "2",
	// 		"limit": "5",
	// 	}
	// 	resp, err := retrieveListOfProducts(customFilter)
	// 	assert.Nil(t, err)
	// 	assertStatusCode(t, resp, http.StatusOK)

	// 	expected := models.ListResponse{
	// 		Limit: 5,
	// 		Page:  2,
	// 	}
	// 	var actual models.ListResponse
	// 	unmarshalBody(resp, &actual)
	// 	compareListResponses(t, expected, actual)
	// })
}

func TestProductUpdateRoute(t *testing.T) {
	// // t.Parallel()
	testSKU := "test-product-updating"

	t.Run("normal use", func(*testing.T) {
		productRootID := createTestProduct(t, testSKU, "")
		JSONBody := `{"quantity":666}`
		resp, err := updateProduct(testSKU, JSONBody)
		require.Nil(t, err)
		assertStatusCode(t, resp, http.StatusOK)

		expected := models.Product{
			ProductRootID:      productRootID,
			Name:               "New Product",
			Subtitle:           "this is a product",
			Description:        "this product is neat or maybe its not who really knows for sure?",
			SKU:                testSKU,
			Manufacturer:       "Manufacturer",
			Brand:              "Brand",
			Quantity:           666,
			QuantityPerPackage: 3,
			Taxable:            false,
			Price:              12.34,
			OnSale:             true,
			SalePrice:          10,
			Cost:               5,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
		}
		var actual models.Product
		unmarshalBody(resp, &actual)
		compareProducts(t, expected, actual)
		deleteTestProduct(t, testSKU)
	})

	t.Run("with completely invalid input", func(*testing.T) {
		resp, err := updateProduct(existentSKU, exampleGarbageInput)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusBadRequest)

		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, expectedBadRequestResponse, actual, "product update route should respond with failure message when you try to update a product with invalid input")
	})

	t.Run("with invalid sku", func(*testing.T) {
		JSONBody := `{"sku": "thí% $kü ïs not åny gõôd"}`
		resp, err := updateProduct(existentSKU, JSONBody)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusBadRequest)
	})

	t.Run("for nonexistent product", func(*testing.T) {
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
	})
}

func TestProductCreationRoute(t *testing.T) {
	testSKU := "test-product-creation"

	t.Run("normal usage", func(*testing.T) {
		newProductJSON := createProductCreationBody(testSKU, "0123456789")
		resp, err := createProduct(newProductJSON)
		assert.Nil(t, err)
		assertStatusCode(t, resp, http.StatusCreated)

		var actual models.ProductRoot
		unmarshalBody(resp, &actual)

		expected := models.ProductRoot{
			Name:               "New Product",
			Subtitle:           "this is a product",
			Description:        "this product is neat or maybe its not who really knows for sure?",
			SKUPrefix:          testSKU,
			Manufacturer:       "Manufacturer",
			Brand:              "Brand",
			Taxable:            false,
			Cost:               5,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
			QuantityPerPackage: 3,
			Options:            []models.ProductOption{},
			Products: []models.Product{
				{
					Name:               "New Product",
					Subtitle:           "this is a product",
					Description:        "this product is neat or maybe its not who really knows for sure?",
					OptionSummary:      "",
					SKU:                testSKU,
					UPC:                "0123456789",
					Manufacturer:       "Manufacturer",
					Brand:              "Brand",
					Quantity:           123,
					Taxable:            false,
					Price:              12.34,
					OnSale:             true,
					SalePrice:          10,
					Cost:               5,
					ProductWeight:      9,
					ProductHeight:      9,
					ProductWidth:       9,
					ProductLength:      9,
					PackageWeight:      9,
					PackageHeight:      9,
					PackageWidth:       9,
					PackageLength:      9,
					QuantityPerPackage: 3,
				},
			},
		}
		compareProductRoots(t, expected, actual)
		deleteTestProductRoot(t, actual.ID)
	})
}

// func TestProductCreationRouteWithOptions(t *testing.T) {
// 	t.Skip()

// 	var productRootID uint64
// 	testSKU := "test-product-creation-with-options"
// 	testProductCreation := func(t *testing.T) {
// 		newProductJSON := fmt.Sprintf(`
// 			{
// 				"name": "New Product",
// 				"subtitle": "this is a product",
// 				"description": "this product is neat or maybe its not who really knows for sure?",
// 				"sku": "%s",
// 				"manufacturer": "Manufacturer",
// 				"brand": "Brand",
// 				"quantity": 123,
// 				"quantity_per_package": 3,
// 				"taxable": false,
// 				"price": 12.34,
// 				"on_sale": true,
// 				"sale_price": 10.00,
// 				"cost": 5,
// 				"product_weight": 9,
// 				"product_height": 9,
// 				"product_width": 9,
// 				"product_length": 9,
// 				"package_weight": 9,
// 				"package_height": 9,
// 				"package_width": 9,
// 				"package_length": 9,
// 				"options": [
// 					{
// 						"name": "color",
// 						"values": [
// 							"red",
// 							"green",
// 							"blue"
// 						]
// 					},
// 					{
// 						"name": "size",
// 						"values": [
// 							"small",
// 							"medium",
// 							"large"
// 						]
// 					}
// 				]
// 			}
// 		`, testSKU)
// 		resp, err := createProduct(newProductJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		productRootID = retrieveIDFromResponseBody(t, body)
// 		actual := cleanAPIResponseBody(body)

// 		tmpl, err := template.New("resp").Parse(`
// 			{
// 				"name": "New Product",
// 				"subtitle": "this is a product",
// 				"description": "this product is neat or maybe its not who really knows for sure?",
// 				"sku_prefix": "{{.SKU}}",
// 				"manufacturer": "Manufacturer",
// 				"brand": "Brand",
// 				"quantity_per_package": 3,
// 				"taxable": false,
// 				"cost": 5,
// 				"product_weight": 9,
// 				"product_height": 9,
// 				"product_width": 9,
// 				"product_length": 9,
// 				"package_weight": 9,
// 				"package_height": 9,
// 				"package_width": 9,
// 				"package_length": 9,
// 				"options": [{
// 					"name": "color",
// 					"values": [{
// 						"value": "red"
// 					}, {
// 						"value": "green"
// 					}, {
// 						"value": "blue"
// 					}]
// 				}, {
// 					"name": "size",
// 					"values": [{
// 						"value": "small"
// 					}, {
// 						"value": "medium"
// 					}, {
// 						"value": "large"
// 					}]
// 				}],
// 				"products": [{
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: red, size: small",
// 					"sku": "{{.SKU}}_red_small",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "red"
// 					}, {
// 						"value": "small"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: red, size: medium",
// 					"sku": "{{.SKU}}_red_medium",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "red"
// 					}, {
// 						"value": "medium"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: red, size: large",
// 					"sku": "{{.SKU}}_red_large",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "red"
// 					}, {
// 						"value": "large"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: green, size: small",
// 					"sku": "{{.SKU}}_green_small",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "green"
// 					}, {
// 						"value": "small"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: green, size: medium",
// 					"sku": "{{.SKU}}_green_medium",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "green"
// 					}, {
// 						"value": "medium"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: green, size: large",
// 					"sku": "{{.SKU}}_green_large",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "green"
// 					}, {
// 						"value": "large"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: blue, size: small",
// 					"sku": "{{.SKU}}_blue_small",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "blue"
// 					}, {
// 						"value": "small"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: blue, size: medium",
// 					"sku": "{{.SKU}}_blue_medium",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "blue"
// 					}, {
// 						"value": "medium"
// 					}]
// 				}, {
// 					"name": "New Product",
// 					"subtitle": "this is a product",
// 					"description": "this product is neat or maybe its not who really knows for sure?",
// 					"option_summary": "color: blue, size: large",
// 					"sku": "{{.SKU}}_blue_large",
// 					"upc": "",
// 					"manufacturer": "Manufacturer",
// 					"brand": "Brand",
// 					"quantity": 123,
// 					"quantity_per_package": 3,
// 					"taxable": false,
// 					"price": 12.34,
// 					"on_sale": true,
// 					"sale_price": 10,
// 					"cost": 5,
// 					"product_weight": 9,
// 					"product_height": 9,
// 					"product_width": 9,
// 					"product_length": 9,
// 					"package_weight": 9,
// 					"package_height": 9,
// 					"package_width": 9,
// 					"package_length": 9,
// 					"applicable_options": [{
// 						"value": "blue"
// 					}, {
// 						"value": "large"
// 					}]
// 				}]
// 			}
// 		`)
// 		assert.Nil(t, err)

// 		b := new(bytes.Buffer)
// 		x := struct {
// 			SKU string
// 		}{
// 			SKU: testSKU,
// 		}
// 		err = tmpl.Execute(b, x)
// 		assert.Nil(t, err)

// 		expected := minifyJSON(t, b.String())
// 		assert.Equal(t, expected, actual, "product creation route should respond with created product body")
// 	}

// 	testDeleteProductRoot := func(t *testing.T) {
// 		resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product",
// 			Test:    testProductCreation,
// 		},
// 		{
// 			Message: "delete created product",
// 			Test:    testDeleteProductRoot,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductDeletion(t *testing.T) {
// 	t.Skip()

// 	testSKU := "test-product-deletion"
// 	testProductCreation := func(t *testing.T) {
// 		newProductJSON := createProductCreationBody(testSKU, "")
// 		resp, err := createProduct(newProductJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)
// 	}

// 	testDeleteProduct := func(t *testing.T) {
// 		resp, err := deleteProduct(testSKU)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		actual := turnResponseBodyIntoString(t, resp)
// 		assert.Empty(t, actual)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product",
// 			Test:    testProductCreation,
// 		},
// 		{
// 			Message: "delete created product",
// 			Test:    testDeleteProduct,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductDeletionRouteForNonexistentProduct(t *testing.T) {
// 	t.Skip()
// 	resp, err := deleteProduct(nonexistentSKU)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 404,
// 			"message": "The product you were looking for (sku 'nonexistent') does not exist"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product deletion route should respond with 404 message when you try to delete a product that doesn't exist")
// }

// func TestProductRootListRetrievalRouteWithDefaultPagination(t *testing.T) {
// 	t.Skip()
// 	resp, err := retrieveProductRoots(nil)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusOK)
// 	body := turnResponseBodyIntoString(t, resp)

// 	lr := parseResponseIntoStruct(t, body)
// 	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
// 	assert.Equal(t, uint8(25), lr.Limit, "product option list route should respond with the default limit when a limit is not specified")
// 	assert.Equal(t, uint64(1), lr.Page, "product option list route should respond with the first page when a page is not specified")
// }

// func TestProductRootListRetrievalRouteWithCustomPagination(t *testing.T) {
// 	t.Skip()
// 	customFilter := map[string]string{
// 		"page":  "2",
// 		"limit": "1",
// 	}
// 	resp, err := retrieveProductRoots(customFilter)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusOK)
// 	body := turnResponseBodyIntoString(t, resp)

// 	lr := parseResponseIntoStruct(t, body)
// 	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
// 	assert.Equal(t, uint8(1), lr.Limit, "product option list route should respond with the default limit when a limit is not specified")
// 	assert.Equal(t, uint64(2), lr.Page, "product option list route should respond with the first page when a page is not specified")
// }

// func TestProductRootRetrievalRoute(t *testing.T) {
// 	t.Skip()
// 	resp, err := retrieveProductRoot(existentID)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusOK)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	expected := minifyJSON(t, `{
// 			"name": "Your Favorite Band's T-Shirt",
// 			"subtitle": "A t-shirt you can wear",
// 			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 			"sku_prefix": "t-shirt",
// 			"manufacturer": "Record Company",
// 			"brand": "Your Favorite Band",
// 			"taxable": true,
// 			"cost": 20,
// 			"product_weight": 1,
// 			"product_height": 5,
// 			"product_width": 5,
// 			"product_length": 5,
// 			"package_weight": 1,
// 			"package_height": 5,
// 			"package_width": 5,
// 			"package_length": 5,
// 			"quantity_per_package": 1,
// 			"options": [{
// 				"name": "color",
// 				"values": null
// 			}, {
// 				"name": "size",
// 				"values": null
// 			}],
// 			"products": [{
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not incharge of your actions",
// 				"option_summary": "Size: Small, Color: Red",
// 				"sku": "t-shirt-small-red",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Medium, Color: Red",
// 				"sku": "t-shirt-medium-red",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Large, Color: Red",
// 				"sku": "t-shirt-large-red",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirtyou can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge ofyour actions",
// 				"option_summary": "Size: Small, Color: Blue",
// 				"sku": "t-shirt-small-blue",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Medium, Color: Blue",
// 				"sku": "t-shirt-medium-blue",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Large, Color: Blue",
// 				"sku": "t-shirt-large-blue",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt youcan wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Small, Color: Green",
// 				"sku": "t-shirt-small-green",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt youcan wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Medium, Color: Green",
// 				"sku": "t-shirt-medium-green",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}, {
// 				"name": "Your Favorite Band's T-Shirt",
// 				"subtitle": "A t-shirt you can wear",
// 				"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 				"option_summary": "Size: Large, Color: Green",
// 				"sku": "t-shirt-large-green",
// 				"upc": "",
// 				"manufacturer": "Record Company",
// 				"brand": "Your Favorite Band",
// 				"quantity": 666,
// 				"taxable": true,
// 				"price": 20,
// 				"on_sale": false,
// 				"sale_price": 0,
// 				"cost": 10,
// 				"product_weight": 1,
// 				"product_height": 5,
// 				"product_width": 5,
// 				"product_length": 5,
// 				"package_weight": 1,
// 				"package_height": 5,
// 				"package_width": 5,
// 				"package_length": 5,
// 				"quantity_per_package": 1
// 			}]
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product retrieval response should contain a complete product")
// }

// func TestProductRootRetrievalRouteForNonexistentRoot(t *testing.T) {
// 	t.Skip()
// 	resp, err := retrieveProductRoot(nonexistentID)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := `{"status":404,"message":"The product_root you were looking for (identified by '999999999') does not exist"}`
// 	assert.Equal(t, expected, actual, "expected and actual bodies should match")
// }

// func TestProductRootDeletionRoute(t *testing.T) {
// 	t.Skip()

// 	var productRootID uint64
// 	testSKU := "test-product-root-deletion"
// 	testProductCreation := func(t *testing.T) {
// 		newProductJSON := createProductCreationBody(testSKU, "")
// 		resp, err := createProduct(newProductJSON)
// 		assert.Nil(t, err)

// 		body := turnResponseBodyIntoString(t, resp)
// 		productRootID = retrieveIDFromResponseBody(t, body)
// 		assertStatusCode(t, resp, http.StatusCreated)
// 	}

// 	testDeleteProductRoot := func(t *testing.T) {
// 		resp, err := deleteProductRoot(strconv.Itoa(int(productRootID)))
// 		assert.Nil(t, err)

// 		body := turnResponseBodyIntoString(t, resp)

// 		assert.Empty(t, body)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product",
// 			Test:    testProductCreation,
// 		},
// 		{
// 			Message: "delete created product root",
// 			Test:    testDeleteProductRoot,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductRootDeletionRouteForNonexistentRoot(t *testing.T) {
// 	t.Skip()
// 	resp, err := deleteProductRoot(nonexistentID)
// 	assert.Nil(t, err)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := `{"status":404,"message":"The product_root you were looking for (identified by '999999999') does not exist"}`
// 	assert.Equal(t, expected, actual, "expected and actual bodies should match")
// 	assertStatusCode(t, resp, http.StatusNotFound)
// }

// func TestProductCreationWithAlreadyExistentSKU(t *testing.T) {
// 	t.Skip()
// 	existentProductJSON := `
// 		{
// 			"name": "Your Favorite Band's T-Shirt",
// 			"subtitle": "A t-shirt you can wear",
// 			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
// 			"sku": "t-shirt",
// 			"upc": "",
// 			"manufacturer": "Record Company",
// 			"brand": "Your Favorite Band",
// 			"quantity": 666,
// 			"quantity_per_package": 1,
// 			"taxable": true,
// 			"price": 20,
// 			"on_sale": false,
// 			"sale_price": 0,
// 			"cost": 10,
// 			"product_weight": 1,
// 			"product_height": 5,
// 			"product_width": 5,
// 			"product_length": 5,
// 			"package_weight": 1,
// 			"package_height": 5,
// 			"package_width": 5,
// 			"package_length": 5
// 		}
// 	`
// 	resp, err := createProduct(existentProductJSON)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 400,
// 			"message": "product with sku 't-shirt' already exists"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product creation route should respond with failure message when you try to create a sku that already exists")
// }

// func TestProductCreationWithInvalidInput(t *testing.T) {
// 	t.Skip()
// 	resp, err := createProduct(exampleGarbageInput)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	assert.Equal(t, expectedBadRequestResponse, actual, "product creation route should respond with failure message when you try to create a product with invalid input")
// }

// func TestProductOptionListRetrievalWithDefaultFilter(t *testing.T) {
// 	t.Skip()
// 	resp, err := retrieveProductOptions("1", nil)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusOK)

// 	body := turnResponseBodyIntoString(t, resp)

// 	lr := parseResponseIntoStruct(t, body)
// 	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
// 	assert.Equal(t, uint8(25), lr.Limit, "product option list route should respond with the default limit when a limit is not specified")
// 	assert.Equal(t, uint64(1), lr.Page, "product option list route should respond with the first page when a page is not specified")
// }

// func TestProductOptionListRetrievalWithCustomFilter(t *testing.T) {
// 	t.Skip()
// 	customFilter := map[string]string{
// 		"page":  "2",
// 		"limit": "1",
// 	}
// 	resp, err := retrieveProductOptions("1", customFilter)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusOK)

// 	body := turnResponseBodyIntoString(t, resp)

// 	lr := parseResponseIntoStruct(t, body)
// 	assert.True(t, len(lr.Data) <= int(lr.Limit), "product option list route should not return more data than the limit")
// 	assert.Equal(t, uint8(1), lr.Limit, "product option list route should respond with the default limit when a limit is not specified")
// 	assert.Equal(t, uint64(2), lr.Page, "product option list route should respond with the first page when a page is not specified")
// }

// func TestProductOptionCreation(t *testing.T) {
// 	t.Skip()

// 	testOptionName := "example_option_to_create"
// 	var createdOptionID uint64
// 	testProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionID = retrieveIDFromResponseBody(t, body)
// 		actual := cleanAPIResponseBody(body)

// 		expected := minifyJSON(t, `
// 			{
// 				"name": "example_option_to_create",
// 				"values": [
// 					{
// 						"value": "one"
// 					},
// 					{
// 						"value": "two"
// 					},
// 					{
// 						"value": "three"
// 					}
// 				]
// 			}
// 		`)
// 		assert.Equal(t, expected, actual, "product option creation route should respond with created product option body")
// 	}

// 	testDeleteProductOption := func(t *testing.T) {
// 		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product option",
// 			Test:    testProductOptionCreation,
// 		},
// 		{
// 			Message: "delete created product option",
// 			Test:    testDeleteProductOption,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionDeletion(t *testing.T) {
// 	t.Skip()

// 	testOptionName := "example_option_to_delete"
// 	var createdOptionID uint64
// 	testProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		resp, err := createProductOptionForProduct(existentID, newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testDeleteProductOption := func(t *testing.T) {
// 		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		body := turnResponseBodyIntoString(t, resp)
// 		assert.NotEmpty(t, body, "deletion body response should not be empty")
// 		// expected := `
// 		// `
// 		// actual := cleanAPIResponseBody(body)
// 		// assert.Equal(t, expected, actual, "product option deletion route should respond with affirmative message upon successful deletion")
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product option",
// 			Test:    testProductOptionCreation,
// 		},
// 		{
// 			Message: "delete created product option",
// 			Test:    testDeleteProductOption,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionDeletionForNonexistentOption(t *testing.T) {
// 	t.Skip()

// 	resp, err := deleteProductOption(nonexistentID)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := `{"status":404,"message":"The product option you were looking for (id '999999999') does not exist"}`
// 	assert.Equal(t, expected, actual, "product option deletion route should respond with affirmative message upon successful deletion")
// }

// func TestProductOptionCreationWithInvalidInput(t *testing.T) {
// 	t.Skip()
// 	resp, err := createProductOptionForProduct(existentID, exampleGarbageInput)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	assert.Equal(t, expectedBadRequestResponse, actual, "product option creation route should respond with failure message when you provide it invalid input")
// }

// func TestProductOptionCreationWithAlreadyExistentName(t *testing.T) {
// 	t.Skip()
// 	var createdOptionID uint64
// 	var createdProductRootID uint64
// 	testOptionName := "already-existent-option"
// 	testSKU := "test-duplicate-option-sku"

// 	testProductCreation := func(t *testing.T) {
// 		newProductJSON := createProductCreationBody(testSKU, "")
// 		resp, err := createProduct(newProductJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)
// 		body := turnResponseBodyIntoString(t, resp)
// 		createdProductRootID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		createdProductRootIDString := strconv.Itoa(int(createdProductRootID))
// 		resp, err := createProductOptionForProduct(createdProductRootIDString, newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testDuplicateProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		createdProductRootIDString := strconv.Itoa(int(createdProductRootID))
// 		resp, err := createProductOptionForProduct(createdProductRootIDString, newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusBadRequest)
// 	}

// 	testDeleteProductOption := func(t *testing.T) {
// 		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		actual := turnResponseBodyIntoString(t, resp)
// 		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product",
// 			Test:    testProductCreation,
// 		},
// 		{
// 			Message: "create product option",
// 			Test:    testProductOptionCreation,
// 		},
// 		{
// 			Message: "create product option again",
// 			Test:    testDuplicateProductOptionCreation,
// 		},
// 		{
// 			Message: "delete created product option",
// 			Test:    testDeleteProductOption,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionUpdate(t *testing.T) {
// 	t.Skip()
// 	testSKU := "testing_product_options"
// 	testOptionName := "example_option_to_update"
// 	var createdOptionID uint64
// 	var createdRootID uint64

// 	updatedOptionName := "not_the_same"

// 	testProductCreation := func(t *testing.T) {
// 		newProductJSON := createProductCreationBody(testSKU, "")
// 		resp, err := createProduct(newProductJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)
// 		body := turnResponseBodyIntoString(t, resp)
// 		createdRootID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		resp, err := createProductOptionForProduct(strconv.Itoa(int(createdRootID)), newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testUpdateProductOption := func(t *testing.T) {
// 		updatedOptionJSON := createProductOptionBody(updatedOptionName)
// 		resp, err := updateProductOption(strconv.Itoa(int(createdOptionID)), updatedOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		body := turnResponseBodyIntoString(t, resp)
// 		actual := cleanAPIResponseBody(body)
// 		expected := minifyJSON(t, `
// 			{
// 				"name": "not_the_same",
// 				"values": [{
// 					"value": "one"
// 				}, {
// 					"value": "two"
// 				}, {
// 					"value": "three"
// 				}]
// 			}
// 		`)
// 		assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
// 	}

// 	testDeleteProductOption := func(t *testing.T) {
// 		resp, err := deleteProductOption(strconv.Itoa(int(createdOptionID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		actual := turnResponseBodyIntoString(t, resp)
// 		assert.Equal(t, "", actual, "product option deletion route should respond with affirmative message upon successful deletion")
// 	}

// 	testDeleteProduct := func(t *testing.T) {
// 		resp, err := deleteProductRoot(strconv.Itoa(int(createdRootID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product to add option to",
// 			Test:    testProductCreation,
// 		},
// 		{
// 			Message: "create product option",
// 			Test:    testProductOptionCreation,
// 		},
// 		{
// 			Message: "update product option",
// 			Test:    testUpdateProductOption,
// 		},
// 		{
// 			Message: "delete created product option",
// 			Test:    testDeleteProductOption,
// 		},
// 		{
// 			Message: "delete created product root",
// 			Test:    testDeleteProduct,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionUpdateWithInvalidInput(t *testing.T) {
// 	t.Skip()
// 	resp, err := updateProductOption(existentID, exampleGarbageInput)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
// }

// func TestProductOptionUpdateForNonexistentOption(t *testing.T) {
// 	t.Skip()
// 	updatedOptionName := "nonexistent-not-the-same"
// 	updatedOptionJSON := createProductOptionBody(updatedOptionName)
// 	resp, err := updateProductOption(nonexistentID, updatedOptionJSON)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 404,
// 			"message": "The product option you were looking for (id '999999999') does not exist"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
// }

// func TestProductOptionValueCreation(t *testing.T) {
// 	t.Skip()

// 	var createdOptionValueID uint64
// 	testValue := "test-value-creation"
// 	testCreateProductOptionValue := func(t *testing.T) {
// 		newOptionValueJSON := createProductOptionValueBody(testValue)
// 		resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionValueID = retrieveIDFromResponseBody(t, body)
// 		actual := cleanAPIResponseBody(body)
// 		expected := minifyJSON(t, fmt.Sprintf(`
// 			{
// 				"value": "%s"
// 			}
// 		`, testValue))
// 		assert.Equal(t, expected, actual, "product option value creation route should respond with created product option body")
// 	}

// 	testDeleteProductOptionValue := func(t *testing.T) {
// 		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product option value",
// 			Test:    testCreateProductOptionValue,
// 		},
// 		{
// 			Message: "delete created product option value",
// 			Test:    testDeleteProductOptionValue,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionValueUpdate(t *testing.T) {
// 	t.Skip()

// 	var createdOptionID uint64
// 	var createdOptionValueID uint64
// 	testSKU := "test-value-update-sku"
// 	testOptionName := "test-value-update-obligatory-option"
// 	testValue := "test-value-update"
// 	testUpdatedValue := "not_the_same_value"

// 	createdProductRootID := createTestProduct(t, testSKU)

// 	testProductOptionCreation := func(t *testing.T) {
// 		newOptionJSON := createProductOptionCreationBody(testOptionName)
// 		createdProductRootIDString := strconv.Itoa(int(createdProductRootID))
// 		resp, err := createProductOptionForProduct(createdProductRootIDString, newOptionJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionID = retrieveIDFromResponseBody(t, body)
// 		actual := cleanAPIResponseBody(body)

// 		expected := minifyJSON(t, fmt.Sprintf(`
// 			{
// 				"name": "%s"
// 				"values": [
// 					{
// 						"value": "one"
// 					},
// 					{
// 						"value": "two"
// 					},
// 					{
// 						"value": "three"
// 					}
// 				]
// 			}
// 		`, testOptionName))
// 		assert.Equal(t, expected, actual, "product option creation route should respond with created product option body")
// 	}

// 	testCreateProductOptionValue := func(t *testing.T) {
// 		newOptionValueJSON := createProductOptionValueBody(testValue)
// 		createdOptionIDString := strconv.Itoa(int(createdOptionID))
// 		resp, err := createProductOptionValueForOption(createdOptionIDString, newOptionValueJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)

// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionValueID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testUpdateProductOptionValue := func(t *testing.T) {
// 		updatedOptionValueJSON := createProductOptionValueBody(testUpdatedValue)
// 		resp, err := updateProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)), updatedOptionValueJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		body := turnResponseBodyIntoString(t, resp)
// 		actual := cleanAPIResponseBody(body)
// 		expected := minifyJSON(t, fmt.Sprintf(`
// 			{
// 				"value": "%s"
// 			}
// 		`, testUpdatedValue))
// 		assert.Equal(t, expected, actual, "product option update response should reflect the updated fields")
// 	}

// 	testDeleteProductOptionValue := func(t *testing.T) {
// 		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product option",
// 			Test:    testProductOptionCreation,
// 		},
// 		{
// 			Message: "create product option value",
// 			Test:    testCreateProductOptionValue,
// 		},
// 		{
// 			Message: "update created product option value",
// 			Test:    testUpdateProductOptionValue,
// 		},
// 		{
// 			Message: "delete created product option value",
// 			Test:    testDeleteProductOptionValue,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionValueDeletion(t *testing.T) {
// 	t.Skip()

// 	var createdOptionValueID uint64
// 	testValue := "test-value-deletion"
// 	testCreateProductOptionValue := func(t *testing.T) {
// 		newOptionValueJSON := createProductOptionValueBody(testValue)
// 		resp, err := createProductOptionValueForOption(existentID, newOptionValueJSON)
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusCreated)
// 		body := turnResponseBodyIntoString(t, resp)
// 		createdOptionValueID = retrieveIDFromResponseBody(t, body)
// 	}

// 	testDeleteProductOptionValue := func(t *testing.T) {
// 		resp, err := deleteProductOptionValueForOption(strconv.Itoa(int(createdOptionValueID)))
// 		assert.Nil(t, err)
// 		assertStatusCode(t, resp, http.StatusOK)

// 		actual := turnResponseBodyIntoString(t, resp)
// 		assert.NotEmpty(t, actual, "product option deletion route should respond with affirmative message upon successful deletion")
// 	}

// 	subtests := []subtest{
// 		{
// 			Message: "create product option value",
// 			Test:    testCreateProductOptionValue,
// 		},
// 		{
// 			Message: "delete created product option value",
// 			Test:    testDeleteProductOptionValue,
// 		},
// 	}
// 	runSubtestSuite(t, subtests)
// }

// func TestProductOptionValueDeletionForNonexistentOptionValue(t *testing.T) {
// 	t.Skip()

// 	resp, err := deleteProductOptionValueForOption(nonexistentID)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := `{"status":404,"message":"The product option value you were looking for (id '999999999') does not exist"}`
// 	assert.Equal(t, expected, actual, "product option deletion route should respond with affirmative message upon successful deletion")
// }

// func TestProductOptionValueCreationWithInvalidInput(t *testing.T) {
// 	t.Skip()
// 	resp, err := createProductOptionValueForOption(existentID, exampleGarbageInput)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	assert.Equal(t, expectedBadRequestResponse, actual, "product option value creation route should respond with failure message when you provide it invalid input")
// }

// func TestProductOptionValueCreationWithAlreadyExistentValue(t *testing.T) {
// 	t.Skip()

// 	alreadyExistentValue := "blue"
// 	existingOptionJSON := createProductOptionValueBody(alreadyExistentValue)
// 	resp, err := createProductOptionValueForOption(existentID, existingOptionJSON)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	actual := turnResponseBodyIntoString(t, resp)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 400,
// 			"message": "product option value 'blue' already exists for option ID 1"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product option value creation route should respond with failure message when you try to create a value that already exists")
// }

// func TestProductOptionValueUpdateWithInvalidInput(t *testing.T) {
// 	t.Skip()
// 	resp, err := updateProductOptionValueForOption(existentID, exampleGarbageInput)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusBadRequest)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	assert.Equal(t, expectedBadRequestResponse, actual, "product option update route should respond with failure message when you provide it invalid input")
// }

// func TestProductOptionValueUpdateForNonexistentOption(t *testing.T) {
// 	t.Skip()

// 	obligatoryValue := "whatever"
// 	updatedOptionValueJSON := createProductOptionValueBody(obligatoryValue)
// 	resp, err := updateProductOptionValueForOption(nonexistentID, updatedOptionValueJSON)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusNotFound)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	expected := minifyJSON(t, `
// 		{
// 			"status": 404,
// 			"message": "The product option value you were looking for (id '999999999') does not exist"
// 		}
// 	`)
// 	assert.Equal(t, expected, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
// }

// func TestProductOptionValueUpdateForAlreadyExistentValue(t *testing.T) {
// 	// Say you have a product option called `color`, and it has three values (`red`, `green`, and `blue`).
// 	// Let's say you try to change `red` to `blue` for whatever reason. That will fail at the database level,
// 	// because the schema ensures a unique combination of value and option ID and archived date.
// 	t.Skip()

// 	duplicatedOptionValueJSON := createProductOptionValueBody("medium")
// 	resp, err := updateProductOptionValueForOption("4", duplicatedOptionValueJSON)
// 	assert.Nil(t, err)
// 	assertStatusCode(t, resp, http.StatusInternalServerError)

// 	body := turnResponseBodyIntoString(t, resp)
// 	actual := cleanAPIResponseBody(body)
// 	assert.Equal(t, expectedInternalErrorResponse, actual, "product option update route should respond with 404 message when you try to delete a product that doesn't exist")
// }
