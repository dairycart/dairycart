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

func setProductReadQueryExpectationBySKU(t *testing.T, mock sqlmock.Sqlmock, sku string, toReturn *models.Product, err error) {
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
	client := NewPostgres()

	exampleSKU := "hello"
	expected := &models.Product{SKU: exampleSKU}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectationBySKU(t, mock, exampleSKU, expected, nil)
		actual, err := client.GetProductBySKU(mockDB, exampleSKU)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected product did not match actual product")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductWithSKUExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, sku string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productWithSKUExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(sku).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductWithSKUExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleSKU := "example"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, nil)
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, sql.ErrNoRows)
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductExists(mockDB, exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.Product, err error) {
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

func TestGetProduct(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.Product{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProduct(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected product did not match actual product")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.Product, rowErr error, err error) {
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
		example.ID,
		example.ProductRootID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.OptionSummary,
		example.SKU,
		example.UPC,
		example.Manufacturer,
		example.Brand,
		example.Quantity,
		example.Taxable,
		example.Price,
		example.OnSale,
		example.SalePrice,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.ProductRootID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.OptionSummary,
		example.SKU,
		example.UPC,
		example.Manufacturer,
		example.Brand,
		example.Quantity,
		example.Taxable,
		example.Price,
		example.OnSale,
		example.SalePrice,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.ProductRootID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.OptionSummary,
		example.SKU,
		example.UPC,
		example.Manufacturer,
		example.Brand,
		example.Quantity,
		example.Taxable,
		example.Price,
		example.OnSale,
		example.SalePrice,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildProductListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.Product{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetProductList(mockDB, exampleQF)

		require.Nil(t, err)
		require.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setProductListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setProductListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on", "available_on"}).AddRow(uint64(1), generateExampleTimeForTests(t), generateExampleTimeForTests(t))
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
	exampleInput := &models.Product{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, actualAvailableOn, err := client.CreateProduct(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")
		require.Equal(t, expected, actualAvailableOn, "expected availability time did not match actual availability time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.Product, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
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
	exampleInput := &models.Product{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateProduct(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteProduct(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProduct(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
