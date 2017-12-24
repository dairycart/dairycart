package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const loginAttemptExhaustionQuery = `
    SELECT count(id) FROM login_attempts
        WHERE username = $1
        AND created_on < NOW()
        AND successful IS false
        AND created_on > (NOW() - (15 * interval '1 minute'))
`

func (pg *postgres) LoginAttemptsHaveBeenExhausted(db storage.Querier, username string) (bool, error) {
	var loginCount uint64
	err := db.QueryRow(loginAttemptExhaustionQuery, username).Scan(&loginCount)
	if err != nil {
		return false, err
	}
	return loginCount >= 10, err
}

const loginAttemptExistenceQuery = `SELECT EXISTS(SELECT id FROM login_attempts WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) LoginAttemptExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(loginAttemptExistenceQuery, id).Scan(&exists)
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

func (pg *postgres) GetLoginAttempt(db storage.Querier, id uint64) (*models.LoginAttempt, error) {
	l := &models.LoginAttempt{}

	err := db.QueryRow(loginAttemptSelectionQuery, id).Scan(&l.ID, &l.Username, &l.Successful, &l.CreatedOn)

	return l, err
}

func buildLoginAttemptListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"username",
			"successful",
			"created_on",
		).
		From("login_attempts")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetLoginAttemptList(db storage.Querier, qf *models.QueryFilter) ([]models.LoginAttempt, error) {
	var list []models.LoginAttempt
	query, args := buildLoginAttemptListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l models.LoginAttempt
		err := rows.Scan(
			&l.ID,
			&l.Username,
			&l.Successful,
			&l.CreatedOn,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildLoginAttemptCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("login_attempts")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetLoginAttemptCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildLoginAttemptCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const loginAttemptCreationQuery = `
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

func (pg *postgres) CreateLoginAttempt(db storage.Querier, nu *models.LoginAttempt) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(loginAttemptCreationQuery, &nu.Username, &nu.Successful).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const loginAttemptUpdateQuery = `
    UPDATE login_attempts
    SET
        username = $1,
        successful = $2,
        updated_on = NOW()
    WHERE id = $3
    RETURNING updated_on;
`

func (pg *postgres) UpdateLoginAttempt(db storage.Querier, updated *models.LoginAttempt) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(loginAttemptUpdateQuery, &updated.Username, &updated.Successful, &updated.ID).Scan(&t)
	return t, err
}

const loginAttemptDeletionQuery = `
    UPDATE login_attempts
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteLoginAttempt(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(loginAttemptDeletionQuery, id).Scan(&t)
	return t, err
}
