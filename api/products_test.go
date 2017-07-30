package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleTimeAvailableString = "2016-12-31T12:00:00Z"
	badSKUUpdateJSON           = `{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`
	exampleProductID           = uint64(2)
	exampleProductUpdateInput  = `
		{
			"sku": "example",
			"name": "Test",
			"quantity": 666,
			"upc": "1234567890",
			"price": 12.34
		}
	`
)

func setExpectationsForProductExistence(mock sqlmock.Sqlmock, SKU string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductListQuery(mock sqlmock.Sqlmock, p *Product, err error) {
	productHeaders := []string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"quantity_per_package",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}
	exampleProductData := []driver.Value{
		p.ID,
		p.ProductRootID,
		p.Name,
		p.Subtitle.String,
		p.Description,
		p.SKU,
		p.UPC.String,
		p.Manufacturer.String,
		p.Brand.String,
		p.Quantity,
		p.QuantityPerPackage,
		p.Taxable,
		p.Price,
		p.OnSale,
		p.SalePrice,
		p.Cost,
		p.ProductWeight,
		p.ProductHeight,
		p.ProductWidth,
		p.ProductLength,
		p.PackageWeight,
		p.PackageHeight,
		p.PackageWidth,
		p.PackageLength,
		p.AvailableOn,
		p.CreatedOn,
		p.UpdatedOn,
		p.ArchivedOn,
	}

	exampleRows := sqlmock.NewRows(productHeaders).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...)

	allProductsRetrievalQuery, _ := buildProductListQuery(defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(allProductsRetrievalQuery)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRetrieval(mock sqlmock.Sqlmock, sku string, p *Product, err error) {
	productHeaders := []string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"quantity_per_package",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}
	exampleProductData := []driver.Value{
		p.ID,
		p.ProductRootID,
		p.Name,
		p.Subtitle.String,
		p.Description,
		p.SKU,
		p.UPC.String,
		p.Manufacturer.String,
		p.Brand.String,
		p.Quantity,
		p.QuantityPerPackage,
		p.Taxable,
		p.Price,
		p.OnSale,
		p.SalePrice,
		p.Cost,
		p.ProductWeight,
		p.ProductHeight,
		p.ProductWidth,
		p.ProductLength,
		p.PackageWeight,
		p.PackageHeight,
		p.PackageWidth,
		p.PackageLength,
		p.AvailableOn,
		p.CreatedOn,
		p.UpdatedOn,
		p.ArchivedOn,
	}

	exampleRows := sqlmock.NewRows(productHeaders).AddRow(exampleProductData...)
	skuRetrievalQuery := formatQueryForSQLMock(completeProductRetrievalQuery)
	mock.ExpectQuery(skuRetrievalQuery).
		WithArgs(sku).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductUpdate(mock sqlmock.Sqlmock, p *Product, err error) {
	productHeaders := []string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"quantity_per_package",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}
	exampleProductData := []driver.Value{
		p.ID,
		p.ProductRootID,
		p.Name,
		p.Subtitle.String,
		p.Description,
		p.SKU,
		p.UPC.String,
		p.Manufacturer.String,
		p.Brand.String,
		p.Quantity,
		p.QuantityPerPackage,
		p.Taxable,
		p.Price,
		p.OnSale,
		p.SalePrice,
		p.Cost,
		p.ProductWeight,
		p.ProductHeight,
		p.ProductWidth,
		p.ProductLength,
		p.PackageWeight,
		p.PackageHeight,
		p.PackageWidth,
		p.PackageLength,
		p.AvailableOn,
		p.CreatedOn,
		p.UpdatedOn,
		p.ArchivedOn,
	}

	exampleRows := sqlmock.NewRows(productHeaders).AddRow(exampleProductData...)
	productUpdateQuery, queryArgs := buildProductUpdateQuery(p)
	args := argsToDriverValues(queryArgs)
	mock.ExpectQuery(formatQueryForSQLMock(productUpdateQuery)).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductCreation(mock sqlmock.Sqlmock, p *Product, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(p.ID, generateExampleTimeForTests())
	productCreationQuery, args := buildProductCreationQuery(p)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(productCreationQuery)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductCreationFromOptions(mock sqlmock.Sqlmock, ps []*Product, optionCount uint, err error, errorOnBridgeEntries bool, errorIndex int) {
	for i, p := range ps {
		p.ID = uint64(i + 1)
		if i == errorIndex && err != nil {
			if errorOnBridgeEntries {
				setExpectationsForProductCreation(mock, p, nil)
				setExpectationsForProductValueBridgeEntryCreation(mock, p.ID, make([]uint64, optionCount), err)
			} else {
				setExpectationsForProductCreation(mock, p, err)
			}
			return
		}
		setExpectationsForProductCreation(mock, p, nil)
		setExpectationsForProductValueBridgeEntryCreation(mock, p.ID, make([]uint64, optionCount), nil)
	}
}

func setExpectationsForProductUpdateHandler(mock sqlmock.Sqlmock, p *Product, err error) {
	productHeaders := []string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"quantity_per_package",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}
	exampleProductData := []driver.Value{
		p.ID,
		p.ProductRootID,
		p.Name,
		p.Subtitle.String,
		p.Description,
		p.SKU,
		p.UPC.String,
		p.Manufacturer.String,
		p.Brand.String,
		p.Quantity,
		p.QuantityPerPackage,
		p.Taxable,
		p.Price,
		p.OnSale,
		p.SalePrice,
		p.Cost,
		p.ProductWeight,
		p.ProductHeight,
		p.ProductWidth,
		p.ProductLength,
		p.PackageWeight,
		p.PackageHeight,
		p.PackageWidth,
		p.PackageLength,
		p.AvailableOn,
		p.CreatedOn,
		p.UpdatedOn,
		p.ArchivedOn,
	}

	exampleRows := sqlmock.NewRows(productHeaders).AddRow(exampleProductData...)
	productUpdateQuery, _ := buildProductUpdateQuery(p)
	mock.ExpectQuery(formatQueryForSQLMock(productUpdateQuery)).
		WithArgs(
			p.Cost,
			p.Name,
			p.Price,
			p.Quantity,
			p.SKU,
			p.UPC.String,
			p.ID,
		).WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductDeletion(mock sqlmock.Sqlmock, sku string, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productDeletionQuery)).
		WithArgs(sku).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestRetrieveProductFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}
	exampleProduct.Subtitle.Valid = true
	exampleProduct.Brand.Valid = true
	exampleProduct.Manufacturer.Valid = true

	setExpectationsForProductRetrieval(testUtil.Mock, exampleSKU, exampleProduct, nil)

	actual, err := retrieveProductFromDB(testUtil.DB, exampleSKU)
	assert.Nil(t, err)
	assert.Equal(t, *exampleProduct, actual, "expected and actual products should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductFromDBWhenDBReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleSKU, exampleProduct, sql.ErrNoRows)

	_, err := retrieveProductFromDB(testUtil.DB, exampleSKU)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDeleteProductBySKU(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductDeletion(testUtil.Mock, exampleSKU, nil)

	err := deleteProductBySKU(testUtil.DB, exampleSKU)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestDeleteProductBySKUReturnsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductDeletion(testUtil.Mock, exampleSKU, arbitraryError)

	err := deleteProductBySKU(testUtil.DB, exampleSKU)
	assert.Equal(t, err, arbitraryError, "deleteProductBySKU should return errors when it encounters them")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateProductInDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductUpdate(testUtil.Mock, exampleProduct, nil)

	err := updateProductInDatabase(testUtil.DB, exampleProduct)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	newID, createdOn, err := createProductInDB(tx, exampleProduct)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductID, newID, "createProductInDB should return the created ID")
	assert.Equal(t, generateExampleTimeForTests(), createdOn, "createProductInDB should return the created ID")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductExistenceHandlerWithExistingProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, true, nil)

	req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductExistenceHandlerWithNonexistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, "unreal", false, nil)

	req, err := http.NewRequest("HEAD", "/v1/product/unreal", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductExistenceHandlerWithExistenceCheckerError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, "unreal", false, arbitraryError)

	req, err := http.NewRequest("HEAD", "/v1/product/unreal", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRetrievalHandlerWithExistingProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, nil)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRetrievalHandlerWithDBError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRetrievalHandlerWithNonexistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductListHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForRowCount(testUtil.Mock, "products", defaultQueryFilter, 3, nil)
	setExpectationsForProductListQuery(testUtil.Mock, exampleProduct, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")

	expected := &ProductsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductsResponse{}
	err = json.NewDecoder(strings.NewReader(testUtil.Response.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	assert.Equal(t, uint64(len(actual.Data)), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductListHandlerWithErrorRetrievingCount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForRowCount(testUtil.Mock, "products", defaultQueryFilter, 3, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductListHandlerWithDBError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForRowCount(testUtil.Mock, "products", defaultQueryFilter, 3, nil)
	setExpectationsForProductListQuery(testUtil.Mock, exampleProduct, arbitraryError)

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	exampleUpdatedProduct := &Product{
		DBRow: DBRow{
			ID:        exampleProduct.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:      "example",
		Name:     "Test",
		UPC:      NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity: 666,
		Cost:     50.00,
		Price:    12.34,
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, nil)
	setExpectationsForProductUpdateHandler(testUtil.Mock, exampleUpdatedProduct, nil)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandlerWithNonexistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandlerWithInputValidationError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPatch, "/v1/product/example", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandlerWithSKUValidationError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, nil)

	req, err := http.NewRequest(http.MethodPatch, "/v1/product/skateboard", strings.NewReader(badSKUUpdateJSON))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandlerWithDBErrorRetrievingProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, arbitraryError)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandlerWithDBErrorUpdatingProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	exampleUpdatedProduct := &Product{
		DBRow: DBRow{
			ID:        exampleProduct.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:      "example",
		Name:     "Test",
		UPC:      NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity: 666,
		Cost:     50.00,
		Price:    12.34,
	}

	setExpectationsForProductRetrieval(testUtil.Mock, exampleProduct.SKU, exampleProduct, nil)
	setExpectationsForProductUpdateHandler(testUtil.Mock, exampleUpdatedProduct, arbitraryError)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, true, nil)
	setExpectationsForProductDeletion(testUtil.Mock, exampleSKU, nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductDeletionHandlerWithNonexistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, false, nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductDeletionHandlerWithErrorEncounteredDeletingProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, true, nil)
	setExpectationsForProductDeletion(testUtil.Mock, exampleSKU, arbitraryError)

	req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	expectedFirstOption := &Product{}
	expectedSecondOption := &Product{}
	expectedThirdOption := &Product{}

	*expectedFirstOption = *exampleProduct
	*expectedSecondOption = *exampleProduct
	*expectedThirdOption = *exampleProduct

	expectedFirstOption.OptionSummary = "something: one"
	expectedFirstOption.SKU = "skateboard_one"
	expectedSecondOption.OptionSummary = "something: two"
	expectedSecondOption.SKU = "skateboard_two"
	expectedThirdOption.OptionSummary = "something: three"
	expectedThirdOption.SKU = "skateboard_three"

	expectedCreatedProducts := []*Product{expectedFirstOption, expectedSecondOption, expectedThirdOption}

	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": [{
				"name": "something",
				"values": [
					"one",
					"two",
					"three"
				]
			}]
		}
	`

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleRoot.ID, nil)
	setExpectationsForMultipleProductOptionValuesCreation(testUtil.Mock, expectedCreatedProductOption.Values, nil, -1)
	setExpectationsForProductCreationFromOptions(testUtil.Mock, expectedCreatedProducts, 1, nil, false, -1)
	testUtil.Mock.ExpectCommit()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorValidatingInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWhereCommitReturnsAnError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": []
		}
	`
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "Skateboard",
		SKU:           "skateboard",
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(arbitraryError)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWhereTransactionFailsToBegin(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": []
		}
	`
	setExpectationsForProductRootSKUExistence(testUtil.Mock, "skateboard", false, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(arbitraryError)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorCreatingProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	exampleProductCreationInput := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1
		}
	`

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, arbitraryError)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithoutOptions(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleProductCreationInput := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1
		}
	`
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "Skateboard",
		SKU:           "skateboard",
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, nil)
	testUtil.Mock.ExpectCommit()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithInvalidProductInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(badSKUUpdateJSON))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerForAlreadyExistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleProductCreationInput := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1
		}
	`
	setExpectationsForProductRootSKUExistence(testUtil.Mock, "skateboard", true, nil)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorCreatingOptions(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	exampleProductCreationInputWithOptions := fmt.Sprintf(`
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 99.99,
			"cost": 50,
			"description": "This is a skateboard. Please wear a helmet.",
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"available_on": "%s",
			"options": [{
				"name": "something",
				"values": [
					"one",
					"two",
					"three"
				]
			}]
		}
	`, exampleTimeAvailableString)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleRoot.ID, arbitraryError)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWhereProductCreationFails(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductCreationInput := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1
		}
	`
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "Skateboard",
		SKU:           "skateboard",
		Price:         12.34,
		Cost:          5,
		Description:   "This is a skateboard. Please wear a helmet.",
		Taxable:       true,
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, arbitraryError)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorCreatingOptionProducts(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	expectedFirstOption := &Product{}
	expectedSecondOption := &Product{}
	expectedThirdOption := &Product{}

	*expectedFirstOption = *exampleProduct
	*expectedSecondOption = *exampleProduct
	*expectedThirdOption = *exampleProduct

	expectedFirstOption.OptionSummary = "something: one"
	expectedFirstOption.SKU = "skateboard_one"
	expectedSecondOption.OptionSummary = "something: two"
	expectedSecondOption.SKU = "skateboard_two"
	expectedThirdOption.OptionSummary = "something: three"
	expectedThirdOption.SKU = "skateboard_three"

	expectedCreatedProducts := []*Product{expectedFirstOption, expectedSecondOption, expectedThirdOption}

	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": [{
				"name": "something",
				"values": [
					"one",
					"two",
					"three"
				]
			}]
		}
	`

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleRoot.ID, nil)
	setExpectationsForMultipleProductOptionValuesCreation(testUtil.Mock, expectedCreatedProductOption.Values, nil, -1)
	setExpectationsForProductCreationFromOptions(testUtil.Mock, expectedCreatedProducts, 1, arbitraryError, false, 0)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorCreatingBridgeEntries(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           NullString{sql.NullString{String: "1234567890", Valid: true}},
		Quantity:      123,
		Price:         12.34,
		Cost:          5,
		Taxable:       true,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	expectedFirstOption := &Product{}
	expectedSecondOption := &Product{}
	expectedThirdOption := &Product{}

	*expectedFirstOption = *exampleProduct
	*expectedSecondOption = *exampleProduct
	*expectedThirdOption = *exampleProduct

	expectedFirstOption.OptionSummary = "something: one"
	expectedFirstOption.SKU = "skateboard_one"
	expectedSecondOption.OptionSummary = "something: two"
	expectedSecondOption.SKU = "skateboard_two"
	expectedThirdOption.OptionSummary = "something: three"
	expectedThirdOption.SKU = "skateboard_three"

	expectedCreatedProducts := []*Product{expectedFirstOption, expectedSecondOption, expectedThirdOption}

	exampleProductCreationInputWithOptions := `
		{
			"sku": "skateboard",
			"name": "Skateboard",
			"upc": "1234567890",
			"quantity": 123,
			"price": 12.34,
			"cost": 5,
			"description": "This is a skateboard. Please wear a helmet.",
			"taxable": true,
			"product_weight": 8,
			"product_height": 7,
			"product_width": 6,
			"product_length": 5,
			"package_weight": 4,
			"package_height": 3,
			"package_width": 2,
			"package_length": 1,
			"options": [{
				"name": "something",
				"values": [
					"one",
					"two",
					"three"
				]
			}]
		}
	`

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleRoot.ID, nil)
	setExpectationsForMultipleProductOptionValuesCreation(testUtil.Mock, expectedCreatedProductOption.Values, nil, -1)
	setExpectationsForProductCreationFromOptions(testUtil.Mock, expectedCreatedProducts, 1, arbitraryError, true, 0)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
