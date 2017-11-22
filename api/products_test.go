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
	"time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleTimeAvailableString = "2016-12-31T12:00:00Z"
	exampleProductID           = uint64(2)
	badSKUUpdateJSON           = `{"sku": "pooƃ ou sᴉ nʞs sᴉɥʇ"}`
)

func createExampleHeadersAndDataFromProduct(p *Product) ([]string, []driver.Value) {
	var headers []string
	var values []driver.Value

	productMap := map[string]driver.Value{
		"id":                   p.ID,
		"product_root_id":      p.ProductRootID,
		"name":                 p.Name,
		"subtitle":             p.Subtitle,
		"description":          p.Description,
		"sku":                  p.SKU,
		"upc":                  p.UPC,
		"manufacturer":         p.Manufacturer,
		"brand":                p.Brand,
		"quantity":             p.Quantity,
		"quantity_per_package": p.QuantityPerPackage,
		"taxable":              p.Taxable,
		"price":                p.Price,
		"on_sale":              p.OnSale,
		"sale_price":           p.SalePrice,
		"cost":                 p.Cost,
		"product_weight":       p.ProductWeight,
		"product_height":       p.ProductHeight,
		"product_width":        p.ProductWidth,
		"product_length":       p.ProductLength,
		"package_weight":       p.PackageWeight,
		"package_height":       p.PackageHeight,
		"package_width":        p.PackageWidth,
		"package_length":       p.PackageLength,
		"available_on":         p.AvailableOn,
		"created_on":           p.CreatedOn,
		"updated_on":           p.UpdatedOn,
		"archived_on":          p.ArchivedOn,
	}

	for header, value := range productMap {
		headers = append(headers, header)
		values = append(values, value)
	}

	return headers, values
}

func setExpectationsForProductExistence(mock sqlmock.Sqlmock, SKU string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(skuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductListQuery(mock sqlmock.Sqlmock, p *Product, err error) {
	productHeaders, exampleProductData := createExampleHeadersAndDataFromProduct(p)

	exampleRows := sqlmock.NewRows(productHeaders).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...)

	allProductsRetrievalQuery, _ := buildProductListQuery(genereateDefaultQueryFilter())
	mock.ExpectQuery(formatQueryForSQLMock(allProductsRetrievalQuery)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

// FIXME: dump SKU from the argument list here, just use whatever sku p has
func setExpectationsForProductRetrieval(mock sqlmock.Sqlmock, sku string, p *Product, err error) {
	productHeaders, exampleProductData := createExampleHeadersAndDataFromProduct(p)
	exampleRows := sqlmock.NewRows(productHeaders).AddRow(exampleProductData...)
	skuRetrievalQuery := formatQueryForSQLMock(completeProductRetrievalQuery)
	mock.ExpectQuery(skuRetrievalQuery).
		WithArgs(sku).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRetrievalNEW(mock sqlmock.Sqlmock, sku string, p *Product, err error) {
	productHeaders, exampleProductData := createExampleHeadersAndDataFromProduct(p)
	exampleRows := sqlmock.NewRows(productHeaders).AddRow(exampleProductData...)

	productQuery := `
		SELECT
			id,
			product_root_id,
			name,
			subtitle,
			description,
			option_summary,
			sku,
			upc,
			manufacturer,
			brand,
			quantity,
			taxable,
			price,
			on_sale,
			sale_price,
			cost,
			product_weight,
			product_height,
			product_width,
			product_length,
			package_weight,
			package_height,
			package_width,
			package_length,
			quantity_per_package,
			available_on,
			created_on,
			updated_on,
			archived_on

		FROM products
		WHERE sku = $1
	`

	skuRetrievalQuery := formatQueryForSQLMock(productQuery)
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
		p.Subtitle,
		p.Description,
		p.SKU,
		p.UPC,
		p.Manufacturer,
		p.Brand,
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
	exampleRows := sqlmock.NewRows([]string{"id", "available_on", "created_on"}).AddRow(p.ID, generateExampleTimeForTests(), generateExampleTimeForTests())
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
		p.Subtitle,
		p.Description,
		p.SKU,
		p.UPC,
		p.Manufacturer,
		p.Brand,
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
			p.UPC,
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
		UPC:           "1234567890",
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
		UPC:           "1234567890",
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
		UPC:           "1234567890",
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

	newID, availableOn, createdOn, err := createProductInDB(tx, exampleProduct)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductID, newID, "createProductInDB should return the created ID")
	assert.Equal(t, generateExampleTimeForTests(), createdOn, "createProductInDB should return the created_on ID")
	assert.Equal(t, generateExampleTimeForTests(), availableOn, "createProductInDB should return the available_on ID")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductExistenceHandler(t *testing.T) {
	exampleSKU := "example"
	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
	})
	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
	t.Run("with error performing check", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductWithSKUExists", mock.Anything, exampleSKU).Return(false, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest("HEAD", "/v1/product/example", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestProductRetrievalHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
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
	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with DB error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)

		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
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
		UPC:           "1234567890",
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

	setExpectationsForRowCount(testUtil.Mock, "products", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductListQuery(testUtil.Mock, exampleProduct, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)

	expected := &ListResponse{
		Page:  1,
		Limit: 25,
		Count: 3,
	}

	actual := &ListResponse{}
	err = json.NewDecoder(strings.NewReader(testUtil.Response.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductListHandlerWithErrorRetrievingCount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForRowCount(testUtil.Mock, "products", genereateDefaultQueryFilter(), 3, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)

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
		UPC:           "1234567890",
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

	setExpectationsForRowCount(testUtil.Mock, "products", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductListQuery(testUtil.Mock, exampleProduct, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/products", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductUpdateHandler(t *testing.T) {

	exampleProduct := &models.Product{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
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
	exampleProductUpdateInput := `
		{
			"sku": "example",
			"name": "Test",
			"quantity": 666,
			"upc": "1234567890",
			"price": 12.34
		}
	`

	// exampleUpdatedProduct := &models.Product{
	// 	ID:        exampleProduct.ID,
	// 	CreatedOn: generateExampleTimeForTests(),
	// 	SKU:       "example",
	// 	Name:      "Test",
	// 	UPC:       "1234567890",
	// 	Quantity:  666,
	// 	Cost:      50.00,
	// 	Price:     12.34,
	// }

	t.Run("normal operation", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assertStatusCode(t, testUtil, http.StatusOK)
		ensureExpectationsWereMet(t, testUtil.Mock)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with database error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with input validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/product/example", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with SKU validation error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, "/v1/product/skateboard", strings.NewReader(badSKUUpdateJSON))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with database error updating product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("UpdateProduct", mock.Anything, mock.Anything).Return(generateExampleTimeForTests(), generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product/%s", exampleProduct.SKU), strings.NewReader(exampleProductUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductDeletionHandler(t *testing.T) {
	exampleProduct := &models.Product{
		ID:        2,
		CreatedOn: generateExampleTimeForTests(),
		SKU:       exampleSKU,
		Name:      "Skateboard",
	}

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).Return(time.Now(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, sql.ErrNoRows).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error beginning transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting bridge entries", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).Return(time.Now(), generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error encountered deleting product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).Return(time.Now(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).Return(generateExampleTimeForTests(), generateArbitraryError()).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with existent product", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		testUtil.MockDB.On("GetProductBySKU", mock.Anything, exampleProduct.SKU).Return(exampleProduct, nil).Once()
		testUtil.MockDB.On("DeleteProductVariantBridgeByProductID", mock.Anything, exampleProduct.ID).Return(time.Now(), nil).Once()
		testUtil.MockDB.On("DeleteProduct", mock.Anything, exampleProduct.ID).Return(generateExampleTimeForTests(), nil).Once()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, "/v1/product/example", nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
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
		UPC:           "1234567890",
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
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        exampleProductOption.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []ProductOptionValue{
			{
				DBRow: DBRow{
					ID:        128, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				DBRow: DBRow{
					ID:        256, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				DBRow: DBRow{
					ID:        512, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

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

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithErrorValidatingInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
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
		UPC:           "1234567890",
		Quantity:      123,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
	testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		UPC:           "1234567890",
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
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		UPC:           "1234567890",
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

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductCreationHandlerWithInvalidProductInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(badSKUUpdateJSON))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
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

	assertStatusCode(t, testUtil, http.StatusBadRequest)
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
		UPC:           "1234567890",
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
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        exampleProductOption.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []ProductOptionValue{
			{
				DBRow: DBRow{
					ID:        128, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				DBRow: DBRow{
					ID:        256, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				DBRow: DBRow{
					ID:        512, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

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
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		UPC:           "1234567890",
		Quantity:      123,
	}
	exampleRoot := createProductRootFromProduct(exampleProduct)

	setExpectationsForProductRootSKUExistence(testUtil.Mock, exampleProduct.SKU, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	setExpectationsForProductCreation(testUtil.Mock, exampleProduct, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		UPC:           "1234567890",
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
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        exampleProductOption.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []ProductOptionValue{
			{
				DBRow: DBRow{
					ID:        128, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				DBRow: DBRow{
					ID:        256, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				DBRow: DBRow{
					ID:        512, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

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
	setExpectationsForProductCreationFromOptions(testUtil.Mock, expectedCreatedProducts, 1, generateArbitraryError(), false, 0)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
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
		UPC:           "1234567890",
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
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	expectedCreatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        exampleProductOption.ID,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: exampleProductOption.ProductRootID,
		Values: []ProductOptionValue{
			{
				DBRow: DBRow{
					ID:        128, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "one",
			}, {
				DBRow: DBRow{
					ID:        256, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "two",
			}, {
				DBRow: DBRow{
					ID:        512, // == exampleProductOptionValue.ID,
					CreatedOn: generateExampleTimeForTests(),
				},
				ProductOptionID: exampleProductOption.ID,
				Value:           "three",
			},
		},
	}

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
	setExpectationsForProductCreationFromOptions(testUtil.Mock, expectedCreatedProducts, 1, generateArbitraryError(), true, 0)
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(exampleProductCreationInputWithOptions))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}
