package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/Masterminds/squirrel"
)

const webhookQueryByEventType = `
    SELECT
        id,
        url,
        event_type,
        content_type,
        created_on,
        updated_on,
        archived_on
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
			&w.ID,
			&w.URL,
			&w.EventType,
			&w.ContentType,
			&w.CreatedOn,
			&w.UpdatedOn,
			&w.ArchivedOn,
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
        id,
        url,
        event_type,
        content_type,
        created_on,
        updated_on,
        archived_on
    FROM
        webhooks
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetWebhook(db storage.Querier, id uint64) (*models.Webhook, error) {
	w := &models.Webhook{}

	err := db.QueryRow(webhookSelectionQuery, id).Scan(&w.ID, &w.URL, &w.EventType, &w.ContentType, &w.CreatedOn, &w.UpdatedOn, &w.ArchivedOn)

	return w, err
}

func buildWebhookListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"url",
			"event_type",
			"content_type",
			"created_on",
			"updated_on",
			"archived_on",
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
			&w.ID,
			&w.URL,
			&w.EventType,
			&w.ContentType,
			&w.CreatedOn,
			&w.UpdatedOn,
			&w.ArchivedOn,
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
            url, event_type, content_type
        )
    VALUES
        (
            $1, $2, $3
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateWebhook(db storage.Querier, nu *models.Webhook) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(webhookCreationQuery, &nu.URL, &nu.EventType, &nu.ContentType).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const webhookUpdateQuery = `
    UPDATE webhooks
    SET
        url = $1, 
        event_type = $2, 
        content_type = $3, 
        updated_on = NOW()
    WHERE id = $4
    RETURNING updated_on;
`

func (pg *postgres) UpdateWebhook(db storage.Querier, updated *models.Webhook) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(webhookUpdateQuery, &updated.URL, &updated.EventType, &updated.ContentType, &updated.ID).Scan(&t)
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
