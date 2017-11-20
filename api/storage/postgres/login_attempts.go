package postgres

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

const loginAttemptSelectionQuery = `
    SELECT
        id,
        username,
        successful,
        created_on
    FROM
        login_attempts
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *Postgres) GetLoginAttempt(id uint64) (*models.LoginAttempt, error) {
	l := &models.LoginAttempt{}

	err := pg.DB.QueryRow(loginAttemptSelectionQuery, id).Scan(&l.ID, &l.Username, &l.Successful, &l.CreatedOn)

	return l, err
}

const loginattemptCreationQuery = `
    INSERT INTO login_attempts
        (
            username, successful
        )
    VALUES
        (
            $1, $2
        )
    RETURNING
        id, created_on;
`

func (pg *Postgres) CreateLoginAttempt(nu *models.LoginAttempt) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(loginattemptCreationQuery, &nu.Username, &nu.Successful).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const loginAttemptUpdateQuery = `
    UPDATE login_attempts
    SET
        username = $1, 
        successful = $2
    WHERE id = $2
    RETURNING updated_on;
`

func (pg *Postgres) UpdateLoginAttempt(updated *models.LoginAttempt) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(loginAttemptUpdateQuery, &updated.Username, &updated.Successful, &updated.ID).Scan(&t)
	return t, err
}

const loginAttemptDeletionQuery = `
    UPDATE login_attempts
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteLoginAttempt(id uint64) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(loginAttemptDeletionQuery, id).Scan(&t)
	return t, err
}
