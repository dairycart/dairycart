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

func setDiscountReadQueryExpectationByCode(t *testing.T, mock sqlmock.Sqlmock, code string, toReturn *models.Discount, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountQueryByCode)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"discount_type",
		"amount",
		"starts_on",
		"expires_on",
		"requires_code",
		"code",
		"limited_use",
		"number_of_uses",
		"login_required",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.Name,
		toReturn.DiscountType,
		toReturn.Amount,
		toReturn.StartsOn,
		toReturn.ExpiresOn,
		toReturn.RequiresCode,
		toReturn.Code,
		toReturn.LimitedUse,
		toReturn.NumberOfUses,
		toReturn.LoginRequired,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(code).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetDiscountByCode(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	exampleCode := "welcome"
	expected := &models.Discount{Code: exampleCode}

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountReadQueryExpectationByCode(t, mock, exampleCode, expected, nil)
		actual, err := client.GetDiscountByCode(mockDB, exampleCode)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected discount did not match actual discount")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestDiscountExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setDiscountExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.DiscountExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setDiscountExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.DiscountExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setDiscountExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.DiscountExists(mockDB, exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.Discount, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"discount_type",
		"amount",
		"starts_on",
		"expires_on",
		"requires_code",
		"code",
		"limited_use",
		"number_of_uses",
		"login_required",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.Name,
		toReturn.DiscountType,
		toReturn.Amount,
		toReturn.StartsOn,
		toReturn.ExpiresOn,
		toReturn.RequiresCode,
		toReturn.Code,
		toReturn.LimitedUse,
		toReturn.NumberOfUses,
		toReturn.LoginRequired,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetDiscount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.Discount{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetDiscount(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected discount did not match actual discount")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.Discount, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"discount_type",
		"amount",
		"starts_on",
		"expires_on",
		"requires_code",
		"code",
		"limited_use",
		"number_of_uses",
		"login_required",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		example.ID,
		example.Name,
		example.DiscountType,
		example.Amount,
		example.StartsOn,
		example.ExpiresOn,
		example.RequiresCode,
		example.Code,
		example.LimitedUse,
		example.NumberOfUses,
		example.LoginRequired,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.Name,
		example.DiscountType,
		example.Amount,
		example.StartsOn,
		example.ExpiresOn,
		example.RequiresCode,
		example.Code,
		example.LimitedUse,
		example.NumberOfUses,
		example.LoginRequired,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.Name,
		example.DiscountType,
		example.Amount,
		example.StartsOn,
		example.ExpiresOn,
		example.RequiresCode,
		example.Code,
		example.LimitedUse,
		example.NumberOfUses,
		example.LoginRequired,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildDiscountListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetDiscountList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.Discount{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetDiscountList(mockDB, exampleQF)

		require.Nil(t, err)
		require.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setDiscountListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetDiscountList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildDiscountListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetDiscountList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setDiscountListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetDiscountList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildDiscountCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM discounts WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildDiscountCountRetrievalQuery(exampleQF)

	require.Equal(t, expected, actual, "expected and actual queries should match")
}

func setDiscountCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildDiscountCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetDiscountCount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	client := NewPostgres()
	expected := uint64(123)
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetDiscountCount(mockDB, exampleQF)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "count retrieval method should return the expected value")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.Discount, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.Name,
			toCreate.DiscountType,
			toCreate.Amount,
			toCreate.StartsOn,
			toCreate.ExpiresOn,
			toCreate.RequiresCode,
			toCreate.Code,
			toCreate.LimitedUse,
			toCreate.NumberOfUses,
			toCreate.LoginRequired,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateDiscount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.Discount{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreateDiscount(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.Discount, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.Name,
			toUpdate.DiscountType,
			toUpdate.Amount,
			toUpdate.StartsOn,
			toUpdate.ExpiresOn,
			toUpdate.RequiresCode,
			toUpdate.Code,
			toUpdate.LimitedUse,
			toUpdate.NumberOfUses,
			toUpdate.LoginRequired,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateDiscountByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.Discount{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateDiscount(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setDiscountDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(discountDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteDiscountByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setDiscountDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteDiscount(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setDiscountDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteDiscount(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
