package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

const passwordResetTokenExistenceQuery = `SELECT EXISTS(SELECT id FROM password_reset_tokens WHERE id = $1 and archived_on IS NULL);`

func (pg *Postgres) PasswordResetTokenExists(db storage.Querier, id uint64) (bool, error) {
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
        id,
        user_id,
        token,
        created_on,
        expires_on,
        password_reset_on
    FROM
        password_reset_tokens
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *Postgres) GetPasswordResetToken(db storage.Querier, id uint64) (*models.PasswordResetToken, error) {
	p := &models.PasswordResetToken{}

	err := db.QueryRow(passwordResetTokenSelectionQuery, id).Scan(&p.ID, &p.UserID, &p.Token, &p.CreatedOn, &p.ExpiresOn, &p.PasswordResetOn)

	return p, err
}

const passwordresettokenCreationQuery = `
    INSERT INTO password_reset_tokens
        (
            user_id, token, expires_on, password_reset_on
        )
    VALUES
        (
            $1, $2, $3, $4
        )
    RETURNING
        id, created_on;
`

func (pg *Postgres) CreatePasswordResetToken(db storage.Querier, nu *models.PasswordResetToken) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := db.QueryRow(passwordresettokenCreationQuery, &nu.UserID, &nu.Token, &nu.ExpiresOn, &nu.PasswordResetOn).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const passwordResetTokenUpdateQuery = `
    UPDATE password_reset_tokens
    SET
        user_id = $1, 
        token = $2, 
        expires_on = $3, 
        password_reset_on = $4
    WHERE id = $4
    RETURNING updated_on;
`

func (pg *Postgres) UpdatePasswordResetToken(db storage.Querier, updated *models.PasswordResetToken) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(passwordResetTokenUpdateQuery, &updated.UserID, &updated.Token, &updated.ExpiresOn, &updated.PasswordResetOn, &updated.ID).Scan(&t)
	return t, err
}

const passwordResetTokenDeletionQuery = `
    UPDATE password_reset_tokens
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeletePasswordResetToken(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(passwordResetTokenDeletionQuery, id).Scan(&t)
	return t, err
}
