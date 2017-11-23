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

func setLoginAttemptExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestLoginAttemptExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"username",
		"successful",
		"created_on",
	}).AddRow(
		toReturn.ID,
		toReturn.Username,
		toReturn.Successful,
		toReturn.CreatedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetLoginAttempt(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.LoginAttempt{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetLoginAttempt(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected loginattempt did not match actual loginattempt")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.LoginAttempt, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"username",
		"successful",
		"created_on",
	}).AddRow(
		example.ID,
		example.Username,
		example.Successful,
		example.CreatedOn,
	).AddRow(
		example.ID,
		example.Username,
		example.Successful,
		example.CreatedOn,
	).AddRow(
		example.ID,
		example.Username,
		example.Successful,
		example.CreatedOn,
	).RowError(1, rowErr)

	query, _ := buildLoginAttemptListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetLoginAttemptList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.LoginAttempt{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		require.Nil(t, err)
		require.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setLoginAttemptListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildLoginAttemptListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setLoginAttemptListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildLoginAttemptCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM login_attempts WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildLoginAttemptCountRetrievalQuery(exampleQF)

	require.Equal(t, expected, actual, "expected and actual queries should match")
}

func setLoginAttemptCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildLoginAttemptCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetLoginAttemptCount(t *testing.T) {
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
		setLoginAttemptCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetLoginAttemptCount(mockDB, exampleQF)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "count retrieval method should return the expected value")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginattemptCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.Username,
			toCreate.Successful,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateLoginAttempt(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.LoginAttempt{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreateLoginAttempt(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.Username,
			toUpdate.Successful,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateLoginAttemptByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.LoginAttempt{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateLoginAttempt(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteLoginAttemptByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteLoginAttempt(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setLoginAttemptDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteLoginAttempt(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
