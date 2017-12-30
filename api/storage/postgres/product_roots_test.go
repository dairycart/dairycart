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

func setProductRootWithSKUPrefixExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, skuPrefix string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootWithSKUPrefixExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(skuPrefix).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductRootWithSKUExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleSKUPrefix := "example"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, nil)
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, sql.ErrNoRows)
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductRootExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductRootExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductRootExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductRootExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"available_on",
		"product_length",
		"updated_on",
		"sku_prefix",
		"package_height",
		"product_weight",
		"product_width",
		"quantity_per_package",
		"name",
		"product_height",
		"package_length",
		"created_on",
		"cost",
		"brand",
		"subtitle",
		"package_weight",
		"archived_on",
		"id",
		"package_width",
		"description",
		"manufacturer",
		"taxable",
	}).AddRow(
		toReturn.AvailableOn,
		toReturn.ProductLength,
		toReturn.UpdatedOn,
		toReturn.SKUPrefix,
		toReturn.PackageHeight,
		toReturn.ProductWeight,
		toReturn.ProductWidth,
		toReturn.QuantityPerPackage,
		toReturn.Name,
		toReturn.ProductHeight,
		toReturn.PackageLength,
		toReturn.CreatedOn,
		toReturn.Cost,
		toReturn.Brand,
		toReturn.Subtitle,
		toReturn.PackageWeight,
		toReturn.ArchivedOn,
		toReturn.ID,
		toReturn.PackageWidth,
		toReturn.Description,
		toReturn.Manufacturer,
		toReturn.Taxable,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductRoot(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.ProductRoot{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProductRoot(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected productroot did not match actual productroot")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.ProductRoot, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"available_on",
		"product_length",
		"updated_on",
		"sku_prefix",
		"package_height",
		"product_weight",
		"product_width",
		"quantity_per_package",
		"name",
		"product_height",
		"package_length",
		"created_on",
		"cost",
		"brand",
		"subtitle",
		"package_weight",
		"archived_on",
		"id",
		"package_width",
		"description",
		"manufacturer",
		"taxable",
	}).AddRow(
		example.AvailableOn,
		example.ProductLength,
		example.UpdatedOn,
		example.SKUPrefix,
		example.PackageHeight,
		example.ProductWeight,
		example.ProductWidth,
		example.QuantityPerPackage,
		example.Name,
		example.ProductHeight,
		example.PackageLength,
		example.CreatedOn,
		example.Cost,
		example.Brand,
		example.Subtitle,
		example.PackageWeight,
		example.ArchivedOn,
		example.ID,
		example.PackageWidth,
		example.Description,
		example.Manufacturer,
		example.Taxable,
	).AddRow(
		example.AvailableOn,
		example.ProductLength,
		example.UpdatedOn,
		example.SKUPrefix,
		example.PackageHeight,
		example.ProductWeight,
		example.ProductWidth,
		example.QuantityPerPackage,
		example.Name,
		example.ProductHeight,
		example.PackageLength,
		example.CreatedOn,
		example.Cost,
		example.Brand,
		example.Subtitle,
		example.PackageWeight,
		example.ArchivedOn,
		example.ID,
		example.PackageWidth,
		example.Description,
		example.Manufacturer,
		example.Taxable,
	).AddRow(
		example.AvailableOn,
		example.ProductLength,
		example.UpdatedOn,
		example.SKUPrefix,
		example.PackageHeight,
		example.ProductWeight,
		example.ProductWidth,
		example.QuantityPerPackage,
		example.Name,
		example.ProductHeight,
		example.PackageLength,
		example.CreatedOn,
		example.Cost,
		example.Brand,
		example.Subtitle,
		example.PackageWeight,
		example.ArchivedOn,
		example.ID,
		example.PackageWidth,
		example.Description,
		example.Manufacturer,
		example.Taxable,
	).RowError(1, rowErr)

	query, _ := buildProductRootListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductRootList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.ProductRoot{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductRootListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductRootList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildProductRootCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM product_roots WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildProductRootCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setProductRootCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildProductRootCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetProductRootCount(t *testing.T) {
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
		setProductRootCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetProductRootCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootCreationQuery)
	tt := buildTestTime(t)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), tt)
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.AvailableOn,
			toCreate.ProductLength,
			toCreate.SKUPrefix,
			toCreate.PackageHeight,
			toCreate.ProductWeight,
			toCreate.ProductWidth,
			toCreate.QuantityPerPackage,
			toCreate.Name,
			toCreate.ProductHeight,
			toCreate.PackageLength,
			toCreate.Cost,
			toCreate.Brand,
			toCreate.Subtitle,
			toCreate.PackageWeight,
			toCreate.PackageWidth,
			toCreate.Description,
			toCreate.Manufacturer,
			toCreate.Taxable,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductRoot(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductRoot{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootCreationQueryExpectation(t, mock, exampleInput, nil)
		expectedCreatedOn := buildTestTime(t)

		actualID, actualCreatedOn, err := client.CreateProductRoot(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expectedCreatedOn, actualCreatedOn, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.AvailableOn,
			toUpdate.ProductLength,
			toUpdate.SKUPrefix,
			toUpdate.PackageHeight,
			toUpdate.ProductWeight,
			toUpdate.ProductWidth,
			toUpdate.QuantityPerPackage,
			toUpdate.Name,
			toUpdate.ProductHeight,
			toUpdate.PackageLength,
			toUpdate.Cost,
			toUpdate.Brand,
			toUpdate.Subtitle,
			toUpdate.PackageWeight,
			toUpdate.PackageWidth,
			toUpdate.Description,
			toUpdate.Manufacturer,
			toUpdate.Taxable,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductRoot{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actual, err := client.UpdateProductRoot(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.DeleteProductRoot(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductRootDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProductRoot(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
