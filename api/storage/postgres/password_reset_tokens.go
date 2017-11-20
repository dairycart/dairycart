package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

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

func (pg *Postgres) GetPasswordResetToken(id uint64) (models.PasswordResetToken, error) {
	var p models.PasswordResetToken

	err := pg.DB.QueryRow(passwordResetTokenSelectionQuery, id).Scan(&p.ID, &p.UserID, &p.Token, &p.CreatedOn, &p.ExpiresOn, &p.PasswordResetOn)

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

func (pg *Postgres) CreatePasswordResetToken(np models.PasswordResetToken) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(passwordresettokenCreationQuery, &np.UserID, &np.Token, &np.ExpiresOn, &np.PasswordResetOn).Scan(&createdID, &createdAt)

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

func (pg *Postgres) UpdatePasswordResetToken(updated models.PasswordResetToken) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(passwordResetTokenUpdateQuery, &updated.UserID, &updated.Token, &updated.ExpiresOn, &updated.PasswordResetOn, &updated.ID).Scan(&t)
	return t, err
}

const passwordResetTokenDeletionQuery = `
    UPDATE password_reset_tokens
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeletePasswordResetToken(id uint64) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(passwordResetTokenDeletionQuery, id).Scan(&t)
	return t, err
}
