package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const webhookQueryByEventType = `
    SELECT
        url,
        created_on,
        updated_on,
        id,
        content_type,
        archived_on,
        event_type
    FROM
        webhooks
    WHERE
        event_type = $1
`

func (pg *postgres) GetWebhooksByEventType(db storage.Querier, eventType string) ([]models.Webhook, error) {
	var list []models.Webhook

	rows, err := db.Query(webhookQueryByEventType, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var w models.Webhook
		err := rows.Scan(
			&w.URL,
			&w.CreatedOn,
			&w.UpdatedOn,
			&w.ID,
			&w.ContentType,
			&w.ArchivedOn,
			&w.EventType,
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

const webhookExistenceQuery = `SELECT EXISTS(SELECT id FROM webhooks WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) WebhookExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(webhookExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const webhookSelectionQuery = `
    SELECT
        url,
        created_on,
        updated_on,
        id,
        content_type,
        archived_on,
        event_type
    FROM
        webhooks
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetWebhook(db storage.Querier, id uint64) (*models.Webhook, error) {
	w := &models.Webhook{}

	err := db.QueryRow(webhookSelectionQuery, id).Scan(&w.URL, &w.CreatedOn, &w.UpdatedOn, &w.ID, &w.ContentType, &w.ArchivedOn, &w.EventType)

	return w, err
}

func buildWebhookListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"url",
			"created_on",
			"updated_on",
			"id",
			"content_type",
			"archived_on",
			"event_type",
		).
		From("webhooks")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetWebhookList(db storage.Querier, qf *models.QueryFilter) ([]models.Webhook, error) {
	var list []models.Webhook
	query, args := buildWebhookListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var w models.Webhook
		err := rows.Scan(
			&w.URL,
			&w.CreatedOn,
			&w.UpdatedOn,
			&w.ID,
			&w.ContentType,
			&w.ArchivedOn,
			&w.EventType,
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

func buildWebhookCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("webhooks")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetWebhookCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildWebhookCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const webhookCreationQuery = `
    INSERT INTO webhooks
        (
            url, content_type, event_type
        )
    VALUES
        (
            $1, $2, $3
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateWebhook(db storage.Querier, nu *models.Webhook) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(webhookCreationQuery, &nu.URL, &nu.ContentType, &nu.EventType).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const webhookUpdateQuery = `
    UPDATE webhooks
    SET
        url = $1,
        content_type = $2,
        event_type = $3,
        updated_on = NOW()
    WHERE id = $4
    RETURNING updated_on;
`

func (pg *postgres) UpdateWebhook(db storage.Querier, updated *models.Webhook) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(webhookUpdateQuery, &updated.URL, &updated.ContentType, &updated.EventType, &updated.ID).Scan(&t)
	return t, err
}

const webhookDeletionQuery = `
    UPDATE webhooks
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteWebhook(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(webhookDeletionQuery, id).Scan(&t)
	return t, err
}
