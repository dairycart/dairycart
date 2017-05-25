package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
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
var exampleProductCreationInput io.Reader

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

	exampleProductCreationInput = strings.NewReader(strings.TrimSpace(`
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

func setExpectationsForProductExistence(mock sqlmock.Sqlmock, SKU string, exists bool) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows)
}

func setExpectationsForProductProgenitorExistence(mock sqlmock.Sqlmock, id int64, exists bool) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatConstantQueryForSQLMock(productProgenitorExistenceQuery)).
		WithArgs(id).
		WillReturnRows(exampleRows)
}

func setExpectationsForProductListQuery(mock sqlmock.Sqlmock) {
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...)

	mock.ExpectQuery(formatConstantQueryForSQLMock(allProductsRetrievalQuery)).
		WillReturnRows(exampleRows)
}

func ensureExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestLoadProductInputWithValidInput(t *testing.T) {
	expected := &Product{
		ProductProgenitorID: 1,
		SKU:                 exampleSKU,
		Name:                "Test",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            666,
		Price:               12.34,
	}

	req := httptest.NewRequest("GET", "http://example.com", exampleProductCreationInput)
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

func TestRetrievePlainProductFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

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
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductListQuery(mock)

	products, err := retrieveProductsFromDB(db)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(products), "there should be 3 products")
	ensureExpectationsWereMet(t, mock)
}

func TestDeleteProductBySKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = deleteProductBySKU(db, exampleSKU)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func setExpectationsForProductUpdate(mock sqlmock.Sqlmock) {
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
}

func TestUpdateProductInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductUpdate(mock)

	err = updateProductInDatabase(db, exampleProduct)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func setExpectationsForProductCreation(mock sqlmock.Sqlmock) {
	mock.ExpectExec(formatConstantQueryForSQLMock(productCreationQuery)).
		WithArgs(
			exampleProduct.ProductProgenitorID,
			exampleProduct.SKU,
			exampleProduct.Name,
			exampleProduct.UPC,
			exampleProduct.Quantity,
			exampleProduct.OnSale,
			exampleProduct.Price,
			exampleProduct.SalePrice,
		).WillReturnResult(sqlmock.NewResult(1, 1))
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductCreation(mock)

	err = createProduct(db, exampleProduct)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

// HTTP handler tests

func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true)

	req, _ := http.NewRequest("HEAD", "/product/example", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, "unreal", false)

	req, _ := http.NewRequest("HEAD", "/product/unreal", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
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
	err = json.NewDecoder(strings.NewReader(res.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "expected and actual products should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductListHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductListQuery(mock)

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
	err = json.NewDecoder(strings.NewReader(res.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	assert.Equal(t, len(actual.Data), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductDeletionHandlerWithExistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true)

	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("DELETE", "/product/example", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductDeletionHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, false)

	req, _ := http.NewRequest("DELETE", "/product/example", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}
