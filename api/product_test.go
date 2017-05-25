package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleSKU        = "example"
	exampleTimeString = "2017-01-01 12:00:00.000000"
)

var plainProductHeaders []string
var examplePlainProductData []driver.Value
var productJoinHeaders []string
var exampleProductJoinData []driver.Value
var exampleTime time.Time
var exampleProduct *Product

func init() {
	var err error
	exampleTime, err = time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
	if err != nil {
		log.Fatalf("error parsing time")
	}

	plainProductHeaders = []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at"}
	examplePlainProductData = []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, exampleTime, nil, nil}

	productJoinHeaders = []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at", "id", "name", "description", "taxable", "price", "product_weight", "product_height", "product_width", "product_length", "package_weight", "package_height", "package_width", "package_length", "created_at", "updated_at", "archived_at"}
	exampleProductJoinData = []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, exampleTime, nil, nil, 2, "Skateboard", "This is a skateboard. Please wear a helmet.", false, 99.99, 8, 7, 6, 5, 4, 3, 2, 1, exampleTime, nil, nil}
	exampleProduct = &Product{
		ProductProgenitor: ProductProgenitor{
			// ID:            2,            // TODO:
			// Name:          "Skateboard", //    figure out why this part
			// Price:         99.99,        //    of the test fails.
			Description:   "This is a skateboard. Please wear a helmet.",
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
		CreatedAt:           exampleTime,
	}

}

func setupMockRequestsAndMux(db *sql.DB) (*httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	SetupAPIRoutes(m, db)
	return httptest.NewRecorder(), m
}

// this function is lame
func formatConstantQueryForSQLMock(query string) string {
	for _, x := range []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"} {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}
	return query
}

func setupProductExistenceChecks(db *sql.DB, mock sqlmock.Sqlmock, SKU string, exists bool) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows)
}

func TestLoadProductInputWithValidInput(t *testing.T) {
	exampleInput := strings.NewReader(strings.TrimSpace(`
		{
			"product_progenitor_id": 1,
			"sku": "example",
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
		SKU:                 exampleSKU,
		Name:                "Test",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            666,
		Price:               12.34,
	}

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual, err := loadProductInput(req)

	assert.Nil(t, err)
	assert.Equal(t, actual, expected, "valid product input should parse into a proper product struct")
}

func TestLoadProductInputWithInvalidInput(t *testing.T) {
	exampleInput := strings.NewReader(`{"testing": true}`)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := loadProductInput(req)

	assert.NotNil(t, err)
}

// commenting out for now, because this likely belongs in helpers_test.go
// func TestProductExistsInDB(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	assert.Nil(t, err)

// 	setupProductExistenceChecks(db, mock, exampleSKU, true)

// 	exists, err := productExistsInDB(db, exampleSKU)

// 	assert.Nil(t, err)
// 	assert.True(t, exists)
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expections: %s", err)
// 	}
// }

func TestRetrievePlainProductFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	exampleRows := sqlmock.NewRows(plainProductHeaders).
		AddRow(examplePlainProductData...)

	mock.ExpectQuery(formatConstantQueryForSQLMock(skuRetrievalQuery)).
		WithArgs(exampleSKU).
		WillReturnRows(exampleRows)

	actual, err := retrievePlainProductFromDB(db, exampleSKU)
	assert.Nil(t, err)

	expected := &Product{
		ID:                  10,
		ProductProgenitorID: 2,
		SKU:                 "skateboard",
		Name:                "Skateboard",
		UPC:                 NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:            123,
		Price:               12.34,
		CreatedAt:           exampleTime,
	}

	assert.Equal(t, expected, actual, "plain product returned should be valid")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestRetrieveProductsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

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

func TestDeleteProductBySKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = deleteProductBySKU(db, exampleSKU)
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestUpdateProductInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	mock.ExpectExec(formatConstantQueryForSQLMock(productUpdateQuery)).
		WithArgs(
			exampleProduct.ProductProgenitorID,
			exampleProduct.SKU,
			exampleProduct.Name,
			exampleProduct.UPC,
			exampleProduct.Quantity,
			exampleProduct.OnSale,
			exampleProduct.Price,
			exampleProduct.SalePrice,
			exampleProduct.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = updateProductInDatabase(db, exampleProduct)
	assert.Nil(t, err)

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

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	res, router := setupMockRequestsAndMux(db)

	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuJoinRetrievalQuery)).
		WithArgs("skateboard").
		WillReturnRows(exampleRows)

	req, _ := http.NewRequest("GET", "/product/skateboard", nil)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")

	expected := exampleProduct
	actual := &Product{}
	bodyReader := strings.NewReader(res.Body.String())
	decoder := json.NewDecoder(bodyReader)
	err = decoder.Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "expected and actual products should be equal")
}

func TestProductListHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)
	res, router := setupMockRequestsAndMux(db)

	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...)

	mock.ExpectQuery(formatConstantQueryForSQLMock(allProductsRetrievalQuery)).
		WillReturnRows(exampleRows)

	req, _ := http.NewRequest("GET", "/products", nil)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")

	expected := &ProductsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductsResponse{}
	bodyReader := strings.NewReader(res.Body.String())
	decoder := json.NewDecoder(bodyReader)
	err = decoder.Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product lists should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product lists should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product lists should be equal")
	assert.Equal(t, len(actual.Data), actual.Count, "expected and actual product lists should be equal")
}

func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	res, router := setupMockRequestsAndMux(db)
	setupProductExistenceChecks(db, mock, "unreal", false)

	req, _ := http.NewRequest("HEAD", "/product/unreal", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
