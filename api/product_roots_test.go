package main

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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

func setExpectationsForProductRootCreation(mock sqlmock.Sqlmock, r *ProductRoot, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(r.ID, generateExampleTimeForTests())
	query, args := buildProductRootCreationQuery(r)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(argsToDriverValues(args)...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setupExpectationsForProductRootRetrieval(mock sqlmock.Sqlmock, r *ProductRoot, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"subtitle",
		"description",
		"sku_prefix",
		"manufacturer",
		"brand",
		"available_on",
		"all_options_populated",
		"taxable",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		r.ID,
		r.Name,
		r.Subtitle.String,
		r.Description,
		r.SKUPrefix,
		r.Manufacturer.String,
		r.Brand.String,
		r.AvailableOn,
		r.AllOptionsPopulated,
		r.Taxable,
		r.Cost,
		r.ProductWeight,
		r.ProductHeight,
		r.ProductWidth,
		r.ProductLength,
		r.PackageWeight,
		r.PackageHeight,
		r.PackageWidth,
		r.PackageLength,
		generateExampleTimeForTests(),
		nil,
		nil,
	)
	mock.ExpectQuery(formatQueryForSQLMock(productRootRetrievalQuery)).
		WithArgs(r.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductRootFromProduct(t *testing.T) {
	exampleInput := &Product{
		Name:               "name",
		Subtitle:           NullString{sql.NullString{String: "subtitle", Valid: true}},
		Description:        "description",
		SKU:                "sku",
		Manufacturer:       NullString{sql.NullString{String: "mfgr", Valid: true}},
		Brand:              NullString{sql.NullString{String: "brand", Valid: true}},
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
	expected := &ProductRoot{
		Name:               "name",
		Subtitle:           NullString{sql.NullString{String: "subtitle", Valid: true}},
		Description:        "description",
		SKUPrefix:          "sku",
		Manufacturer:       NullString{sql.NullString{String: "mfgr", Valid: true}},
		Brand:              NullString{sql.NullString{String: "brand", Valid: true}},
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
	exampleRoot := &ProductRoot{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
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

	exampleRoot := &ProductRoot{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
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
	exampleRoot.Subtitle.Valid = true
	exampleRoot.Manufacturer.Valid = true
	exampleRoot.Brand.Valid = true

	setupExpectationsForProductRootRetrieval(testUtil.Mock, exampleRoot, nil)

	actual, err := retrieveProductRootFromDB(testUtil.DB, exampleRoot.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleRoot, actual, "product root retrieved by query should match")

	ensureExpectationsWereMet(t, testUtil.Mock)
}
