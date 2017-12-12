package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const webhookExecutionLogExistenceQuery = `SELECT EXISTS(SELECT id FROM webhook_execution_logs WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) WebhookExecutionLogExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(webhookExecutionLogExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const webhookExecutionLogSelectionQuery = `
    SELECT
        id,
        succeeded,
        status_code,
        executed_on,
        webhook_id
    FROM
        webhook_execution_logs
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetWebhookExecutionLog(db storage.Querier, id uint64) (*models.WebhookExecutionLog, error) {
	w := &models.WebhookExecutionLog{}

	err := db.QueryRow(webhookExecutionLogSelectionQuery, id).Scan(&w.ID, &w.Succeeded, &w.StatusCode, &w.ExecutedOn, &w.WebhookID)

	return w, err
}

func buildWebhookExecutionLogListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"succeeded",
			"status_code",
			"executed_on",
			"webhook_id",
		).
		From("webhook_execution_logs")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetWebhookExecutionLogList(db storage.Querier, qf *models.QueryFilter) ([]models.WebhookExecutionLog, error) {
	var list []models.WebhookExecutionLog
	query, args := buildWebhookExecutionLogListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var w models.WebhookExecutionLog
		err := rows.Scan(
			&w.ID,
			&w.Succeeded,
			&w.StatusCode,
			&w.ExecutedOn,
			&w.WebhookID,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildWebhookExecutionLogCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("webhook_execution_logs")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetWebhookExecutionLogCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildWebhookExecutionLogCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const webhookExecutionLogCreationQuery = `
    INSERT INTO webhook_execution_logs
        (
            succeeded, status_code, executed_on, webhook_id
        )
    VALUES
        (
            $1, $2, $3, $4
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateWebhookExecutionLog(db storage.Querier, nu *models.WebhookExecutionLog) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(webhookExecutionLogCreationQuery, &nu.Succeeded, &nu.StatusCode, &nu.ExecutedOn, &nu.WebhookID).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const webhookExecutionLogUpdateQuery = `
    UPDATE webhook_execution_logs
    SET
        succeeded = $1, 
        status_code = $2, 
        executed_on = $3, 
        webhook_id = $4
    WHERE id = $4
    RETURNING updated_on;
`

func (pg *postgres) UpdateWebhookExecutionLog(db storage.Querier, updated *models.WebhookExecutionLog) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(webhookExecutionLogUpdateQuery, &updated.Succeeded, &updated.StatusCode, &updated.ExecutedOn, &updated.WebhookID, &updated.ID).Scan(&t)
	return t, err
}

const webhookExecutionLogDeletionQuery = `
    UPDATE webhook_execution_logs
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteWebhookExecutionLog(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(webhookExecutionLogDeletionQuery, id).Scan(&t)
	return t, err
}
