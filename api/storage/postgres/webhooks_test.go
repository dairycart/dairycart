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

func setWebhookReadQueryExpectationByEventType(t *testing.T, mock sqlmock.Sqlmock, example *models.Webhook, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"url",
		"event_type",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	mock.ExpectQuery(formatQueryForSQLMock(webhookQueryByEventType)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetWebhooksByEventType(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	exampleEventType := "product_updated"
	example := &models.Webhook{}

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookReadQueryExpectationByEventType(t, mock, example, nil, nil)
		actual, err := client.GetWebhooksByEventType(mockDB, exampleEventType)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setWebhookReadQueryExpectationByEventType(t, mock, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetWebhooksByEventType(mockDB, exampleEventType)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		mock.ExpectQuery(formatQueryForSQLMock(webhookQueryByEventType)).
			WillReturnRows(exampleRows)

		actual, err := client.GetWebhooksByEventType(mockDB, exampleEventType)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setWebhookReadQueryExpectationByEventType(t, mock, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetWebhooksByEventType(mockDB, exampleEventType)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(webhookExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestWebhookExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setWebhookExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.WebhookExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setWebhookExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.WebhookExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setWebhookExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.WebhookExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.Webhook, err error) {
	t.Helper()
	query := formatQueryForSQLMock(webhookSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"url",
		"event_type",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.URL,
		toReturn.EventType,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetWebhook(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.Webhook{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetWebhook(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected webhook did not match actual webhook")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.Webhook, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"url",
		"event_type",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.URL,
		example.EventType,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildWebhookListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetWebhookList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.Webhook{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetWebhookList(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setWebhookListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetWebhookList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildWebhookListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetWebhookList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setWebhookListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetWebhookList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildWebhookCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM webhooks WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildWebhookCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setWebhookCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildWebhookCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetWebhookCount(t *testing.T) {
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
		setWebhookCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetWebhookCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.Webhook, err error) {
	t.Helper()
	query := formatQueryForSQLMock(webhookCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.URL,
			toCreate.EventType,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateWebhook(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.Webhook{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreateWebhook(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.Webhook, err error) {
	t.Helper()
	query := formatQueryForSQLMock(webhookUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.URL,
			toUpdate.EventType,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateWebhookByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.Webhook{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateWebhook(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setWebhookDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(webhookDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteWebhookByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setWebhookDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteWebhook(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setWebhookDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteWebhook(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
