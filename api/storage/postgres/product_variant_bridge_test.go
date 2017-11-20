package postgres

import (
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductVariantBridgeReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductVariantBridge, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productVariantBridgeSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"product_option_value_id",
		"created_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.ProductID,
		toReturn.ProductOptionValueID,
		toReturn.CreatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductVariantBridgeByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := &models.ProductVariantBridge{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetProductVariantBridge(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected productvariantbridge did not match actual productvariantbridge")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductVariantBridgeCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductVariantBridge, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productvariantbridgeCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.ProductID,
			toCreate.ProductOptionValueID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductVariantBridge(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductVariantBridge{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreateProductVariantBridge(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductVariantBridgeUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductVariantBridge, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productVariantBridgeUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.ProductID,
			toUpdate.ProductOptionValueID,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductVariantBridgeByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductVariantBridge{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateProductVariantBridge(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductVariantBridgeDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productVariantBridgeDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductVariantBridgeByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProductVariantBridge(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
