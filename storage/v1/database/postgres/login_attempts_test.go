package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/models/v1"

	// external dependencies
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setExpectationsForLoginAttemptExhaustionQuery(mock sqlmock.Sqlmock, username string, exhausted bool, shouldError bool) {
	count := 0
	if exhausted {
		count = 666
	}

	var argToReturn interface{} = count
	if shouldError {
		argToReturn = "hello"
	}

	exampleRows := sqlmock.NewRows([]string{""}).AddRow(argToReturn)
	query := formatQueryForSQLMock(loginAttemptExhaustionQuery)
	mock.ExpectQuery(query).
		WithArgs(username).
		WillReturnRows(exampleRows)
}

func TestLoginAttemptsHaveBeenExhausted(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()
	exampleUsername := "username"

	t.Run("optimal behavior", func(*testing.T) {
		expected := false
		setExpectationsForLoginAttemptExhaustionQuery(mock, exampleUsername, expected, false)
		actual, err := client.LoginAttemptsHaveBeenExhausted(mockDB, exampleUsername)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("with db error", func(*testing.T) {
		expected := false
		setExpectationsForLoginAttemptExhaustionQuery(mock, exampleUsername, expected, true)
		actual, err := client.LoginAttemptsHaveBeenExhausted(mockDB, exampleUsername)

		assert.NotNil(t, err)
		assert.Equal(t, expected, actual)
	})
}

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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setLoginAttemptExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.LoginAttemptExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.LoginAttempt{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetLoginAttempt(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected loginattempt did not match actual loginattempt")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	assert.NoError(t, err)
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

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setLoginAttemptListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildLoginAttemptListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setLoginAttemptListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetLoginAttemptList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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

	assert.Equal(t, expected, actual, "expected and actual queries should match")
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
	assert.NoError(t, err)
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

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptCreationQuery)
	tt := buildTestTime(t)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), tt)
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
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.LoginAttempt{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptCreationQueryExpectation(t, mock, exampleInput, nil)
		expectedCreatedOn := buildTestTime(t)

		actualID, actualCreatedOn, err := client.CreateLoginAttempt(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expectedCreatedOn, actualCreatedOn, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(buildTestTime(t))
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.LoginAttempt{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actual, err := client.UpdateLoginAttempt(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteLoginAttemptByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.DeleteLoginAttempt(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setLoginAttemptDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteLoginAttempt(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
