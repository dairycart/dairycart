package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	updatedAtReplacementPattern = `,"(updated_at|archived_at)":{"Time":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z","Valid":(true|false)}`
)

func setupMockRequestsAndMux(db *sql.DB) (*httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	SetupAPIRoutes(m, db)
	return httptest.NewRecorder(), m
}

func formatConstantQueryForSQLMock(query string) string {
	charsToReplace := []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"}

	for _, x := range charsToReplace {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}

	return query
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

func TestProductExistsInDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	exampleRows := sqlmock.NewRows([]string{""}).AddRow("true")
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs("example").
		WillReturnRows(exampleRows)
	exists, err := productExistsInDB(db, "example")

	assert.Nil(t, err)
	assert.True(t, exists)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestRetrieveProductsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	productJoinHeaders := []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at", "id", "name", "description", "taxable", "price", "product_weight", "product_height", "product_width", "product_length", "package_weight", "package_height", "package_width", "package_length", "created_at", "updated_at", "archived_at"}
	exampleProductJoinData := []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, "2017-05-23 12:36:43.932053", nil, nil, 2, "Skateboard", "This is a skateboard. Please wear a helmet.", false, 99.99, 8, 7, 6, 5, 4, 3, 2, 1, "2017-05-23 12:36:43.932053", nil, nil}
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...)

	mock.ExpectQuery(formatConstantQueryForSQLMock(allProductsRetrievalQuery)).
		WillReturnRows(exampleRows)

	products, err := retrieveProductsFromDB(db)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(products), "there should be 3 products")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

// HTTP handler tests
func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	res, router := setupMockRequestsAndMux(db)

	exampleRows := sqlmock.NewRows([]string{""}).AddRow("true")
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs("skateboard").
		WillReturnRows(exampleRows)

	req, _ := http.NewRequest("HEAD", "/product/skateboard", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	res, router := setupMockRequestsAndMux(db)

	exampleRows := sqlmock.NewRows([]string{""}).AddRow("false")
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs("unreal").
		WillReturnRows(exampleRows)

	req, _ := http.NewRequest("HEAD", "/product/unreal", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	res, router := setupMockRequestsAndMux(db)

	productJoinHeaders := []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at", "id", "name", "description", "taxable", "price", "product_weight", "product_height", "product_width", "product_length", "package_weight", "package_height", "package_width", "package_length", "created_at", "updated_at", "archived_at"}
	exampleProductJoinData := []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, time.Now(), nil, nil, 2, "Skateboard", "This is a skateboard. Please wear a helmet.", false, 99.99, 8, 7, 6, 5, 4, 3, 2, 1, time.Now(), nil, nil}
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuJoinRetrievalQuery)).
		WithArgs("skateboard").
		WillReturnRows(exampleRows)

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
	bodyReader := strings.NewReader(res.Body.String())
	decoder := json.NewDecoder(bodyReader)
	err = decoder.Decode(actualProduct)
	assert.Nil(t, err)

	assertProductEqualityForTests(t, expectedProduct, actualProduct)
}
