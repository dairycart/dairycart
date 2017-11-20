package postgres

import (
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductOptionValueReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"product_option_id",
		"value",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.ProductOptionID,
		toReturn.Value,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductOptionValueByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := &models.ProductOptionValue{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetProductOptionValue(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected productoptionvalue did not match actual productoptionvalue")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productoptionvalueCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.ProductOptionID,
			toCreate.Value,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductOptionValue(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductOptionValue{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreateProductOptionValue(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductOptionValue, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.ProductOptionID,
			toUpdate.Value,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductOptionValueByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductOptionValue{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateProductOptionValue(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductOptionValueDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productOptionValueDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductOptionValueByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setProductOptionValueDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProductOptionValue(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
