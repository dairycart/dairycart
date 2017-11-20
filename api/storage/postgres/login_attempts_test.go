package postgres

import (
	"testing"

	// internal dependencies
	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setLoginAttemptReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn models.LoginAttempt, err error) {
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

func TestGetLoginAttemptByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := models.LoginAttempt{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetLoginAttempt(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected loginattempt did not match actual loginattempt")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginattemptCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
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
	exampleInput := models.LoginAttempt{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreateLoginAttempt(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate models.LoginAttempt, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
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
	exampleInput := models.LoginAttempt{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateLoginAttempt(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setLoginAttemptDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(loginAttemptDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteLoginAttemptByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setLoginAttemptDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteLoginAttempt(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
