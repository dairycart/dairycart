package postgres

import (
	"testing"

	// internal dependencies
	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setProductReadQueryExpectationBySKU(t *testing.T, mock sqlmock.Sqlmock, sku string, toReturn models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productQueryBySKU)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"option_summary",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
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
		toReturn.ProductRootID,
		toReturn.Name,
		toReturn.Subtitle,
		toReturn.Description,
		toReturn.OptionSummary,
		toReturn.SKU,
		toReturn.UPC,
		toReturn.Manufacturer,
		toReturn.Brand,
		toReturn.Quantity,
		toReturn.Taxable,
		toReturn.Price,
		toReturn.OnSale,
		toReturn.SalePrice,
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
	mock.ExpectQuery(query).WithArgs(sku).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductBySKU(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleSKU := "hello"
	expected := models.Product{SKU: exampleSKU}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectationBySKU(t, mock, exampleSKU, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetProductBySKU(exampleSKU)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected product did not match actual product")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"option_summary",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
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
		toReturn.ProductRootID,
		toReturn.Name,
		toReturn.Subtitle,
		toReturn.Description,
		toReturn.OptionSummary,
		toReturn.SKU,
		toReturn.UPC,
		toReturn.Manufacturer,
		toReturn.Brand,
		toReturn.Quantity,
		toReturn.Taxable,
		toReturn.Price,
		toReturn.OnSale,
		toReturn.SalePrice,
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

func TestGetProductByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := models.Product{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetProduct(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected product did not match actual product")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on", "available_on"}).AddRow(uint64(1), generateExampleTimeForTests(), generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.ProductRootID,
			toCreate.Name,
			toCreate.Subtitle,
			toCreate.Description,
			toCreate.OptionSummary,
			toCreate.SKU,
			toCreate.UPC,
			toCreate.Manufacturer,
			toCreate.Brand,
			toCreate.Quantity,
			toCreate.Taxable,
			toCreate.Price,
			toCreate.OnSale,
			toCreate.SalePrice,
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

func TestCreateProduct(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := models.Product{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, actualAvailableOn, err := client.CreateProduct(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")
		require.Equal(t, expected, actualAvailableOn, "expected availability time did not match actual availability time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.ProductRootID,
			toUpdate.Name,
			toUpdate.Subtitle,
			toUpdate.Description,
			toUpdate.OptionSummary,
			toUpdate.SKU,
			toUpdate.UPC,
			toUpdate.Manufacturer,
			toUpdate.Brand,
			toUpdate.Quantity,
			toUpdate.Taxable,
			toUpdate.Price,
			toUpdate.OnSale,
			toUpdate.SalePrice,
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

func TestUpdateProductByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := models.Product{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateProduct(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setProductDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteProduct(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
