package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

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

func (pg *postgres) GetUser(db storage.Querier, id uint64) (*models.User, error) {
	u := &models.User{}

	err := db.QueryRow(userSelectionQuery, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Password, &u.Salt, &u.IsAdmin, &u.PasswordLastChangedOn, &u.CreatedOn, &u.UpdatedOn, &u.ArchivedOn)

	return u, err
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

func (pg *postgres) CreateUser(db storage.Querier, nu *models.User) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(userCreationQuery, &nu.FirstName, &nu.LastName, &nu.Username, &nu.Email, &nu.Password, &nu.Salt, &nu.IsAdmin, &nu.PasswordLastChangedOn).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
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

func (pg *postgres) UpdateUser(db storage.Querier, updated *models.User) (time.Time, error) {
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

func (pg *postgres) DeleteUser(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(userDeletionQuery, id).Scan(&t)
	return t, err
}
