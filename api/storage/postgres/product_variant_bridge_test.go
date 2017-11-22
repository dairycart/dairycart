package postgres

import (
	"database/sql"
	"database/sql/driver"
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
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductVariantBridgeExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductVariantBridgeExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductVariantBridgeExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductVariantBridgeExists(mockDB, exampleID)

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

func TestGetProductVariantBridge(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.ProductVariantBridge{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProductVariantBridge(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected productvariantbridge did not match actual productvariantbridge")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductVariantBridgeListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.ProductVariantBridge, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"product_option_value_id",
		"created_on",
		"archived_on",
	}).AddRow(
		example.ID,
		example.ProductID,
		example.ProductOptionValueID,
		example.CreatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.ProductID,
		example.ProductOptionValueID,
		example.CreatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.ProductID,
		example.ProductOptionValueID,
		example.CreatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildProductVariantBridgeListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductVariantBridgeList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.ProductVariantBridge{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetProductVariantBridgeList(mockDB, exampleQF)

		require.Nil(t, err)
		require.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setProductVariantBridgeListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductVariantBridgeList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductVariantBridgeListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductVariantBridgeList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setProductVariantBridgeListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductVariantBridgeList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
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
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreateProductVariantBridge(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildMultiProductVariantBridgeCreationQuery(t *testing.T) {
	expectedQuery := `
        INSERT INTO product_variant_bridge
            (
                product_id, product_option_value_id
            )
        VALUES
            (
                ($1, $2),
                ($1, $3),
                ($1, $4),
                ($1, $5),
                ($1, $6)
            )
        RETURNING
            id, created_on;
    `
	expectedValues := []interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
		uint64(4),
		uint64(5),
		uint64(6),
	}

	exampleInput := []uint64{2, 3, 4, 5, 6}
	actualQuery, actualValues := buildMultiProductVariantBridgeCreationQuery(1, exampleInput)

	require.Equal(t, expectedQuery, actualQuery)
	require.Equal(t, actualValues, expectedValues)
}

func setMultipleProductVariantBridgeCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, productID uint64, optionValueIDs []uint64, err error) {
	t.Helper()
	query, args := buildMultiProductVariantBridgeCreationQuery(productID, optionValueIDs)
	queryToExpect := formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	mock.ExpectExec(queryToExpect).
		WithArgs(argsToExpect...).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func TestCreateMultipleProductVariantBridgesForProductID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		exampleProductID := uint64(1)
		exampleOptionValueIDs := []uint64{2, 3, 4, 5, 6}
		setMultipleProductVariantBridgeCreationQueryExpectation(t, mock, exampleProductID, exampleOptionValueIDs, nil)
		err := client.CreateMultipleProductVariantBridgesForProductID(mockDB, exampleProductID, exampleOptionValueIDs)

		require.Nil(t, err)
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
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateProductVariantBridge(mockDB, exampleInput)

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
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteProductVariantBridge(mockDB, exampleID)

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
		actual, err := client.DeleteProductVariantBridge(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductVariantBridgeDeletionByProductIDQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productVariantBridgeDeletionQueryByProductID)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductVariantBridgeByProductID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductVariantBridgeDeletionByProductIDQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteProductVariantBridgeByProductID(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductVariantBridgeDeletionByProductIDQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProductVariantBridgeByProductID(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
