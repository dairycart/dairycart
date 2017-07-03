package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dummySalt []byte
var userTableHeaders []string
var exampleUserData []driver.Value

const (
	examplePassword       = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
	hashedExamplePassword = "$2a$13$hsflIwHM55jooxaTmYahhOO8LdfI.utMBjpHe5Fr311W4PpRxqyXm"
)

func init() {
	os.Setenv("DAIRYSECRET", "do-not-use-secrets-like-this-plz")
	dummySalt = []byte("farts")
	userTableHeaders = strings.Split(usersTableHeaders, ", ")
	exampleUserData = []driver.Value{
		1, "Frank", "Zappa", "frankzappa", "frank@zappa.com", hashedExamplePassword, dummySalt, true, nil, generateExampleTimeForTests(), nil, nil,
	}
}

func setExpectationsForUserExistence(mock sqlmock.Sqlmock, username string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(userExistenceQuery)).
		WithArgs(username).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserExistenceByID(mock sqlmock.Sqlmock, id string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(userExistenceQueryByID)).
		WithArgs(id).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserCreation(mock sqlmock.Sqlmock, u *User, err error) {
	// can't expect args here because we can't predict the salt/hash
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(u.ID)
	query, _ := buildUserCreationQuery(u)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserRetrieval(mock sqlmock.Sqlmock, email string, err error) {
	exampleRows := sqlmock.NewRows(userTableHeaders).AddRow(exampleUserData...)
	query, rawArgs := buildUserSelectionQuery(email)
	query = formatQueryForSQLMock(query)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserDeletion(mock sqlmock.Sqlmock, id uint64, err error) {
	mock.ExpectExec(formatQueryForSQLMock(userDeletionQuery)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForPasswordResetCreation(mock sqlmock.Sqlmock, id uint64, resetToken string, err error) {
	query, rawArgs := buildPasswordResetRowCreationQuery(id, resetToken)
	args := argsToDriverValues(rawArgs)
	mock.ExpectExec(formatQueryForSQLMock(query)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForPasswordResetCreationWithNoSpecificResetToken(mock sqlmock.Sqlmock, id uint64, err error) {
	query, _ := buildPasswordResetRowCreationQuery(id, "reset-token")
	mock.ExpectExec(formatQueryForSQLMock(query)).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
}

func setExpectationsForPasswordResetEntryExistenceByResetToken(mock sqlmock.Sqlmock, resetToken string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(passwordResetExistenceQuery)
	mock.ExpectQuery(query).
		WithArgs(resetToken).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForPasswordResetEntryExistenceByUserID(mock sqlmock.Sqlmock, userID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := formatQueryForSQLMock(passwordResetExistenceQueryForUserID)
	mock.ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestValidateSessionCookieMiddleware(t *testing.T) {
	t.Parallel()

	handlerWasCalled := false
	exampleHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
	}

	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.Nil(t, err)

	session, err := testUtil.Store.Get(req, dairycartCookieName)
	assert.Nil(t, err)
	session.Values[sessionAuthorizedKeyName] = true
	session.Save(req, testUtil.Response)

	validateSessionCookieMiddleware(testUtil.Response, req, testUtil.Store, exampleHandler)
	assert.True(t, handlerWasCalled)
}

func TestValidateSessionCookieMiddlewareWithInvalidCookie(t *testing.T) {
	t.Parallel()

	handlerWasCalled := false
	exampleHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
	}

	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.Nil(t, err)

	validateSessionCookieMiddleware(testUtil.Response, req, testUtil.Store, exampleHandler)
	assert.False(t, handlerWasCalled)
}

func TestCreateUserFromInput(t *testing.T) {
	t.Parallel()

	exampleUserInput := &UserCreationInput{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
		IsAdmin:   true,
	}
	expected := &User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		IsAdmin:   true,
	}

	actual, err := createUserFromInput(exampleUserInput)
	assert.Nil(t, err)

	assert.Equal(t, expected.FirstName, actual.FirstName, "FirstName fields should match")
	assert.Equal(t, expected.LastName, actual.LastName, "LastName fields should match")
	assert.Equal(t, expected.Email, actual.Email, "Email fields should match")
	assert.Equal(t, expected.IsAdmin, actual.IsAdmin, "IsAdmin fields should match")
	assert.NotEqual(t, expected.Password, actual.Password, "Generated User password should not have the same password as the user input")
	assert.Equal(t, saltSize, len(actual.Salt), fmt.Sprintf("Generated salt should be %d bytes large", saltSize))
}

func TestGenerateSalt(t *testing.T) {
	t.Parallel()
	salt, err := generateSalt()
	assert.Nil(t, err)
	assert.Equal(t, saltSize, len(salt), fmt.Sprintf("Generated salt should be %d bytes large", saltSize))
}

func TestSaltAndHashPassword(t *testing.T) {
	t.Parallel()
	salt := []byte(strings.Repeat("go", 64))
	saltedPass := append(salt, examplePassword...)

	actual, err := saltAndHashPassword(examplePassword, salt)
	assert.Nil(t, err)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(actual), saltedPass))
}

func TestCreateUserInDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)
	newID, err := createUserInDB(testUtil.DB, exampleUser)

	assert.Nil(t, err)
	assert.Equal(t, exampleUser.ID, newID, "createProductInDB should return the created ID")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateUserInDBWhenErrorOccurs(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	setExpectationsForUserCreation(testUtil.Mock, exampleUser, arbitraryError)
	_, err := createUserInDB(testUtil.DB, exampleUser)

	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveUserFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	expected := User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Username:  "frankzappa",
		Email:     "frank@zappa.com",
		Password:  hashedExamplePassword,
		IsAdmin:   true,
		// for some reason I'm too stupid to understand, go wants to copy this value with a cap of 8
		Salt: dummySalt, //[0:len(dummySalt):len(dummySalt)],
	}

	setExpectationsForUserRetrieval(testUtil.Mock, expected.Username, nil)
	actual, err := retrieveUserFromDB(testUtil.DB, expected.Username)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "expected and actual users should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestPasswordIsValid(t *testing.T) {
	t.Parallel()

	input := &UserLoginInput{
		Password: examplePassword,
	}

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordIsValid(input, exampleUser)
	assert.True(t, actual)
}

func TestPasswordIsValidFailsWhenPasswordsDoNotMatch(t *testing.T) {
	t.Parallel()

	input := &UserLoginInput{
		Password: "password",
	}

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordIsValid(input, exampleUser)
	assert.False(t, actual)
}

func TestPasswordIsValidWithVeryLongPassword(t *testing.T) {
	t.Parallel()

	input := &UserLoginInput{
		Password: examplePassword,
	}

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordIsValid(input, exampleUser)
	assert.True(t, actual)
}

func TestArchiveUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForUserDeletion(testUtil.Mock, 1, nil)

	err := archiveUser(testUtil.DB, 1)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreatePasswordResetEntryInDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	userID := uint64(1)
	resetToken := "reset-token"
	setExpectationsForPasswordResetCreation(testUtil.Mock, userID, resetToken, nil)

	err := createPasswordResetEntryInDatabase(testUtil.DB, userID, resetToken)
	assert.Nil(t, err)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestUserCreationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Username, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerFailsWhenCreatingAdminUsersWithoutAlreadyHavingAdminUserStatusJeezLouiseThisIsALongFunctionName(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s",
			"is_admin": true
		}
	`, examplePassword)

	// exampleUser := &User{
	// 	DBRow: DBRow{
	// 		ID:        1,
	// 		CreatedOn: generateExampleTimeForTests(),
	// 	},
	// 	FirstName: "Frank",
	// 	LastName:  "Zappa",
	// 	Email:     "frank@zappa.com",
	// 	Username:  "frankzappa",
	// 	Password:  examplePassword,
	// }

	// setExpectationsForUserExistence(testUtil.Mock, exampleUser.Username, false, nil)
	// setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusForbidden, testUtil.Response.Code, "status code should be 401")
	// ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerForAlreadyExistentUsername(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	setExpectationsForUserExistence(testUtil.Mock, "frankzappa", true, nil)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerWithErrorCreatingUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Username, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusCreated, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerWhenErrorEncounteredInsertingIntoDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  "password",
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Username, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, arbitraryError)

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, nil)

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithInvalidLoginInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler shouldn't attach a cookie when request is invalid")
	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithNoMatchingUserInDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, sql.ErrNoRows)

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithErrorRetrievingUserFromDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, arbitraryError)

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithInvalidPassword(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := `
		{
			"username": "frankzappa",
			"password": "password"
		}
	`

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, nil)

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 401, testUtil.Response.Code, "status code should be 401")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLogoutHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/logout", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "logout handler should attach a cookie when request is valid")
	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
}

func TestUserDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, true, nil)
	setExpectationsForUserDeletion(testUtil.Mock, exampleID, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserDeletionHandlerForNonexistentUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, false, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserDeletionHandlerWithArbitraryErrorWhenDeletingUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, true, nil)
	setExpectationsForUserDeletion(testUtil.Mock, exampleID, arbitraryError)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}
	userID := strconv.Itoa(int(exampleUser.ID))

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
	assert.Nil(t, err)

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, nil)
	setExpectationsForPasswordResetEntryExistenceByUserID(testUtil.Mock, userID, false, nil)
	setExpectationsForPasswordResetCreationWithNoSpecificResetToken(testUtil.Mock, exampleUser.ID, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandlerWithNonexistentUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
	assert.Nil(t, err)

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, sql.ErrNoRows)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandlerWithErrorRetrievingUserFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
	assert.Nil(t, err)

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, arbitraryError)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandlerWithAlreadyExistentPasswordResetEntry(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}
	userID := strconv.Itoa(int(exampleUser.ID))

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
	assert.Nil(t, err)

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, nil)
	setExpectationsForPasswordResetEntryExistenceByUserID(testUtil.Mock, userID, true, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusBadRequest, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserForgottenPasswordHandlerWithErrorCreatingResetToken(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}
	userID := strconv.Itoa(int(exampleUser.ID))

	req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
	assert.Nil(t, err)

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Username, nil)
	setExpectationsForPasswordResetEntryExistenceByUserID(testUtil.Mock, userID, false, nil)
	setExpectationsForPasswordResetCreationWithNoSpecificResetToken(testUtil.Mock, exampleUser.ID, arbitraryError)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusInternalServerError, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestPasswordResetValidationHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleResetToken := "reset-token"
	req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
	assert.Nil(t, err)

	setExpectationsForPasswordResetEntryExistenceByResetToken(testUtil.Mock, exampleResetToken, true, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusOK, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestPasswordResetValidationHandlerForNonexistentToken(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleResetToken := "reset-token"
	req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
	assert.Nil(t, err)

	setExpectationsForPasswordResetEntryExistenceByResetToken(testUtil.Mock, exampleResetToken, false, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, http.StatusNotFound, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
