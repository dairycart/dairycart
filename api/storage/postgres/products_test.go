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
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	exampleSKU := "hello"
	expected := &models.Product{SKU: exampleSKU}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectationBySKU(t, mock, exampleSKU, expected, nil)
		actual, err := client.GetProductBySKU(mockDB, exampleSKU)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected product did not match actual product")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
func setProductReadQueryExpectationByProductRootID(t *testing.T, mock sqlmock.Sqlmock, example *models.Product, rowErr error, err error) {
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

	mock.ExpectQuery(formatQueryForSQLMock(productQueryByProductRootID)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductsByProductRootID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	exampleProductRootID := uint64(1)
	example := &models.Product{ProductRootID: exampleProductRootID}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectationByProductRootID(t, mock, example, nil, nil)
		actual, err := client.GetProductsByProductRootID(mockDB, exampleProductRootID)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setProductReadQueryExpectationByProductRootID(t, mock, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductsByProductRootID(mockDB, exampleProductRootID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		mock.ExpectQuery(formatQueryForSQLMock(productQueryByProductRootID)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductsByProductRootID(mockDB, exampleProductRootID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setProductReadQueryExpectationByProductRootID(t, mock, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductsByProductRootID(mockDB, exampleProductRootID)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleSKU := "example"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, nil)
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, sql.ErrNoRows)
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductWithSKUExistenceQueryExpectation(t, mock, exampleSKU, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductWithSKUExists(mockDB, exampleSKU)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.Product{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProduct(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected product did not match actual product")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
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

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setProductListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setProductListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildProductCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM products WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildProductCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setProductCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildProductCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetProductCount(t *testing.T) {
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
		setProductCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetProductCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.Product{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, actualAvailableOn, err := client.CreateProduct(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")
		assert.Equal(t, expected, actualAvailableOn, "expected availability time did not match actual availability time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.Product{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateProduct(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteProduct(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProduct(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductWithProductRootIDDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productWithProductRootIDDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestArchiveProductsWithProductRootID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductWithProductRootIDDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.ArchiveProductsWithProductRootID(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductWithProductRootIDDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.ArchiveProductsWithProductRootID(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
