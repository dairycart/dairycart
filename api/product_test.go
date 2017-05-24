package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	updatedAtReplacementPattern = `,"(updated_at|archived_at)":{"Time":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z","Valid":(true|false)}`
)

func setup() (*httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	SetupAPIRoutes(m, testDB)
	return httptest.NewRecorder(), m
}

func replaceTimeStringsForTests(body string) string {
	// we can't reliably predict what the `updated_at` or `archived_at` columns
	// could possibly equal, so we strip them out of the body becuase we're bad
	// at programming.
	re := regexp.MustCompile(updatedAtReplacementPattern)
	return re.ReplaceAllString(body, "")
}

func assertProductEqualityForTests(t *testing.T, expected *Product, actual *Product) {
	assert.Equal(t, expected.ID, actual.ID, "product IDs should be equal")
	assert.Equal(t, expected.Name, actual.Name, "product Names should be equal")
	assert.Equal(t, expected.Description, actual.Description, "product Descriptions should be equal")
	assert.Equal(t, expected.Price, actual.Price, "product Prices should be equal")
	assert.Equal(t, expected.ProductWeight, actual.ProductWeight, "product ProductWeights should be equal")
	assert.Equal(t, expected.ProductHeight, actual.ProductHeight, "product ProductHeights should be equal")
	assert.Equal(t, expected.ProductWidth, actual.ProductWidth, "product ProductWidths should be equal")
	assert.Equal(t, expected.ProductLength, actual.ProductLength, "product ProductLengths should be equal")
	assert.Equal(t, expected.PackageWeight, actual.PackageWeight, "product PackageWeights should be equal")
	assert.Equal(t, expected.PackageHeight, actual.PackageHeight, "product PackageHeights should be equal")
	assert.Equal(t, expected.PackageWidth, actual.PackageWidth, "product PackageWidths should be equal")
	assert.Equal(t, expected.PackageLength, actual.PackageLength, "product PackageLengths should be equal")
	assert.Equal(t, expected.ProductProgenitorID, actual.ProductProgenitorID, "product ProductProgenitorIDs should be equal")
	assert.Equal(t, expected.SKU, actual.SKU, "product SKUs should be equal")
	assert.Equal(t, expected.Name, actual.Name, "product Names should be equal")
	assert.Equal(t, expected.UPC, actual.UPC, "product UPCs should be equal")
	assert.Equal(t, expected.Quantity, actual.Quantity, "product Quantitys should be equal")
	assert.Equal(t, expected.Price, actual.Price, "product Prices should be equal")
}

func TestLoadProductInputWithValidInput(t *testing.T) {
	exampleProduct := strings.NewReader(strings.TrimSpace(`
		{
			"product_progenitor_id": 1,
			"sku": "example_sku",
			"name": "Test",
			"quantity": 666,
			"upc": "1234567890",
			"on_sale": false,
			"price": 12.34,
			"sale_price": null
		}
	`))
	expected := &Product{
		ProductProgenitorID: 1,
		SKU:                 "example_sku",
		Name:                "Test",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            666,
		Price:               12.34,
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", exampleProduct)
	actual, err := loadProductInput(req)

	assert.Nil(t, err)
	assert.Equal(t, actual, expected, "valid product input should parse into a proper product struct")
}
func TestRetrieveProductsFromDB(t *testing.T) {
	products, err := retrieveProductsFromDB(testDB)

	assert.Nil(t, err)
	assert.Equal(t, len(products), 10, "there should be 9 products")
}

// HTTP handler tests
func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	res, router := setup()
	req, _ := http.NewRequest("HEAD", "/product/skateboard", nil)

	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
}
func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	res, router := setup()
	req, _ := http.NewRequest("HEAD", "/product/unreal", nil)

	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
}

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	res, router := setup()
	req, _ := http.NewRequest("GET", "/product/skateboard", nil)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")

	expectedProduct := &Product{
		ProductProgenitor: ProductProgenitor{
			ID:            2,
			Name:          "Skateboard",
			Description:   "This is a skateboard. Please wear a helmet.",
			Price:         99.99,
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
		},
		ID:                  10,
		ProductProgenitorID: 2,
		SKU:                 "skateboard",
		Name:                "Skateboard",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            123,
		Price:               12.34,
	}

	actualProduct := &Product{}
	rawResponseBody := res.Body.String()
	re := regexp.MustCompile(updatedAtReplacementPattern)
	responseBody := re.ReplaceAllString(rawResponseBody, "")
	bodyReader := strings.NewReader(responseBody)
	decoder := json.NewDecoder(bodyReader)
	err := decoder.Decode(actualProduct)
	assert.Nil(t, err)

	assertProductEqualityForTests(t, expectedProduct, actualProduct)
}
