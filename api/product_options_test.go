package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductOptionCreationBody = `
		{
			"name": "something",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`
	exampleProductOptionUpdateBody = `
		{
			"name": "something else"
		}
	`
)

func setExpectationsForProductOptionExistenceByID(mock sqlmock.Sqlmock, a *ProductOption, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionExistenceQuery)
	stringID := strconv.Itoa(int(a.ID))
	mock.ExpectQuery(query).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionExistenceByName(mock sqlmock.Sqlmock, a *ProductOption, productID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(productOptionExistenceQueryForProductByName)
	mock.ExpectQuery(query).
		WithArgs(a.Name, productID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionRetrievalQuery(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "name", "product_root_id", "created_on", "updated_on", "archived_on"}).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...)
	query := formatQueryForSQLMock(productOptionRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionListQuery(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "name", "product_root_id", "created_on", "updated_on", "archived_on"}).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...)
	query, _ := buildProductOptionListQuery(exampleProductID, genereateDefaultQueryFilter())
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionListQueryWithoutFilter(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "name", "product_root_id", "created_on", "updated_on", "archived_on"}).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductRootID, generateExampleTimeForTests(), nil, nil}...)
	query, _ := buildProductOptionListQuery(exampleProductID, nil)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionCreation(mock sqlmock.Sqlmock, a *ProductOption, productRootID uint64, err error) {
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(a.ID, generateExampleTimeForTests())
	query, args := buildProductOptionCreationQuery(a, productRootID)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionUpdate(mock sqlmock.Sqlmock, a *ProductOption, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	query, args := buildProductOptionUpdateQuery(a)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductOptionDeletion(mock sqlmock.Sqlmock, id uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productOptionDeletionQuery)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForProductOptionValuesDeletionByOptionID(mock sqlmock.Sqlmock, id uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(productOptionValuesDeletionQueryByOptionID)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestGenerateCartesianProductForOptions(t *testing.T) {
	t.Parallel()

	small := ProductOptionValue{DBRow: DBRow{ID: 1}, Value: "small"}
	medium := ProductOptionValue{DBRow: DBRow{ID: 2}, Value: "medium"}
	large := ProductOptionValue{DBRow: DBRow{ID: 3}, Value: "large"}
	red := ProductOptionValue{DBRow: DBRow{ID: 4}, Value: "red"}
	green := ProductOptionValue{DBRow: DBRow{ID: 5}, Value: "green"}
	blue := ProductOptionValue{DBRow: DBRow{ID: 6}, Value: "blue"}
	xtraLarge := ProductOptionValue{DBRow: DBRow{ID: 7}, Value: "xtra-large"}
	polyester := ProductOptionValue{DBRow: DBRow{ID: 8}, Value: "polyester"}
	cotton := ProductOptionValue{DBRow: DBRow{ID: 9}, Value: "cotton"}

	tt := []struct {
		in       []*ProductOption
		expected []simpleProductOption
		len      int
	}{
		{
			in: []*ProductOption{
				{Name: "Size", Values: []ProductOptionValue{small, medium, large}},
				{Name: "Color", Values: []ProductOptionValue{red, green, blue}},
			},
			expected: []simpleProductOption{
				{IDs: []uint64{small.ID, red.ID}, OptionSummary: "Size: small, Color: red", SKUPostfix: "small_red", OriginalValues: []ProductOptionValue{small, red}},
				{IDs: []uint64{small.ID, green.ID}, OptionSummary: "Size: small, Color: green", SKUPostfix: "small_green", OriginalValues: []ProductOptionValue{small, green}},
				{IDs: []uint64{small.ID, blue.ID}, OptionSummary: "Size: small, Color: blue", SKUPostfix: "small_blue", OriginalValues: []ProductOptionValue{small, blue}},
				{IDs: []uint64{medium.ID, red.ID}, OptionSummary: "Size: medium, Color: red", SKUPostfix: "medium_red", OriginalValues: []ProductOptionValue{medium, red}},
				{IDs: []uint64{medium.ID, green.ID}, OptionSummary: "Size: medium, Color: green", SKUPostfix: "medium_green", OriginalValues: []ProductOptionValue{medium, green}},
				{IDs: []uint64{medium.ID, blue.ID}, OptionSummary: "Size: medium, Color: blue", SKUPostfix: "medium_blue", OriginalValues: []ProductOptionValue{medium, blue}},
				{IDs: []uint64{large.ID, red.ID}, OptionSummary: "Size: large, Color: red", SKUPostfix: "large_red", OriginalValues: []ProductOptionValue{large, red}},
				{IDs: []uint64{large.ID, green.ID}, OptionSummary: "Size: large, Color: green", SKUPostfix: "large_green", OriginalValues: []ProductOptionValue{large, green}},
				{IDs: []uint64{large.ID, blue.ID}, OptionSummary: "Size: large, Color: blue", SKUPostfix: "large_blue", OriginalValues: []ProductOptionValue{large, blue}},
			},
			len: 9,
		},
		{
			// test that name: value pairs can be completely different sizes
			in: []*ProductOption{
				{Name: "Size", Values: []ProductOptionValue{small, medium, large, xtraLarge}},
				{Name: "Color", Values: []ProductOptionValue{red, green, blue}},
				{Name: "Fabric", Values: []ProductOptionValue{polyester, cotton}},
			},
			expected: []simpleProductOption{
				{
					IDs:            []uint64{small.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: red, Fabric: polyester",
					SKUPostfix:     "small_red_polyester",
					OriginalValues: []ProductOptionValue{small, red, polyester},
				},
				{
					IDs:            []uint64{small.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: red, Fabric: cotton",
					SKUPostfix:     "small_red_cotton",
					OriginalValues: []ProductOptionValue{small, red, cotton},
				},
				{
					IDs:            []uint64{small.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: green, Fabric: polyester",
					SKUPostfix:     "small_green_polyester",
					OriginalValues: []ProductOptionValue{small, green, polyester},
				},
				{
					IDs:            []uint64{small.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: green, Fabric: cotton",
					SKUPostfix:     "small_green_cotton",
					OriginalValues: []ProductOptionValue{small, green, cotton},
				},
				{
					IDs:            []uint64{small.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: blue, Fabric: polyester",
					SKUPostfix:     "small_blue_polyester",
					OriginalValues: []ProductOptionValue{small, blue, polyester},
				},
				{
					IDs:            []uint64{small.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: blue, Fabric: cotton",
					SKUPostfix:     "small_blue_cotton",
					OriginalValues: []ProductOptionValue{small, blue, cotton},
				},
				{
					IDs:            []uint64{medium.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: red, Fabric: polyester",
					SKUPostfix:     "medium_red_polyester",
					OriginalValues: []ProductOptionValue{medium, red, polyester},
				},
				{
					IDs:            []uint64{medium.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: red, Fabric: cotton",
					SKUPostfix:     "medium_red_cotton",
					OriginalValues: []ProductOptionValue{medium, red, cotton},
				},
				{
					IDs:            []uint64{medium.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: green, Fabric: polyester",
					SKUPostfix:     "medium_green_polyester",
					OriginalValues: []ProductOptionValue{medium, green, polyester},
				},
				{
					IDs:            []uint64{medium.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: green, Fabric: cotton",
					SKUPostfix:     "medium_green_cotton",
					OriginalValues: []ProductOptionValue{medium, green, cotton},
				},
				{
					IDs:            []uint64{medium.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: blue, Fabric: polyester",
					SKUPostfix:     "medium_blue_polyester",
					OriginalValues: []ProductOptionValue{medium, blue, polyester},
				},
				{
					IDs:            []uint64{medium.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: blue, Fabric: cotton",
					SKUPostfix:     "medium_blue_cotton",
					OriginalValues: []ProductOptionValue{medium, blue, cotton},
				},
				{
					IDs:            []uint64{large.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: red, Fabric: polyester",
					SKUPostfix:     "large_red_polyester",
					OriginalValues: []ProductOptionValue{large, red, polyester},
				},
				{
					IDs:            []uint64{large.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: red, Fabric: cotton",
					SKUPostfix:     "large_red_cotton",
					OriginalValues: []ProductOptionValue{large, red, cotton},
				},
				{
					IDs:            []uint64{large.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: green, Fabric: polyester",
					SKUPostfix:     "large_green_polyester",
					OriginalValues: []ProductOptionValue{large, green, polyester},
				},
				{
					IDs:            []uint64{large.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: green, Fabric: cotton",
					SKUPostfix:     "large_green_cotton",
					OriginalValues: []ProductOptionValue{large, green, cotton},
				},
				{
					IDs:            []uint64{large.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: blue, Fabric: polyester",
					SKUPostfix:     "large_blue_polyester",
					OriginalValues: []ProductOptionValue{large, blue, polyester},
				},
				{
					IDs:            []uint64{large.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: blue, Fabric: cotton",
					SKUPostfix:     "large_blue_cotton",
					OriginalValues: []ProductOptionValue{large, blue, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: red, Fabric: polyester",
					SKUPostfix:     "xtra-large_red_polyester",
					OriginalValues: []ProductOptionValue{xtraLarge, red, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: red, Fabric: cotton",
					SKUPostfix:     "xtra-large_red_cotton",
					OriginalValues: []ProductOptionValue{xtraLarge, red, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: green, Fabric: polyester",
					SKUPostfix:     "xtra-large_green_polyester",
					OriginalValues: []ProductOptionValue{xtraLarge, green, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: green, Fabric: cotton",
					SKUPostfix:     "xtra-large_green_cotton",
					OriginalValues: []ProductOptionValue{xtraLarge, green, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: blue, Fabric: polyester",
					SKUPostfix:     "xtra-large_blue_polyester",
					OriginalValues: []ProductOptionValue{xtraLarge, blue, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: blue, Fabric: cotton",
					SKUPostfix:     "xtra-large_blue_cotton",
					OriginalValues: []ProductOptionValue{xtraLarge, blue, cotton},
				},
			},
			len: 24,
		},
	}

	for _, tc := range tt {
		actual := generateCartesianProductForOptions(tc.in)
		assert.Equal(t, tc.len, len(actual), "there should be %d simpleProductOptions, but we generated %d", tc.len, len(actual))
		assert.Equal(t, tc.expected, actual, "expected output should match actual output")
	}
}

func TestRetrieveProductOptionFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)

	actual, err := retrieveProductOptionFromDB(testUtil.DB, exampleProductOption.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption, actual, "expected and actual product options should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveProductOptionFromDBWithNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, sql.ErrNoRows)

	_, err := retrieveProductOptionFromDB(testUtil.DB, exampleProductOption.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, exampleProductOption, exampleProductID, nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	newOptionID, createdOn, err := createProductOptionInDB(tx, exampleProductOption, exampleProductID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductOption.ID, newOptionID, "Creating a product option should return the created product option ID")
	assert.Equal(t, exampleProductOption.CreatedOn, createdOn, "Creating a product option should return the created product option creation date")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductOptionInput := &ProductOptionCreationInput{
		Name:   "something",
		Values: []string{"one", "two", "three"},
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

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	actual, err := createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, exampleProductID)
	assert.Nil(t, err)
	assert.Equal(t, expectedCreatedProductOption, actual, "output from product option creation should match expectation")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductOptionInput := &ProductOptionCreationInput{
		Name:   "something",
		Values: []string{"one", "two", "three"},
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

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, generateArbitraryError())
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductRootID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateProductOptionAndValuesInDBFromInputWithIssueCreatingOptionValue(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleProductOptionInput := &ProductOptionCreationInput{
		Name:   "something",
		Values: []string{"one", "two", "three"},
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

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], generateArbitraryError())
	testUtil.Mock.ExpectCommit()

	tx, err := testUtil.DB.Begin()
	assert.Nil(t, err)

	_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionInput, expectedCreatedProductOption.ProductRootID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUpdateProductOptionInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	setExpectationsForProductOptionUpdate(testUtil.Mock, expectedCreatedProductOption, nil)

	updatedOn, err := updateProductOptionInDB(testUtil.DB, exampleProductOption)
	assert.Nil(t, err)
	assert.Equal(t, generateExampleTimeForTests(), updatedOn, "updateProductOptionInDB should return the time the option was updated on")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveProductOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.DB.Beginx()
	assert.Nil(t, err)
	setExpectationsForProductOptionDeletion(testUtil.Mock, 1, nil)

	err = archiveProductOption(tx, 1)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestArchiveProductOptionValuesForOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	testUtil.Mock.ExpectBegin()
	tx, err := testUtil.DB.Beginx()
	assert.Nil(t, err)
	setExpectationsForProductOptionValuesDeletionByOptionID(testUtil.Mock, 1, nil)

	err = archiveProductOptionValuesForOption(tx, 1)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionListHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForRowCount(testUtil.Mock, "product_options", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductOptionListQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProductID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusOK)

	expected := &ListResponse{
		Page:  1,
		Limit: 25,
		Count: 3,
	}

	actual := &ListResponse{}
	bodyString := testUtil.Response.Body.String()
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product option pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product option limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product option counts should be equal")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionListHandlerWithErrorRetrievingCount(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForRowCount(testUtil.Mock, "product_options", genereateDefaultQueryFilter(), 3, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProductID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)

	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionListHandlerWithErrorsRetrievingValues(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForRowCount(testUtil.Mock, "product_options", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductOptionListQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProductID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionListHandlerWithDBErrors(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	setExpectationsForRowCount(testUtil.Mock, "product_options", genereateDefaultQueryFilter(), 3, nil)
	setExpectationsForProductOptionListQuery(testUtil.Mock, exampleProductOption, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product", strconv.Itoa(int(exampleProductID)), "options")
	req, err := http.NewRequest(http.MethodGet, productOptionEndpoint, nil)
	assert.Nil(t, err)

	testUtil.Router.ServeHTTP(testUtil.Response, req)
	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerFailureToSetupTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerFailureToCommitTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, sql.ErrNoRows)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[0], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[1], nil)
	setExpectationsForProductOptionValueCreation(testUtil.Mock, &expectedCreatedProductOption.Values[2], nil)
	testUtil.Mock.ExpectCommit()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusCreated)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithNonExistentProduct(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, false, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithInvalidOptionCreationInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWhenOptionWithTheSameNameExists(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, true, nil)

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionCreationHandlerWithProblemsCreatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
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

	productIDString := strconv.Itoa(int(exampleProductOption.ProductRootID))
	setExpectationsForProductRootExistence(testUtil.Mock, productIDString, true, nil)
	setExpectationsForProductOptionExistenceByName(testUtil.Mock, expectedCreatedProductOption, productIDString, false, nil)

	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionCreation(testUtil.Mock, expectedCreatedProductOption, exampleProductID, generateArbitraryError())
	testUtil.Mock.ExpectRollback()

	productOptionEndpoint := buildRoute("v1", "product", productIDString, "options")
	req, err := http.NewRequest(http.MethodPost, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleUpdatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID: exampleProductOption.ID,
		},
		Name:          "something else",
		ProductRootID: exampleProductOption.ProductRootID,
	}

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(testUtil.Mock, exampleUpdatedProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithNonexistentOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, false, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusBadRequest)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithErrorRetrievingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithErrorUpdatingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleUpdatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID: exampleProductOption.ID,
		},
		Name:          "something else",
		ProductRootID: exampleProductOption.ProductRootID,
	}

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(testUtil.Mock, exampleUpdatedProductOption, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionUpdateHandlerWithErrorRetrievingValues(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	exampleUpdatedProductOption := &ProductOption{
		DBRow: DBRow{
			ID: exampleProductOption.ID,
		},
		Name:          "something else",
		ProductRootID: exampleProductOption.ProductRootID,
	}

	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	setExpectationsForProductOptionRetrievalQuery(testUtil.Mock, exampleProductOption, nil)
	setExpectationsForProductOptionUpdate(testUtil.Mock, exampleUpdatedProductOption, nil)
	setExpectationsForProductOptionValueRetrievalByOptionID(testUtil.Mock, exampleProductOption, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodPatch, productOptionEndpoint, strings.NewReader(exampleProductOptionUpdateBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValuesDeletionByOptionID(testUtil.Mock, exampleProductOption.ID, nil)
	setExpectationsForProductOptionDeletion(testUtil.Mock, exampleProductOption.ID, nil)
	testUtil.Mock.ExpectCommit()

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusOK)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandlerForNonexistentOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, false, nil)
	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusNotFound)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandlerWithErrorCreatingTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandlerWithErrorDeletingOptionValues(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValuesDeletionByOptionID(testUtil.Mock, exampleProductOption.ID, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandlerWithErrorDeletingOption(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValuesDeletionByOptionID(testUtil.Mock, exampleProductOption.ID, nil)
	setExpectationsForProductOptionDeletion(testUtil.Mock, exampleProductOption.ID, generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestProductOptionDeletionHandlerWithErrorCommittingTransaction(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)
	exampleProductOption := &ProductOption{
		DBRow: DBRow{
			ID:        123,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:          "something",
		ProductRootID: 2,
	}
	optionIDString := strconv.Itoa(int(exampleProductOption.ID))

	setExpectationsForProductOptionExistenceByID(testUtil.Mock, exampleProductOption, true, nil)
	testUtil.Mock.ExpectBegin()
	setExpectationsForProductOptionValuesDeletionByOptionID(testUtil.Mock, exampleProductOption.ID, nil)
	setExpectationsForProductOptionDeletion(testUtil.Mock, exampleProductOption.ID, nil)
	testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())

	productOptionEndpoint := buildRoute("v1", "product_options", optionIDString)
	req, err := http.NewRequest(http.MethodDelete, productOptionEndpoint, strings.NewReader(exampleProductOptionCreationBody))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assertStatusCode(t, testUtil, http.StatusInternalServerError)
	ensureExpectationsWereMet(t, testUtil.Mock)
}
