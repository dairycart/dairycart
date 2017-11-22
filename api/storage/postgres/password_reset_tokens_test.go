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
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := Postgres{}

	t.Run("existing", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setPasswordResetTokenExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.PasswordResetTokenExists(mockDB, exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.PasswordResetToken{ID: exampleID}
	client := Postgres{}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetPasswordResetToken(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected passwordresettoken did not match actual passwordresettoken")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordresettokenCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
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
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.PasswordResetToken{ID: expectedID}
	client := Postgres{}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreatePasswordResetToken(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
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
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.PasswordResetToken{ID: uint64(1)}
	client := Postgres{}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdatePasswordResetToken(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeletePasswordResetTokenByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := Postgres{}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeletePasswordResetToken(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setPasswordResetTokenDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeletePasswordResetToken(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
