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

func setPasswordResetTokenExistenceQueryByUserIDExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenExistenceQueryByUserID)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestPasswordResetTokenForUserIDExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByUserIDExpectation(t, mock, exampleID, true, nil)
		actual, err := client.PasswordResetTokenForUserIDExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByUserIDExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.PasswordResetTokenForUserIDExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByUserIDExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.PasswordResetTokenForUserIDExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenExistenceQueryByTokenExpectation(t *testing.T, mock sqlmock.Sqlmock, token string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenExistenceQueryByToken)

	mock.ExpectQuery(query).
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestPasswordResetTokenWithTokenExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleToken := "deadbeef"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByTokenExpectation(t, mock, exampleToken, true, nil)
		actual, err := client.PasswordResetTokenWithTokenExists(mockDB, exampleToken)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByTokenExpectation(t, mock, exampleToken, true, sql.ErrNoRows)
		actual, err := client.PasswordResetTokenWithTokenExists(mockDB, exampleToken)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryByTokenExpectation(t, mock, exampleToken, true, errors.New("pineapple on pizza"))
		actual, err := client.PasswordResetTokenWithTokenExists(mockDB, exampleToken)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestPasswordResetTokenExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"token",
		"created_on",
		"expires_on",
		"password_reset_on",
	}).AddRow(
		toReturn.ID,
		toReturn.UserID,
		toReturn.Token,
		toReturn.CreatedOn,
		toReturn.ExpiresOn,
		toReturn.PasswordResetOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetPasswordResetToken(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.PasswordResetToken{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetPasswordResetToken(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected passwordresettoken did not match actual passwordresettoken")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.PasswordResetToken, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"token",
		"created_on",
		"expires_on",
		"password_reset_on",
	}).AddRow(
		example.ID,
		example.UserID,
		example.Token,
		example.CreatedOn,
		example.ExpiresOn,
		example.PasswordResetOn,
	).AddRow(
		example.ID,
		example.UserID,
		example.Token,
		example.CreatedOn,
		example.ExpiresOn,
		example.PasswordResetOn,
	).AddRow(
		example.ID,
		example.UserID,
		example.Token,
		example.CreatedOn,
		example.ExpiresOn,
		example.PasswordResetOn,
	).RowError(1, rowErr)

	query, _ := buildPasswordResetTokenListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetPasswordResetTokenList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.PasswordResetToken{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetPasswordResetTokenList(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setPasswordResetTokenListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetPasswordResetTokenList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildPasswordResetTokenListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetPasswordResetTokenList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setPasswordResetTokenListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetPasswordResetTokenList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildPasswordResetTokenCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM password_reset_tokens WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildPasswordResetTokenCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setPasswordResetTokenCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildPasswordResetTokenCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetPasswordResetTokenCount(t *testing.T) {
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
		setPasswordResetTokenCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetPasswordResetTokenCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenCreationQuery)
	tt := buildTestTime(t)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), tt)
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.UserID,
			toCreate.Token,
			toCreate.ExpiresOn,
			toCreate.PasswordResetOn,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreatePasswordResetToken(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.PasswordResetToken{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenCreationQueryExpectation(t, mock, exampleInput, nil)
		expectedCreatedOn := buildTestTime(t)

		actualID, actualCreatedOn, err := client.CreatePasswordResetToken(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expectedCreatedOn, actualCreatedOn, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.UserID,
			toUpdate.Token,
			toUpdate.ExpiresOn,
			toUpdate.PasswordResetOn,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdatePasswordResetTokenByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.PasswordResetToken{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actual, err := client.UpdatePasswordResetToken(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeletePasswordResetTokenByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.DeletePasswordResetToken(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setPasswordResetTokenDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeletePasswordResetToken(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
