package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

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
		UPC:                 NullString{sql.NullString{String: "1234567890", Valid: true}},
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

func setExpectationsForProductExistence(mock sqlmock.Sqlmock, SKU string, exists bool) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	skuExistenceQuery := buildProductExistenceQuery(SKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows)
}

func setExpectationsForProductListQuery(mock sqlmock.Sqlmock) {
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...)

	allProductsRetrievalQuery, _ := buildAllProductsRetrievalQuery(defaultQueryFilter)
	mock.ExpectQuery(formatConstantQueryForSQLMock(allProductsRetrievalQuery)).
		// WithArgs(args).
		WillReturnRows(exampleRows)
}

func ensureExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestValidateProductUpdateInputWithValidInput(t *testing.T) {
	expected := &Product{
		ProductProgenitorID: 1,
		SKU:                 exampleSKU,
		Name:                "Test",
		UPC:                 NullString{sql.NullString{String: "1234567890"}},
		Quantity:            666,
		Price:               12.34,
	}

	req := httptest.NewRequest("GET", "http://example.com", exampleProductCreationInput)
	actual, err := validateProductUpdateInput(req)

	assert.Nil(t, err)
	assert.Equal(t, actual, expected, "valid product input should parse into a proper product struct")
}

func TestValidateProductUpdateInputWithInvalidInput(t *testing.T) {
	exampleInput := strings.NewReader(`{"testing": true}`)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateProductUpdateInput(req)

	assert.NotNil(t, err)
}

func TestValidateProductUpdateInputWithInvalidSKU(t *testing.T) {
	exampleInput := strings.NewReader(`{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateProductUpdateInput(req)

	assert.NotNil(t, err)
}

func setupExpectationsForProductRetrieval(mock sqlmock.Sqlmock) {
	exampleRows := sqlmock.NewRows(plainProductHeaders).
		AddRow(examplePlainProductData...)

	skuRetrievalQuery := buildProductRetrievalQuery(exampleSKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuRetrievalQuery)).
		WithArgs(exampleSKU).
		WillReturnRows(exampleRows)
}

func TestRetrieveProductsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductListQuery(mock)

	products, err := retrieveProductsFromDB(db, defaultQueryFilter)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(products), "there should be 3 products")
	ensureExpectationsWereMet(t, mock)
}

func TestDeleteProductBySKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	skuDeletionQuery := buildProductDeletionQuery(exampleSKU)
	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = deleteProductBySKU(db, exampleSKU)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func setExpectationsForProductUpdate(mock sqlmock.Sqlmock) {
	exampleRows := sqlmock.NewRows(plainProductHeaders).
		AddRow(examplePlainProductData...)

	productUpdateQuery, _ := buildProductUpdateQuery(exampleProduct)
	mock.ExpectQuery(formatConstantQueryForSQLMock(productUpdateQuery)).
		WithArgs(
			exampleProduct.Name,
			exampleProduct.OnSale,
			exampleProduct.Price,
			exampleProduct.Quantity,
			exampleProduct.SalePrice,
			exampleProduct.SKU,
			exampleProduct.UPC,
			exampleProduct.ID,
		).WillReturnRows(exampleRows)
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
	exampleRows := sqlmock.NewRows(plainProductHeaders).
		AddRow(examplePlainProductData...)

	productCreationQuery, _ := buildProductCreationQuery(exampleProduct)
	mock.ExpectQuery(formatConstantQueryForSQLMock(productCreationQuery)).
		WithArgs(
			exampleProduct.ProductProgenitorID,
			exampleProduct.SKU,
			exampleProduct.Name,
			exampleProduct.UPC,
			exampleProduct.Quantity,
			exampleProduct.OnSale,
			exampleProduct.Price,
			exampleProduct.SalePrice,
		).WillReturnRows(exampleRows)
}

func TestCreateProductInDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductCreation(mock)

	err = createProductInDB(db, exampleProduct)
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

	req, _ := http.NewRequest("HEAD", "/v1/product/example", nil)
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

	req, _ := http.NewRequest("HEAD", "/v1/product/unreal", nil)
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
	skuJoinRetrievalQuery := buildCompleteProductRetrievalQuery(exampleProduct.SKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuJoinRetrievalQuery)).
		WithArgs(exampleProduct.SKU).
		WillReturnRows(exampleRows)

	req, _ := http.NewRequest("GET", "/v1/product/skateboard", nil)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductListHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductListQuery(mock)

	req, _ := http.NewRequest("GET", "/v1/products", nil)

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
	assert.Equal(t, uint64(len(actual.Data)), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductDeletionHandlerWithExistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true)

	skuDeletionQuery := buildProductDeletionQuery(exampleSKU)
	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("DELETE", "/v1/product/example", nil)
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

	req, _ := http.NewRequest("DELETE", "/v1/product/example", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}
