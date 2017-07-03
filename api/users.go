package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	saltSize                 = 1 << 5
	hashCost                 = bcrypt.DefaultCost + 3
	resetTokenSize           = 1 << 7
	dairycartCookieName      = "dairycart"
	sessionAdminKeyName      = "is_admin"
	sessionAuthorizedKeyName = "authenticated"

	usersTableHeaders       = `id, first_name, last_name, username, email, password, salt, is_admin, password_last_changed_on, created_on, updated_on, archived_on`
	userExistenceQuery      = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND archived_on IS NULL)`
	adminUserExistenceQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE is_admin is true AND archived_on IS NULL)`
	userExistenceQueryByID  = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND archived_on IS NULL)`
	userDeletionQuery       = `UPDATE users SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`

	passwordResetExistenceQueryForUserID = `SELECT EXISTS(SELECT 1 FROM password_reset_tokens WHERE user_id = $1 AND NOW() < expires_on)`
	passwordResetExistenceQuery          = `SELECT EXISTS(SELECT 1 FROM password_reset_tokens WHERE token = $1 AND NOW() < expires_on)`
)

// User represents a Dairycart user
type User struct {
	DBRow
	FirstName             string   `json:"first_name"`
	LastName              string   `json:"last_name"`
	Username              string   `json:"username"`
	Email                 string   `json:"email"`
	Password              string   `json:"password"`
	Salt                  []byte   `json:"salt"`
	IsAdmin               bool     `json:"is_admin"`
	PasswordLastChangedOn NullTime `json:"password_last_changed_on,omitempty"`
}

// DisplayUser represents a Dairycart user we can return in responses
type DisplayUser struct {
	DBRow
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserCreationInput represents the payload used to create a Dairycart user
type UserCreationInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Username  string `json:"username"   validate:"required"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,gte=64"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserLoginInput represents the payload used to log in a Dairycart user
type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateInput represents the payload used to update a Dairycart user
type UserUpdateInput struct {
	FirstName       string `json:"first_name"       validate:"required"`
	LastName        string `json:"last_name"        validate:"required"`
	Username        string `json:"username"         validate:"required"`
	Email           string `json:"email"            validate:"required,email"`
	CurrentPassword string `json:"current_password" validate:"required,gte=64"`
	NewPassword     string `json:"new_password"     validate:"required,gte=64"`
	IsAdmin         bool   `json:"is_admin"`
}

func validateSessionCookieMiddleware(res http.ResponseWriter, req *http.Request, store *sessions.CookieStore, next http.HandlerFunc) {
	session, err := store.Get(req, dairycartCookieName)
	if auth, ok := session.Values[sessionAuthorizedKeyName].(bool); !ok || !auth || err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		errRes := &ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		}
		json.NewEncoder(res).Encode(errRes)
		return
	}
	next(res, req)
}

func createUserFromInput(in *UserCreationInput) (*User, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	saltedAndHashedPassword, err := saltAndHashPassword(in.Password, salt)
	if err != nil {
		return nil, err
	}

	user := &User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Username:  in.Username,
		Password:  string(saltedAndHashedPassword),
		Salt:      salt,
		IsAdmin:   in.IsAdmin,
	}
	return user, nil
}

func generateSalt() ([]byte, error) {
	b := make([]byte, saltSize)
	_, err := rand.Read(b)
	return b, err
}

func saltAndHashPassword(password string, salt []byte) (string, error) {
	passwordToHash := append(salt, password...)
	saltedAndHashedPassword, err := bcrypt.GenerateFromPassword(passwordToHash, hashCost)
	return string(saltedAndHashedPassword), err
}

func createUserInDB(db *sqlx.DB, u *User) (uint64, error) {
	var newUserID uint64
	query, args := buildUserCreationQuery(u)
	err := db.QueryRow(query, args...).Scan(&newUserID)
	return newUserID, err
}

func retrieveUserFromDB(db *sqlx.DB, username string) (User, error) {
	var u User
	query, args := buildUserSelectionQuery(username)
	err := db.Get(&u, query, args...)
	return u, err
}

func passwordIsValid(in *UserLoginInput, u User) bool {
	saltedInputPassword := append(u.Salt, in.Password...)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), saltedInputPassword)
	return err == nil
}

func archiveUser(db *sqlx.DB, id uint64) error {
	_, err := db.Exec(userDeletionQuery, id)
	return err
}

func createPasswordResetEntryInDatabase(db *sqlx.DB, userID uint64, resetToken string) error {
	/*
		NOTE: this docstring is mostly for my own future reference

		I will work to implement the creation of these rows and the validation of their contents,
		but I won't be implementing the actual emailing of users with these reset tokens just yet.
		Mostly because email is a ~*~spooky business~*~ and I have absolutely no idea how to test
		that stuff without getting even real complicated. Towards the end of development, when I
		feel like Dairycart is closer to being ready to release, I will implement this feature and
		test it manually on occasion. RIP to my sweet test coverage number.
	*/
	query, args := buildPasswordResetRowCreationQuery(userID, resetToken)
	_, err := db.Exec(query, args...)
	return err
}

func buildUserCreationHandler(db *sqlx.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userInput := &UserCreationInput{}
		err := validateRequestInput(req, userInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInternalIssue(res, err, "read session data")
			return
		}
		if userInput.IsAdmin {
			// only an admin user can create an admin user
			if admin, ok := session.Values[sessionAdminKeyName].(bool); !ok || !admin {
				http.Error(res, "Forbidden", http.StatusForbidden)
				return
			}
		}

		// can't create a user with an email that already exists!
		exists, err := rowExistsInDB(db, userExistenceQuery, userInput.Username)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, errors.New("username already taken"))
			return
		}

		newUser, err := createUserFromInput(userInput)
		if err != nil {
			notifyOfInternalIssue(res, err, "creating user")
			return
		}

		createdUserID, err := createUserInDB(db, newUser)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert user in database")
			return
		}

		responseUser := &DisplayUser{
			DBRow: DBRow{
				ID:        createdUserID,
				CreatedOn: time.Now(),
			},
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
			IsAdmin:   newUser.IsAdmin,
		}
		session.Values[sessionAuthorizedKeyName] = true
		session.Values[sessionAdminKeyName] = newUser.IsAdmin
		session.Save(req, res)

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(responseUser)
	}
}

func buildUserLoginHandler(db *sqlx.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginInput := &UserLoginInput{}
		err := validateRequestInput(req, loginInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		username := loginInput.Username

		user, err := retrieveUserFromDB(db, username)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user", username)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		loginValid := passwordIsValid(loginInput, user)
		if !loginValid {
			notifyOfInvalidAuthenticationAttempt(res, username)
			return
		}

		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInternalIssue(res, err, "read session data")
			return
		}

		statusToWrite := http.StatusUnauthorized
		if loginValid {
			statusToWrite = http.StatusOK
			session.Values[sessionAuthorizedKeyName] = true
			session.Values[sessionAdminKeyName] = user.IsAdmin
			session.Save(req, res)
		}
		res.WriteHeader(statusToWrite)
	}
}

func buildUserLogoutHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInternalIssue(res, err, "read session data")
			return
		}
		session.Values[sessionAuthorizedKeyName] = false
		session.Save(req, res)
		res.WriteHeader(http.StatusOK)
	}
}

func buildUserDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userID := chi.URLParam(req, "user_id")
		// we can eat this error because Mux takes care of validating route params for us
		userIDInt, _ := strconv.ParseInt(userID, 10, 64)

		// can't delete a user with an email that already exists!
		exists, err := rowExistsInDB(db, userExistenceQueryByID, userID)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "user", userID)
			return
		}

		err = archiveUser(db, uint64(userIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "archive user")
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func buildUserForgottenPasswordHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginInput := &UserLoginInput{}
		err := validateRequestInput(req, loginInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		username := loginInput.Username

		user, err := retrieveUserFromDB(db, username)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user", username)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		userIDString := strconv.Itoa(int(user.ID))
		exists, err := rowExistsInDB(db, passwordResetExistenceQueryForUserID, userIDString)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, errors.New("user has existent, non-expired password reset request"))
			return
		}

		resetToken := uniuri.NewLen(resetTokenSize)
		err = createPasswordResetEntryInDatabase(db, user.ID, resetToken)
		if err != nil {
			notifyOfInternalIssue(res, err, "read session data")
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func buildUserPasswordResetTokenValidationHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		resetToken := chi.URLParam(req, "reset_token")

		exists, err := rowExistsInDB(db, passwordResetExistenceQuery, resetToken)
		if err != nil || !exists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}