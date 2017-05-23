package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setup() (*httptest.ResponseRecorder, *mux.Router) {
	//mux router with added question routes
	m := mux.NewRouter()
	SetupAPIRoutes(m, testDB)

	//The response recorder used to record HTTP responses
	return httptest.NewRecorder(), m
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
			// CreatedAt:     "",
		},
		ID:                  10,
		ProductProgenitorID: 2,
		SKU:                 "skateboard",
		Name:                "Skateboard",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            123,
		Price:               12.34,
		// CreatedAt: "",
	}
	expectedBody, err := json.Marshal(expectedProduct)
	assert.Nil(t, err)

	assert.Equal(t, string(expectedBody), res.Body.String(), "response should be a marshaled Product struct")
}
