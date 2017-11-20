package postgres

import (
	"database/sql"
	"errors"
	"strconv"
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductVariantBridgeExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productVariantBridgeExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductVariantBridgeExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("existing", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.ProductVariantBridgeExists(exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		client := Postgres{DB: mockDB}
		actual, err := client.ProductVariantBridgeExists(exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		client := Postgres{DB: mockDB}
		actual, err := client.ProductVariantBridgeExists(exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

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
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
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
		expected := generateExampleTimeForTests(t)
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
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
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
		expected := generateExampleTimeForTests(t)
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
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
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
		expected := generateExampleTimeForTests(t)
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProductVariantBridge(exampleID, nil)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductVariantBridgeDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProductVariantBridge(exampleID, tx)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
