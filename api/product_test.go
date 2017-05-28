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

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleSKU        = "example"
	exampleTimeString = "2017-01-01 12:00:00.000000"

	exampleProductCreationInput = `
		{
			"sku": "new-product",
			"name": "New Product",
			"upc": "0123456789",
			"quantity": 123,
			"on_sale": false,
			"price": 0,
			"sale_price": null,
			"description": "This is a new product.",
			"taxable": true,
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

	exampleProductUpdateInput = `
		{
			"sku": "example",
			"name": "Test",
			"quantity": 666,
			"upc": "1234567890",
			"on_sale": false,
			"price": 12.34,
			"sale_price": null
		}
	`
)

var arbitraryError error
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

	arbitraryError = fmt.Errorf("arbitrary error")
	plainProductHeaders = []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at"}
	examplePlainProductData = []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, exampleTime, nil, nil}

	productJoinHeaders = []string{"id", "product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price", "created_at", "updated_at", "archived_at", "id", "name", "description", "taxable", "price", "product_weight", "product_height", "product_width", "product_length", "package_weight", "package_height", "package_width", "package_length", "created_at", "updated_at", "archived_at"}
	exampleProductJoinData = []driver.Value{10, 2, "skateboard", "Skateboard", "1234567890", 123, false, 12.34, nil, exampleTime, nil, nil, 2, "Skateboard", "This is a skateboard. Please wear a helmet.", false, 99.99, 8, 7, 6, 5, 4, 3, 2, 1, exampleTime, nil, nil}
	exampleProduct = &Product{
		ProductProgenitor: ProductProgenitor{
			ID:            2,
			Name:          "Skateboard",
			Price:         99.99,
			Description:   "This is a skateboard. Please wear a helmet.",
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
			CreatedAt:     exampleTime,
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
}

func setExpectationsForProductExistence(mock sqlmock.Sqlmock, SKU string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	skuExistenceQuery := buildProductExistenceQuery(SKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductListQuery(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...).
		AddRow(exampleProductJoinData...)

	allProductsRetrievalQuery, _ := buildAllProductsRetrievalQuery(defaultQueryFilter)
	mock.ExpectQuery(formatConstantQueryForSQLMock(allProductsRetrievalQuery)).WillReturnRows(exampleRows).WillReturnError(err)
}

func ensureExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestValidateProductUpdateInputWithValidInput(t *testing.T) {
	expected := &Product{
		SKU:       exampleSKU,
		Name:      "Test",
		UPC:       NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:  666,
		Price:     12.34,
		SalePrice: NullFloat64{sql.NullFloat64{Float64: 0, Valid: true}},
	}

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleProductUpdateInput))
	actual, err := validateProductUpdateInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "valid product input should parse into a proper product struct")
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

func setExpectationsForProductRetrieval(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...)

	skuRetrievalQuery := formatConstantQueryForSQLMock(buildCompleteProductRetrievalQuery(exampleSKU))
	mock.ExpectQuery(skuRetrievalQuery).WithArgs(exampleSKU).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestRetrieveProductFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductRetrieval(mock, nil)

	actual, err := retrieveProductFromDB(db, exampleSKU)
	assert.Nil(t, err)
	assert.Equal(t, exampleProduct, actual, "expected and actual products should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductFromDBWhenDBReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductRetrieval(mock, sql.ErrNoRows)

	_, err = retrieveProductFromDB(db, exampleSKU)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductListQuery(mock, nil)

	products, err := retrieveProductsFromDB(db, defaultQueryFilter)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(products), "there should be 3 products")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductsFromDBWhenDBReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductListQuery(mock, sql.ErrNoRows)

	_, err = retrieveProductsFromDB(db, defaultQueryFilter)
	assert.NotNil(t, err)
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

func setExpectationsForProductUpdate(mock sqlmock.Sqlmock, err error) {
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
		).WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductUpdate(mock, nil)

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

func TestValidateProductCreationInput(t *testing.T) {
	expected := &ProductCreationInput{
		Description:   "This is a new product.",
		Taxable:       true,
		ProductWeight: 9,
		ProductHeight: 9,
		ProductWidth:  9,
		ProductLength: 9,
		PackageWeight: 9,
		PackageHeight: 9,
		PackageWidth:  9,
		PackageLength: 9,
		SKU:           "new-product",
		Name:          "New Product",
		UPC:           "0123456789",
		Quantity:      123,
	}

	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(exampleProductCreationInput))
	actual, err := validateProductCreationInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "valid product input should parse into a proper product struct")
}

func TestValidateProductCreationInputWithInvalidInput(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateProductCreationInput(req)

	assert.NotNil(t, err)
}

func TestValidateProductCreationInputWithInvalidSKU(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader(`{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`))
	_, err := validateProductCreationInput(req)

	assert.NotNil(t, err)
}

//////////////////////////////////////////////////////////////
//                        ,-.             __                //
//  HTTP                ,'   `---.___.---'  `.              //
//    handler         ,'   ,-                 `-._          //
//      tests       ,'    /                       \         //
//               ,\/     /                        \\        //
//           )`._)>)     |                         \\       //
//           `>,'    _   \                  /       |\      //
//             )      \   |   |            |        |\\     //
//    .   ,   /        \  |    `.          |        | ))    //
//    \`. \`-'          )-|      `.        |        /((     //
//     \ `-`   a`     _/ ;\ _     )`-.___.--\      /  `'    //
//      `._         ,'    \`j`.__/        \  `.    \        //
//        / ,    ,'       _)\   /`        _) ( \   /        //
//        \__   /        /nn_) (         /nn__\_) (         //
//          `--'     hjw   /nn__\             /nn__\        //
//////////////////////////////////////////////////////////////

func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true, nil)

	req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, "unreal", false, nil)

	req, err := http.NewRequest("HEAD", "/v1/product/unreal", nil)
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductExistenceHandlerWithExistenceCheckerError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, "unreal", false, arbitraryError)

	req, err := http.NewRequest("HEAD", "/v1/product/unreal", nil)
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func setExpectationsForSingleProductRetrieval(mock sqlmock.Sqlmock, err error) {

	exampleRows := sqlmock.NewRows(productJoinHeaders).
		AddRow(exampleProductJoinData...)
	skuJoinRetrievalQuery := buildCompleteProductRetrievalQuery(exampleProduct.SKU)
	mock.ExpectQuery(formatConstantQueryForSQLMock(skuJoinRetrievalQuery)).
		WithArgs(exampleProduct.SKU).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForSingleProductRetrieval(mock, nil)

	req, err := http.NewRequest("GET", "/v1/product/skateboard", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductRetrievalHandlerWithDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForSingleProductRetrieval(mock, sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/v1/product/skateboard", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductListHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductListQuery(mock, nil)

	req, err := http.NewRequest("GET", "/v1/products", nil)
	assert.Nil(t, err)

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

func TestProductListHandlerWithDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductListQuery(mock, arbitraryError)

	req, err := http.NewRequest("GET", "/v1/products", nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func setExpectationsForProductUpdateHandler(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(plainProductHeaders).
		AddRow(examplePlainProductData...)

	productUpdateQuery, _ := buildProductUpdateQuery(exampleProduct)
	mock.ExpectQuery(formatConstantQueryForSQLMock(productUpdateQuery)).
		WithArgs(
			"Test",
			false,
			12.34000015258789,
			666,
			float64(0),
			"example",
			nil,
			0,
		).WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestProductUpdateHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true, nil)
	setExpectationsForProductUpdateHandler(mock, nil)

	req, err := http.NewRequest("PUT", "/v1/product/example", strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductUpdateHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, false, nil)

	req, err := http.NewRequest("PUT", "/v1/product/example", strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductDeletionHandlerWithExistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, true, nil)

	skuDeletionQuery := buildProductDeletionQuery(exampleSKU)
	mock.ExpectExec(formatConstantQueryForSQLMock(skuDeletionQuery)).
		WithArgs(exampleSKU).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("DELETE", "/v1/product/example", nil)
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductDeletionHandlerWithNonexistentProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductExistence(mock, exampleSKU, false, nil)

	req, err := http.NewRequest("DELETE", "/v1/product/example", nil)
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 404, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}
