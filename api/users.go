package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairymodels/v1"

	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost                 = bcrypt.DefaultCost + 3
	saltSize                 = 32
	resetTokenSize           = 128
	minimumPasswordSize      = 32
	dairycartCookieName      = "dairycart"
	sessionAdminKeyName      = "is_admin"
	sessionUserIDKeyName     = "user_id"
	sessionAuthorizedKeyName = "authenticated"
)

// DisplayUser represents a Dairycart user we can return in responses
// TODO: the main reason for doing this is so we don't end up returning
// the password hash to the user, but there's bound to be a way to reuse
// that struct
type DisplayUser struct {
	ID         uint64            `json:"id"`
	FirstName  string            `json:"first_name"`
	LastName   string            `json:"last_name"`
	Email      string            `json:"email"`
	IsAdmin    bool              `json:"is_admin"`
	CreatedOn  time.Time         `json:"created_on"`
	UpdatedOn  *models.Dairytime `json:"updated_on,omitempty"`
	ArchivedOn *models.Dairytime `json:"archived_on,omitempty"`
}

// UserLoginInput represents the payload used to log in a Dairycart user
type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func passwordIsValid(s string) bool {
	return len(s) >= minimumPasswordSize
}

func createUserFromInput(in *models.UserCreationInput) (*models.User, error) {
	salt, err := generateSalt()
	// COVERAGE NOTE: I cannot seem to synthesize this error for the sake of testing, so if you're
	// seeing this in a coverage report and the line below is red, just know that I tried. :(
	if err != nil {
		return nil, err
	}

	saltedAndHashedPassword, err := saltAndHashPassword(in.Password, salt)
	// COVERAGE NOTE: see above
	if err != nil {
		return nil, err
	}

	user := &models.User{
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

func createUserFromUpdateInput(in *models.UserUpdateInput, hashedPassword string) *models.User {
	out := &models.User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  hashedPassword,
	}
	return out
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

func passwordMatches(password string, u *models.User) bool {
	saltedInputPassword := append(u.Salt, password...)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), saltedInputPassword)
	return err == nil
}

func validateUserCreationInput(in *models.UserCreationInput) error {
	if in == nil {
		return errors.New("invalid user creation input")
	}
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		log.Printf("error parsing email address '%s': %v\n", in.Email, err)
		return errors.New("email address must be valid")
	}
	if in.FirstName == "" {
		return errors.New("first name must not be empty")
	}
	if in.LastName == "" {
		return errors.New("last name must not be empty")
	}
	if in.Username == "" {
		return errors.New("username must not be empty")
	}
	if len(in.Password) < minimumPasswordSize {
		return errors.New(fmt.Sprintf("password must be at least %d characters", minimumPasswordSize))
	}
	return nil
}

func buildUserCreationHandler(db *sql.DB, client database.Storer, store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userInput := &models.UserCreationInput{}
		err := validateRequestInput(req, userInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		err = validateUserCreationInput(userInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInvalidRequestCookie(res)
			return
		}

		if userInput.IsAdmin {
			// only an admin user can create an admin user
			if admin, ok := session.Values[sessionAdminKeyName].(bool); !ok || !admin {
				res.WriteHeader(http.StatusForbidden)
				errRes := &ErrorResponse{
					Status:  http.StatusForbidden,
					Message: "User is not authorized to create admin users",
				}
				json.NewEncoder(res).Encode(errRes)
				return
			}
		}

		// can't create a user with an email that already exists!
		exists, err := client.UserWithUsernameExists(db, userInput.Username)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, errors.New("username already taken"))
			return
		}

		newUser, err := createUserFromInput(userInput)
		// COVERAGE NOTE: see note in createUserFromInput
		if err != nil {
			notifyOfInternalIssue(res, err, "creating user")
			return
		}

		createdUserID, createdOn, err := client.CreateUser(db, newUser)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert user in database")
			return
		}

		responseUser := &DisplayUser{
			ID:        createdUserID,
			CreatedOn: createdOn,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
			IsAdmin:   newUser.IsAdmin,
		}
		session.Values[sessionUserIDKeyName] = createdUserID
		session.Values[sessionAuthorizedKeyName] = true
		session.Values[sessionAdminKeyName] = newUser.IsAdmin
		session.Save(req, res)

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(responseUser)
	}
}

func buildUserLoginHandler(db *sql.DB, client database.Storer, store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginInput := &UserLoginInput{}
		err := validateRequestInput(req, loginInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		username := loginInput.Username

		exhaustedAttempts, err := client.LoginAttemptsHaveBeenExhausted(db, username)
		if exhaustedAttempts {
			notifyOfExaustedAuthenticationAttempts(res)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		// TODO: we should ensure there isn't an unsatisfied password reset token requested before allowing login

		user, err := client.GetUserByUsername(db, username)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user", username)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		loginValid := passwordMatches(loginInput.Password, user)
		_, _, err = client.CreateLoginAttempt(db, &models.LoginAttempt{Username: username, Successful: loginValid})
		if err != nil {
			notifyOfInternalIssue(res, err, "create login attempt entry")
			return
		}

		if !loginValid {
			notifyOfInvalidAuthenticationAttempt(res)
			return
		}

		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInvalidRequestCookie(res)
			return
		}

		statusToWrite := http.StatusUnauthorized
		if loginValid {
			statusToWrite = http.StatusOK
			session.Values[sessionUserIDKeyName] = user.ID
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
			notifyOfInvalidRequestCookie(res)
			return
		}
		session.Values[sessionAuthorizedKeyName] = false
		session.Save(req, res)
		res.WriteHeader(http.StatusOK)
	}
}

func buildUserDeletionHandler(db *sql.DB, client database.Storer, store *sessions.CookieStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userID := chi.URLParam(req, "user_id")
		// we can eat this error because Mux takes care of validating route params for us
		userIDInt, _ := strconv.ParseInt(userID, 10, 64)
		userIDInt64 := uint64(userIDInt)

		// can't delete a user that doesn't already exist!
		user, err := client.GetUser(db, userIDInt64)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user", userID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		session, err := store.Get(req, dairycartCookieName)
		if err != nil {
			notifyOfInvalidRequestCookie(res)
			return
		}

		// only an admin user can delete an admin user
		admin, ok := session.Values[sessionAdminKeyName].(bool)
		if !ok || !admin {
			res.WriteHeader(http.StatusForbidden)
			errRes := &ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "User is not authorized to delete users",
			}
			json.NewEncoder(res).Encode(errRes)
			return
		} else if admin {
			archivedOn, err := client.DeleteUser(db, userIDInt64)
			user.ArchivedOn = &models.Dairytime{Time: archivedOn}
			if err != nil {
				notifyOfInternalIssue(res, err, "archive user")
				return
			}
		}

		json.NewEncoder(res).Encode(user)
	}
}

func buildUserForgottenPasswordHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginInput := &UserLoginInput{}
		err := validateRequestInput(req, loginInput)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		username := loginInput.Username

		user, err := client.GetUserByUsername(db, username)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user", username)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		exists, err := client.PasswordResetTokenForUserIDExists(db, user.ID)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, errors.New("user has existent, non-expired password reset request"))
			return
		}

		resetToken := &models.PasswordResetToken{
			UserID: user.ID,
			Token:  uniuri.NewLen(resetTokenSize),
		}

		resetToken.ID, resetToken.CreatedOn, err = client.CreatePasswordResetToken(db, resetToken)
		if err != nil {
			notifyOfInternalIssue(res, err, "read session data")
			return
		}

		json.NewEncoder(res).Encode(resetToken)
	}
}

// TODO: rethinking having this as a mere validation handler, instead of a password resetting handler
func buildUserPasswordResetTokenValidationHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		resetToken := chi.URLParam(req, "reset_token")

		exists, err := client.PasswordResetTokenWithTokenExists(db, resetToken)
		if err != nil || !exists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func buildUserUpdateHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userID := chi.URLParam(req, "user_id")
		// eating these errors because Chi should validate these for us.
		userIDInt, _ := strconv.Atoi(userID)
		userIDInt64 := uint64(userIDInt)

		updatedUserInfo := &models.UserUpdateInput{}
		err := validateRequestInput(req, updatedUserInfo)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		newPassword := updatedUserInfo.NewPassword
		passwordChanged := newPassword != ""
		if passwordChanged && !passwordIsValid(newPassword) {
			notifyOfInvalidRequestBody(res, errors.New("provided password is invalid"))
			return
		}

		existingUser, err := client.GetUser(db, userIDInt64)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "user ID", userID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve user")
			return
		}

		loginValid := passwordMatches(updatedUserInfo.CurrentPassword, existingUser)
		if !loginValid {
			notifyOfInvalidAuthenticationAttempt(res)
			return
		}

		// TODO: Evaluate whether or not I should be reusing the salt here.
		hashedPassword := existingUser.Password
		if passwordChanged {
			var err error
			hashedPassword, err = saltAndHashPassword(newPassword, existingUser.Salt)
			// COVERAGE NOTE: see note in createUserFromInput
			if err != nil {
				notifyOfInternalIssue(res, err, "update user")
				return
			}
		}

		updatedUser := createUserFromUpdateInput(updatedUserInfo, hashedPassword)

		mergo.Merge(updatedUser, existingUser)

		// FIXME: this isn't how this should be done
		if passwordChanged {
			updatedUser.PasswordLastChangedOn = &models.Dairytime{Time: time.Now()}
		}

		updatedOn, err := client.UpdateUser(db, updatedUser)
		if err != nil {
			notifyOfInternalIssue(res, err, "update user")
			return
		}
		updatedUser.UpdatedOn = &models.Dairytime{Time: updatedOn}

		json.NewEncoder(res).Encode(updatedUser)
	}
}
