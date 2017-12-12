package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const passwordResetTokenExistenceQueryByUserID = `SELECT EXISTS(SELECT id FROM password_reset_tokens WHERE user_id = $1 AND NOW() < expires_on);`

func (pg *postgres) PasswordResetTokenForUserIDExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(passwordResetTokenExistenceQueryByUserID, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const passwordResetTokenExistenceQueryByToken = `SELECT EXISTS(SELECT id FROM password_reset_tokens WHERE token = $1 AND NOW() < expires_on);`

func (pg *postgres) PasswordResetTokenWithTokenExists(db storage.Querier, token string) (bool, error) {
	var exists string

	err := db.QueryRow(passwordResetTokenExistenceQueryByToken, token).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const passwordResetTokenExistenceQuery = `SELECT EXISTS(SELECT id FROM password_reset_tokens WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) PasswordResetTokenExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(passwordResetTokenExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const passwordResetTokenSelectionQuery = `
    SELECT
        password_reset_on,
        token,
        id,
        created_on,
        user_id,
        expires_on
    FROM
        password_reset_tokens
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetPasswordResetToken(db storage.Querier, id uint64) (*models.PasswordResetToken, error) {
	p := &models.PasswordResetToken{}

	err := db.QueryRow(passwordResetTokenSelectionQuery, id).Scan(&p.PasswordResetOn, &p.Token, &p.ID, &p.CreatedOn, &p.UserID, &p.ExpiresOn)

	return p, err
}

func buildPasswordResetTokenListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"password_reset_on",
			"token",
			"id",
			"created_on",
			"user_id",
			"expires_on",
		).
		From("password_reset_tokens")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetPasswordResetTokenList(db storage.Querier, qf *models.QueryFilter) ([]models.PasswordResetToken, error) {
	var list []models.PasswordResetToken
	query, args := buildPasswordResetTokenListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.PasswordResetToken
		err := rows.Scan(
			&p.PasswordResetOn,
			&p.Token,
			&p.ID,
			&p.CreatedOn,
			&p.UserID,
			&p.ExpiresOn,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildPasswordResetTokenCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("password_reset_tokens")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetPasswordResetTokenCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildPasswordResetTokenCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const passwordResetTokenCreationQuery = `
    INSERT INTO password_reset_tokens
        (
            password_reset_on, token, user_id, expires_on
        )
    VALUES
        (
            $1, $2, $3, $4
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreatePasswordResetToken(db storage.Querier, nu *models.PasswordResetToken) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(passwordResetTokenCreationQuery, &nu.PasswordResetOn, &nu.Token, &nu.UserID, &nu.ExpiresOn).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const passwordResetTokenUpdateQuery = `
    UPDATE password_reset_tokens
    SET
        password_reset_on = $1, 
        token = $2, 
        user_id = $3, 
        expires_on = $4
    WHERE id = $4
    RETURNING updated_on;
`

func (pg *postgres) UpdatePasswordResetToken(db storage.Querier, updated *models.PasswordResetToken) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(passwordResetTokenUpdateQuery, &updated.PasswordResetOn, &updated.Token, &updated.UserID, &updated.ExpiresOn, &updated.ID).Scan(&t)
	return t, err
}

const passwordResetTokenDeletionQuery = `
    UPDATE password_reset_tokens
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeletePasswordResetToken(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(passwordResetTokenDeletionQuery, id).Scan(&t)
	return t, err
}
