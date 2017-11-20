package postgres

import (
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductRootReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"subtitle",
		"description",
		"sku_prefix",
		"manufacturer",
		"brand",
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
		"quantity_per_package",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.Name,
		toReturn.Subtitle,
		toReturn.Description,
		toReturn.SkuPrefix,
		toReturn.Manufacturer,
		toReturn.Brand,
		toReturn.Taxable,
		toReturn.Cost,
		toReturn.ProductWeight,
		toReturn.ProductHeight,
		toReturn.ProductWidth,
		toReturn.ProductLength,
		toReturn.PackageWeight,
		toReturn.PackageHeight,
		toReturn.PackageWidth,
		toReturn.PackageLength,
		toReturn.QuantityPerPackage,
		toReturn.AvailableOn,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := &models.ProductRoot{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetProductRoot(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected productroot did not match actual productroot")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productrootCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.Name,
			toCreate.Subtitle,
			toCreate.Description,
			toCreate.SkuPrefix,
			toCreate.Manufacturer,
			toCreate.Brand,
			toCreate.Taxable,
			toCreate.Cost,
			toCreate.ProductWeight,
			toCreate.ProductHeight,
			toCreate.ProductWidth,
			toCreate.ProductLength,
			toCreate.PackageWeight,
			toCreate.PackageHeight,
			toCreate.PackageWidth,
			toCreate.PackageLength,
			toCreate.QuantityPerPackage,
			toCreate.AvailableOn,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductRoot(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductRoot{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreateProductRoot(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.Name,
			toUpdate.Subtitle,
			toUpdate.Description,
			toUpdate.SkuPrefix,
			toUpdate.Manufacturer,
			toUpdate.Brand,
			toUpdate.Taxable,
			toUpdate.Cost,
			toUpdate.ProductWeight,
			toUpdate.ProductHeight,
			toUpdate.ProductWidth,
			toUpdate.ProductLength,
			toUpdate.PackageWeight,
			toUpdate.PackageHeight,
			toUpdate.PackageWidth,
			toUpdate.PackageLength,
			toUpdate.QuantityPerPackage,
			toUpdate.AvailableOn,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductRoot{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateProductRoot(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProductRoot(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
