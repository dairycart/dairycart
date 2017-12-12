package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"testing"

	// internal dependencies
	"github.com/dairycart/dairymodels/v1"

	// external dependencies
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductOptionValueForOptionIDExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, optionID uint64, value string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueForOptionIDExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(optionID, value).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductOptionValueForOptionIDExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleOptionID := uint64(1)
	exampleValue := "example"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductOptionValueForOptionIDExistenceQueryExpectation(t, mock, exampleOptionID, exampleValue, true, nil)
		actual, err := client.ProductOptionValueForOptionIDExists(mockDB, exampleOptionID, exampleValue)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setProductOptionValueForOptionIDExistenceQueryExpectation(t, mock, exampleOptionID, exampleValue, true, sql.ErrNoRows)
		actual, err := client.ProductOptionValueForOptionIDExists(mockDB, exampleOptionID, exampleValue)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setProductOptionValueForOptionIDExistenceQueryExpectation(t, mock, exampleOptionID, exampleValue, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductOptionValueForOptionIDExists(mockDB, exampleOptionID, exampleValue)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueDeletionQueryByOptionIDExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueArchiveQueryByOptionID)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestArchiveProductOptionValuesForOption(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueDeletionQueryByOptionIDExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.ArchiveProductOptionValuesForOption(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueForOptionIDReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, optionID uint64, example *models.ProductOptionValue, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"value",
		"created_on",
		"id",
		"archived_on",
		"updated_on",
		"product_option_id",
	}).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).RowError(1, rowErr)

	mock.ExpectQuery(formatQueryForSQLMock(productOptionValueRetrievalQueryByOptionID)).
		WithArgs(optionID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductOptionValuesForOption(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.ProductOptionValue{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueForOptionIDReadQueryExpectation(t, mock, exampleID, example, nil, nil)
		actual, err := client.GetProductOptionValuesForOption(mockDB, exampleID)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setProductOptionValueForOptionIDReadQueryExpectation(t, mock, exampleID, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductOptionValuesForOption(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		mock.ExpectQuery(formatQueryForSQLMock(productOptionValueRetrievalQueryByOptionID)).
			WillReturnRows(exampleRows)
		actual, err := client.GetProductOptionValuesForOption(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setProductOptionValueForOptionIDReadQueryExpectation(t, mock, exampleID, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductOptionValuesForOption(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductOptionValueExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductOptionValueExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductOptionValueExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setProductOptionValueExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductOptionValueExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setProductOptionValueExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductOptionValueExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"value",
		"created_on",
		"id",
		"archived_on",
		"updated_on",
		"product_option_id",
	}).AddRow(
		toReturn.Value,
		toReturn.CreatedOn,
		toReturn.ID,
		toReturn.ArchivedOn,
		toReturn.UpdatedOn,
		toReturn.ProductOptionID,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductOptionValue(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.ProductOptionValue{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProductOptionValue(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected productoptionvalue did not match actual productoptionvalue")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.ProductOptionValue, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"value",
		"created_on",
		"id",
		"archived_on",
		"updated_on",
		"product_option_id",
	}).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).AddRow(
		example.Value,
		example.CreatedOn,
		example.ID,
		example.ArchivedOn,
		example.UpdatedOn,
		example.ProductOptionID,
	).RowError(1, rowErr)

	query, _ := buildProductOptionValueListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductOptionValueList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.ProductOptionValue{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetProductOptionValueList(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setProductOptionValueListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductOptionValueList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductOptionValueListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductOptionValueList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setProductOptionValueListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductOptionValueList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildProductOptionValueCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM product_option_values WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildProductOptionValueCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setProductOptionValueCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildProductOptionValueCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetProductOptionValueCount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()
	expected := uint64(123)
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetProductOptionValueCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueCreationQuery)
	tt := buildTestTime(t)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), tt)
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.Value,
			toCreate.ProductOptionID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductOptionValue(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductOptionValue{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actualID, actualCreationDate, err := client.CreateProductOptionValue(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.Value,
			toUpdate.ProductOptionID,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductOptionValueByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductOptionValue{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actual, err := client.UpdateProductOptionValue(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductOptionValueByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.DeleteProductOptionValue(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductOptionValueDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProductOptionValue(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueWithProductRootIDDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueWithProductRootIDDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestArchiveProductOptionValuesWithProductRootID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueWithProductRootIDDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.ArchiveProductOptionValuesWithProductRootID(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductOptionValueWithProductRootIDDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.ArchiveProductOptionValuesWithProductRootID(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
