// +build !migrated

package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func createExampleHeadersAndDataFromProductRoot(r *models.ProductRoot) ([]string, []driver.Value) {
	var headers []string
	var values []driver.Value

	productMap := map[string]driver.Value{
		"id":                   r.ID,
		"name":                 r.Name,
		"subtitle":             r.Subtitle,
		"description":          r.Description,
		"sku_prefix":           r.SKUPrefix,
		"manufacturer":         r.Manufacturer,
		"brand":                r.Brand,
		"taxable":              r.Taxable,
		"cost":                 r.Cost,
		"product_weight":       r.ProductWeight,
		"product_height":       r.ProductHeight,
		"product_width":        r.ProductWidth,
		"product_length":       r.ProductLength,
		"package_weight":       r.PackageWeight,
		"package_height":       r.PackageHeight,
		"package_width":        r.PackageWidth,
		"package_length":       r.PackageLength,
		"quantity_per_package": r.QuantityPerPackage,
		"available_on":         r.AvailableOn,
		"created_on":           r.CreatedOn,
		"updated_on":           r.UpdatedOn,
		"archived_on":          r.ArchivedOn,
	}

	for header, value := range productMap {
		headers = append(headers, header)
		values = append(values, value)
	}

	return headers, values
}

func setExpectationsForProductRootSKUExistence(mock sqlmock.Sqlmock, SKU string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(productRootSkuExistenceQuery)).
		WithArgs(SKU).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRootExistence(mock sqlmock.Sqlmock, id string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productRootExistenceQuery)
	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRootRetrieval(mock sqlmock.Sqlmock, r *models.ProductRoot, err error) {
	productRootHeaders, exampleProductRootData := createExampleHeadersAndDataFromProductRoot(r)

	exampleRows := sqlmock.NewRows(productRootHeaders).AddRow(exampleProductRootData...)
	query := formatQueryForSQLMock(productRootRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(r.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRootListQuery(mock sqlmock.Sqlmock, r *models.ProductRoot, err error) {
	productRootHeaders, exampleProductRootData := createExampleHeadersAndDataFromProductRoot(r)

	exampleRows := sqlmock.NewRows(productRootHeaders).
		AddRow(exampleProductRootData...).
		AddRow(exampleProductRootData...).
		AddRow(exampleProductRootData...)

	rootsRetrievalQuery, _ := buildProductRootListQuery(genereateDefaultQueryFilter())
	mock.ExpectQuery(formatQueryForSQLMock(rootsRetrievalQuery)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRootCreation(mock sqlmock.Sqlmock, r *models.ProductRoot, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(r.ID, generateExampleTimeForTests())
	query, args := buildProductRootCreationQuery(r)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(argsToDriverValues(args)...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductRootDeletion(mock sqlmock.Sqlmock, rootID uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productRootDeletionQuery)).
		WithArgs(rootID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForProductsAssociatedWithRootDeletion(mock sqlmock.Sqlmock, rootID uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productDeletionQueryByRootID)).
		WithArgs(rootID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestDeleteProductsAssociatedWithRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleID := uint64(666)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.PlainDB.Begin()
	assert.Nil(t, err)
	setExpectationsForProductsAssociatedWithRootDeletion(testUtil.Mock, exampleID, nil)

	err = deleteProductsAssociatedWithRoot(tx, exampleID)
	assert.Nil(t, err)
}

func setExpectationsForProductOptionsAssociatedWithRootDeletion(mock sqlmock.Sqlmock, rootID uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productOptionDeletionQueryByRootID)).
		WithArgs(rootID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestDeleteProductOptionsAssociatedWithRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleID := uint64(666)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.PlainDB.Begin()
	assert.Nil(t, err)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleID, nil)

	err = deleteProductOptionsAssociatedWithRoot(tx, exampleID)
	assert.Nil(t, err)
}

func setExpectationsForProductOptionValuesAssociatedWithRootDeletion(mock sqlmock.Sqlmock, rootID uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productOptionValueDeletionQueryByRootID)).
		WithArgs(rootID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestDeleteProductOptionValuesAssociatedWithRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleID := uint64(666)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.PlainDB.Begin()
	assert.Nil(t, err)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleID, nil)

	err = deleteProductOptionValuesAssociatedWithRoot(tx, exampleID)
	assert.Nil(t, err)
}

func setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(mock sqlmock.Sqlmock, rootID uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productVariantBridgeDeletionQueryByRootID)).
		WithArgs(rootID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestDeleteVariantBridgeEntriesAssociatedWithRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleID := uint64(666)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.PlainDB.Begin()
	assert.Nil(t, err)
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleID, nil)

	err = deleteVariantBridgeEntriesAssociatedWithRoot(tx, exampleID)
	assert.Nil(t, err)
}

func setupExpectationsForProductRootRetrieval(mock sqlmock.Sqlmock, r *models.ProductRoot, err error) {
	exampleProductRootHeaders, exampleProductRootValues := createExampleHeadersAndDataFromProductRoot(r)
	exampleRows := sqlmock.NewRows(exampleProductRootHeaders).AddRow(exampleProductRootValues...)
	mock.ExpectQuery(formatQueryForSQLMock(productRootRetrievalQuery)).
		WithArgs(r.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func createExampleHeadersAndDataFromProduct(p *models.Product) ([]string, []driver.Value) {
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

func setExpectationsForProductAssociatedWithRootListQuery(mock sqlmock.Sqlmock, r *models.ProductRoot, p *models.Product, err error) {
	productHeaders, exampleProductData := createExampleHeadersAndDataFromProduct(p)
	exampleRows := sqlmock.NewRows(productHeaders).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...).
		AddRow(exampleProductData...)

	query, _ := buildProductAssociatedWithRootListQuery(r.ID)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(r.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestDeleteProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleID := uint64(666)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.PlainDB.Begin()
	assert.Nil(t, err)
	setExpectationsForProductRootDeletion(testUtil.Mock, exampleID, nil)

	err = deleteProductRoot(tx, exampleID)
	assert.Nil(t, err)
}

func TestCreateProductRootFromProduct(t *testing.T) {
	exampleInput := &models.Product{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		Manufacturer:       "mfgr",
		Brand:              "brand",
		QuantityPerPackage: 666,
		Taxable:            true,
		Cost:               12.34,
		ProductWeight:      1,
		ProductHeight:      1,
		ProductWidth:       1,
		ProductLength:      1,
		PackageWeight:      1,
		PackageHeight:      1,
		PackageWidth:       1,
		PackageLength:      1,
		AvailableOn:        generateExampleTimeForTests(),
	}
	expected := &models.ProductRoot{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKUPrefix:          "sku",
		Manufacturer:       "mfgr",
		Brand:              "brand",
		QuantityPerPackage: 666,
		Taxable:            true,
		Cost:               12.34,
		ProductWeight:      1,
		ProductHeight:      1,
		ProductWidth:       1,
		ProductLength:      1,
		PackageWeight:      1,
		PackageHeight:      1,
		PackageWidth:       1,
		PackageLength:      1,
		AvailableOn:        generateExampleTimeForTests(),
	}
	actual := createProductRootFromProduct(exampleInput)

	assert.Equal(t, expected, actual, "expected output should match actual output")
}

func TestCreateProductRootInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleRoot := &models.ProductRoot{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "Skateboard",
		Description:   "This is a skateboard. Please wear a helmet.",
		Cost:          50.00,
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductRootCreation(testUtil.Mock, exampleRoot, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	newRootID, createdOn, err := createProductRootInDB(tx, exampleRoot)
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), newRootID, "createProductRootInDB should return the correct ID for a new root ")
	assert.Equal(t, generateExampleTimeForTests(), createdOn, "createProductRootInDB should return the correct creation time for a new root ")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductRootFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleRoot := models.ProductRoot{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "Skateboard",
		Description:   "This is a skateboard. Please wear a helmet.",
		Cost:          50.00,
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}

	setupExpectationsForProductRootRetrieval(testUtil.Mock, &exampleRoot, nil)

	actual, err := retrieveProductRootFromDB(testUtil.DB, exampleRoot.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleRoot, actual, "product root retrieved by query should match")

	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestSingleProductRootRetrievalHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := &models.Product{
		ID:          2,
		CreatedOn:   generateExampleTimeForTests(),
		Name:        "Skateboard",
		Description: "This is a skateboard. Please wear a helmet.",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, nil)
	setExpectationsForProductOptionListQueryWithoutFilter(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestSingleProductRootRetrievalHandlerWhenNoSuchRootExists(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestSingleProductRootRetrievalHandlerWithErrorQueryingDatabaseForProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestSingleProductRootRetrievalHandlerWithErrorRetrievingAssociatedProducts(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := &models.Product{
		ID:          2,
		CreatedOn:   generateExampleTimeForTests(),
		Name:        "Skateboard",
		Description: "This is a skateboard. Please wear a helmet.",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestSingleProductRootRetrievalHandlerWitherrorRetrievingProductOptions(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := &models.Product{
		ID:          2,
		CreatedOn:   generateExampleTimeForTests(),
		Name:        "Skateboard",
		Description: "This is a skateboard. Please wear a helmet.",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, nil)
	setExpectationsForProductOptionListQueryWithoutFilter(testUtil.Mock, exampleProductOption, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootListRetrievalHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := &models.Product{
		ID:          2,
		CreatedOn:   generateExampleTimeForTests(),
		SKU:         "skateboard",
		Name:        "Skateboard",
		UPC:         "1234567890",
		Quantity:    123,
		Price:       12.34,
		Cost:        5,
		Taxable:     true,
		Description: "This is a skateboard. Please wear a helmet.",
	}

	setExpectationsForRowCount(testUtil.Mock, "product_roots", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductRootListQuery(testUtil.Mock, exampleProductRoot, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootListRetrievalHandlerWithErrorGettingRowCount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForRowCount(testUtil.Mock, "product_roots", genereateDefaultQueryFilter(), 3, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootListRetrievalHandlerWithErrorRetrievingProductRoots(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForRowCount(testUtil.Mock, "product_roots", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductRootListQuery(testUtil.Mock, exampleProductRoot, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootListRetrievalHandlerWithErrorRetrivingAssociatedProducts(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}
	exampleProduct := &models.Product{
		ID:          2,
		CreatedOn:   generateExampleTimeForTests(),
		SKU:         "skateboard",
		Name:        "Skateboard",
		UPC:         "1234567890",
		Quantity:    123,
		Price:       12.34,
		Cost:        5,
		Taxable:     true,
		Description: "This is a skateboard. Please wear a helmet.",
	}

	setExpectationsForRowCount(testUtil.Mock, "product_roots", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductRootListQuery(testUtil.Mock, exampleProductRoot, nil)
	setExpectationsForProductAssociatedWithRootListQuery(testUtil.Mock, exampleProductRoot, exampleProduct, generateArbitraryError())

	req, err := http.NewRequest(http.MethodGet, "/v1/product_roots", nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	testUtil.Mock.ExpectCommit()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithNonexistentProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorRetrievingProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, generateArbitraryError())

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorBeginningTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorDeletingBridgeEntries(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorDeletingOptionValues(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorDeletingOptions(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorDeletingProducts(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorDeletingProductRoot(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductRootDeletion(testUtil.Mock, exampleProductRoot.ID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductRootDeletionHandlerWithErrorCommittingTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductRoot := &models.ProductRoot{
		ID:           2,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "root_name",
		Subtitle:     "subtitle",
		Description:  "description",
		SKUPrefix:    "sku_prefix",
		Manufacturer: "manufacturer",
		Brand:        "brand",
	}

	setExpectationsForProductRootRetrieval(testUtil.Mock, exampleProductRoot, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForVariantBridgeEntriesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionValuesAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductOptionsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductsAssociatedWithRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	setExpectationsForProductRootDeletion(testUtil.Mock, exampleProductRoot.ID, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_root/%d", exampleProductRoot.ID), nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}
