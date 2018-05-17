package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const userQueryByUsername = `
    SELECT
        id,
        first_name,
        last_name,
        username,
        email,
        password,
        salt,
        is_admin,
        password_last_changed_on,
        created_on,
        updated_on,
        archived_on
    FROM
        users
    WHERE
        archived_on is null
    AND
        username = $1
`

func (pg *postgres) GetUserByUsername(db database.Querier, username string) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow(userQueryByUsername, username).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Password, &u.Salt, &u.IsAdmin, &u.PasswordLastChangedOn, &u.CreatedOn, &u.UpdatedOn, &u.ArchivedOn)
	return u, err
}

const userWithUsernameExistenceQuery = `SELECT EXISTS(SELECT id FROM users WHERE username = $1 and archived_on IS NULL);`

func (pg *postgres) UserWithUsernameExists(db database.Querier, sku string) (bool, error) {
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

func (pg *postgres) UserExists(db database.Querier, id uint64) (bool, error) {
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
        id,
        first_name,
        last_name,
        username,
        email,
        password,
        salt,
        is_admin,
        password_last_changed_on,
        created_on,
        updated_on,
        archived_on
    FROM
        users
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetUser(db database.Querier, id uint64) (*models.User, error) {
	u := &models.User{}

	err := db.QueryRow(userSelectionQuery, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Password, &u.Salt, &u.IsAdmin, &u.PasswordLastChangedOn, &u.CreatedOn, &u.UpdatedOn, &u.ArchivedOn)

	return u, err
}

func buildUserListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"first_name",
			"last_name",
			"username",
			"email",
			"password",
			"salt",
			"is_admin",
			"password_last_changed_on",
			"created_on",
			"updated_on",
			"archived_on",
		).
		From("users")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetUserList(db database.Querier, qf *models.QueryFilter) ([]models.User, error) {
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
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.Salt,
			&u.IsAdmin,
			&u.PasswordLastChangedOn,
			&u.CreatedOn,
			&u.UpdatedOn,
			&u.ArchivedOn,
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

func (pg *postgres) GetUserCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildUserCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const userCreationQuery = `
    INSERT INTO users
        (
            first_name, last_name, username, email, password, salt, is_admin, password_last_changed_on
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateUser(db database.Querier, nu *models.User) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(userCreationQuery, &nu.FirstName, &nu.LastName, &nu.Username, &nu.Email, &nu.Password, &nu.Salt, &nu.IsAdmin, &nu.PasswordLastChangedOn).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const userUpdateQuery = `
    UPDATE users
    SET
        first_name = $1,
        last_name = $2,
        username = $3,
        email = $4,
        password = $5,
        salt = $6,
        is_admin = $7,
        password_last_changed_on = $8,
        updated_on = NOW()
    WHERE id = $9
    RETURNING updated_on;
`

func (pg *postgres) UpdateUser(db database.Querier, updated *models.User) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(userUpdateQuery, &updated.FirstName, &updated.LastName, &updated.Username, &updated.Email, &updated.Password, &updated.Salt, &updated.IsAdmin, &updated.PasswordLastChangedOn, &updated.ID).Scan(&t)
	return t, err
}

const userDeletionQuery = `
    UPDATE users
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteUser(db database.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(userDeletionQuery, id).Scan(&t)
	return t, err
}
