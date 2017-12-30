package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const userQueryByUsername = `
    SELECT
        password,
        is_admin,
        password_last_changed_on,
        salt,
        last_name,
        id,
        username,
        created_on,
        archived_on,
        first_name,
        updated_on,
        email
    FROM
        users
    WHERE
        archived_on is null
    AND
        username = $1
`

func (pg *postgres) GetUserByUsername(db storage.Querier, username string) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow(userQueryByUsername, username).Scan(&u.Password, &u.IsAdmin, &u.PasswordLastChangedOn, &u.Salt, &u.LastName, &u.ID, &u.Username, &u.CreatedOn, &u.ArchivedOn, &u.FirstName, &u.UpdatedOn, &u.Email)
	return u, err
}

const userWithUsernameExistenceQuery = `SELECT EXISTS(SELECT id FROM users WHERE username = $1 and archived_on IS NULL);`

func (pg *postgres) UserWithUsernameExists(db storage.Querier, sku string) (bool, error) {
	var exists string

	err := db.QueryRow(userWithUsernameExistenceQuery, sku).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const userExistenceQuery = `SELECT EXISTS(SELECT id FROM users WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) UserExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(userExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const userSelectionQuery = `
    SELECT
        password,
        is_admin,
        password_last_changed_on,
        salt,
        last_name,
        id,
        username,
        created_on,
        archived_on,
        first_name,
        updated_on,
        email
    FROM
        users
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetUser(db storage.Querier, id uint64) (*models.User, error) {
	u := &models.User{}

	err := db.QueryRow(userSelectionQuery, id).Scan(&u.Password, &u.IsAdmin, &u.PasswordLastChangedOn, &u.Salt, &u.LastName, &u.ID, &u.Username, &u.CreatedOn, &u.ArchivedOn, &u.FirstName, &u.UpdatedOn, &u.Email)

	return u, err
}

func buildUserListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"password",
			"is_admin",
			"password_last_changed_on",
			"salt",
			"last_name",
			"id",
			"username",
			"created_on",
			"archived_on",
			"first_name",
			"updated_on",
			"email",
		).
		From("users")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetUserList(db storage.Querier, qf *models.QueryFilter) ([]models.User, error) {
	var list []models.User
	query, args := buildUserListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.Password,
			&u.IsAdmin,
			&u.PasswordLastChangedOn,
			&u.Salt,
			&u.LastName,
			&u.ID,
			&u.Username,
			&u.CreatedOn,
			&u.ArchivedOn,
			&u.FirstName,
			&u.UpdatedOn,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildUserCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("users")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetUserCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildUserCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const userCreationQuery = `
    INSERT INTO users
        (
            password, is_admin, password_last_changed_on, salt, last_name, username, first_name, email
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateUser(db storage.Querier, nu *models.User) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(userCreationQuery, &nu.Password, &nu.IsAdmin, &nu.PasswordLastChangedOn, &nu.Salt, &nu.LastName, &nu.Username, &nu.FirstName, &nu.Email).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const userUpdateQuery = `
    UPDATE users
    SET
        password = $1,
        is_admin = $2,
        password_last_changed_on = $3,
        salt = $4,
        last_name = $5,
        username = $6,
        first_name = $7,
        email = $8,
        updated_on = NOW()
    WHERE id = $9
    RETURNING updated_on;
`

func (pg *postgres) UpdateUser(db storage.Querier, updated *models.User) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(userUpdateQuery, &updated.Password, &updated.IsAdmin, &updated.PasswordLastChangedOn, &updated.Salt, &updated.LastName, &updated.Username, &updated.FirstName, &updated.Email, &updated.ID).Scan(&t)
	return t, err
}

const userDeletionQuery = `
    UPDATE users
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteUser(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(userDeletionQuery, id).Scan(&t)
	return t, err
}
