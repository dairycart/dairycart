package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

const loginAttemptExistenceQuery = `SELECT EXISTS(SELECT id FROM login_attempts WHERE id = $1 and archived_on IS NULL);`

func (pg *Postgres) LoginAttemptExists(id uint64) (bool, error) {
	var exists string

	err := pg.DB.QueryRow(loginAttemptExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

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

func (pg *Postgres) DeleteLoginAttempt(id uint64, tx *sql.Tx) (t time.Time, err error) {
	if tx != nil {
		err = tx.QueryRow(loginAttemptDeletionQuery, id).Scan(&t)
	} else {
		err = pg.DB.QueryRow(loginAttemptDeletionQuery, id).Scan(&t)
	}
	return
}
