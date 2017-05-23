package api

import (
	"database/sql"
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

func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	res, router := setup()
	req, _ := http.NewRequest("HEAD", "/product/skateboard", nil)

	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
}

func TestRetrieveProductsFromDB(t *testing.T) {
	products, err := retrieveProductsFromDB(testDB)

	assert.Nil(t, err)
	assert.Equal(t, len(products), 10, "there should be 9 products")
}
