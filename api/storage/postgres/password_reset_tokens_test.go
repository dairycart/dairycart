package postgres

import (
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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

func TestGetPasswordResetTokenByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := &models.PasswordResetToken{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetPasswordResetToken(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected passwordresettoken did not match actual passwordresettoken")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordresettokenCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
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

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreatePasswordResetToken(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.PasswordResetToken, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
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

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdatePasswordResetToken(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setPasswordResetTokenDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(passwordResetTokenDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeletePasswordResetTokenByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setPasswordResetTokenDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeletePasswordResetToken(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
