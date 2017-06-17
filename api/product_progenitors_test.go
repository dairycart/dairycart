package api

import (
	"database/sql/driver"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var exampleProgenitor *ProductProgenitor
var productProgenitorHeaders []string
var exampleProgenitorData []driver.Value

func init() {
	exampleProgenitor = &ProductProgenitor{
		ID:            2,
		Name:          "Skateboard",
		Description:   "This is a skateboard. Please wear a helmet.",
		Price:         99.99,
		Cost:          50.00,
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		CreatedOn:     exampleTime,
	}
	productProgenitorHeaders = []string{"id", "name", "description", "taxable", "price", "cost", "product_weight", "product_height", "product_width", "product_length", "package_weight", "package_height", "package_width", "package_length", "created_on", "updated_on", "archived_on"}
	exampleProgenitorData = []driver.Value{2, "Skateboard", "This is a skateboard. Please wear a helmet.", false, 99.99, 50.00, 8, 7, 6, 5, 4, 3, 2, 1, exampleTime, nil, nil}

}

func setExpectationsForProductProgenitorExistence(mock sqlmock.Sqlmock, id string, exists bool) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productProgenitorExistenceQuery)
	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(exampleRows)
}

func setExpectationsForProductProgenitorCreation(mock sqlmock.Sqlmock, g *ProductProgenitor, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProgenitor.ID)
	query, args := buildProgenitorCreationQuery(g)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(argsToDriverValues(args)...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setupExpectationsForProductProgenitorRetrieval(mock sqlmock.Sqlmock) {
	exampleRows := sqlmock.NewRows(productProgenitorHeaders).AddRow(exampleProgenitorData...)
	mock.ExpectQuery(formatQueryForSQLMock(productProgenitorRetrievalQuery)).
		WithArgs(exampleProgenitor.ID).
		WillReturnRows(exampleRows)
}

func TestNewProductProgenitorFromProductCreationInput(t *testing.T) {
	t.Parallel()
	expected := &ProductProgenitor{
		Name:          "Example",
		Description:   "this is a description",
		Taxable:       true,
		Price:         10,
		ProductWeight: 10,
		ProductHeight: 10,
		ProductWidth:  10,
		ProductLength: 10,
		PackageWeight: 10,
		PackageHeight: 10,
		PackageWidth:  10,
		PackageLength: 10,
	}
	input := &ProductCreationInput{
		Name:          "Example",
		Description:   "this is a description",
		Taxable:       true,
		Price:         10,
		ProductWeight: 10,
		ProductHeight: 10,
		ProductWidth:  10,
		ProductLength: 10,
		PackageWeight: 10,
		PackageHeight: 10,
		PackageWidth:  10,
		PackageLength: 10,
	}
	actual := newProductProgenitorFromProductCreationInput(input)
	assert.Equal(t, expected, actual, "Output of newProductProgenitorFromProductCreationInput was unexpected")
}

func TestCreateProductProgenitorInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductProgenitorCreation(mock, exampleProgenitor, nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	newProgenitorID, err := createProductProgenitorInDB(tx, exampleProgenitor)
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), newProgenitorID, "createProductProgenitorInDB should return the correct ID for a new progenitor")

	err = tx.Commit()
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductProgenitorFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setupExpectationsForProductProgenitorRetrieval(mock)

	actual, err := retrieveProductProgenitorFromDB(db, exampleProgenitor.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProgenitor, actual, "product progenitor retrieved by query should match")
}
