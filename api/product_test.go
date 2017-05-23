package api

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
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
	expectedUPC := NullString{
		sql.NullString{
			String: "1234567890",
		},
	}
	expected := &Product{
		ProductProgenitorID: 1,
		SKU:                 "example_sku",
		Name:                "Test",
		UPC:                 expectedUPC,
		Quantity:            666,
		Price:               12.34,
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", exampleProduct)
	actual, err := loadProductInput(req)

	assert.Nil(t, err)
	assert.Equal(t, actual, expected, "valid product input should parse into a proper product struct")
}

func TestNonExistentProductResponder(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	respondThatProductDoesNotExist(req, w, "example")

	assert.Equal(t, w.Body.String(), "No product with the sku 'example' found\n", "response should indicate the product was not found")
	assert.Equal(t, w.Code, 404, "status code should be 404")
}
