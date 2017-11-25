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

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
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

func setExpectationsForUserCreation(mock sqlmock.Sqlmock, u *models.User, err error) {
	// can't expect args here because we can't predict the salt/hash
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(u.ID, generateExampleTimeForTests())
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

func setExpectationsForUserRetrievalByID(mock sqlmock.Sqlmock, userID uint64, err error) {
	exampleRows := sqlmock.NewRows(userTableHeaders).AddRow(exampleUserData...)
	query, rawArgs := buildUserSelectionQueryByID(userID)
	query = formatQueryForSQLMock(query)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserUpdate(mock sqlmock.Sqlmock, u *models.User, passwordChanged bool, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	rawQuery, rawArgs := buildUserUpdateQuery(u, passwordChanged)
	query := formatQueryForSQLMock(rawQuery)
	args := argsToDriverValues(rawArgs)
	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserUpdateWithoutSpecifyingPassword(mock sqlmock.Sqlmock, u *models.User, passwordChanged bool, err error) {
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	rawQuery, _ := buildUserUpdateQuery(u, passwordChanged)
	query := formatQueryForSQLMock(rawQuery)
	mock.ExpectQuery(query).
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

func setExpectationsForLoginAttemptExhaustionQuery(mock sqlmock.Sqlmock, username string, exhausted bool, err error) {
	count := 0
	if exhausted {
		count = 666
	}
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(count)
	query := formatQueryForSQLMock(loginAttemptExhaustionQuery)
	mock.ExpectQuery(query).
		WithArgs(username).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForLoginAttemptCreationQuery(mock sqlmock.Sqlmock, username string, successful bool, err error) {
	query, rawArgs := buildLoginAttemptCreationQuery(username, successful)
	args := argsToDriverValues(rawArgs)
	mock.ExpectExec(formatQueryForSQLMock(query)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(1, 1)).
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

func TestPasswordIsValid(t *testing.T) {
	inputOutputMap := map[string]bool{
		// the worst password ever
		"password": false,
		// should pass, but only barely
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaA1!": true,
		// the example password we've already been using all over the place
		examplePassword: true,
	}

	for in, expected := range inputOutputMap {
		actual := passwordIsValid(in)
		msg := fmt.Sprintf("expected password `%s` to be considered valid, but it was considered invalid", in)
		if !expected {
			msg = fmt.Sprintf("expected password `%s` to be considered invalid, but it was considered valid", in)
		}
		assert.Equal(t, expected, actual, msg)
	}
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
	expected := &models.User{
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

func TestCreateUserFromUpdateInput(t *testing.T) {
	t.Parallel()

	exampleUserUpdateInput := &UserUpdateInput{
		FirstName: "FirstName",
		LastName:  "LastName",
		Username:  "Username",
		Email:     "Email",
	}
	expected := &models.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Username:  "Username",
		Email:     "Email",
		Password:  hashedExamplePassword,
	}
	actual := createUserFromUpdateInput(exampleUserUpdateInput, hashedExamplePassword)

	assert.Equal(t, expected, actual, "expected and actual output were not equal")
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

	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)
	newID, createdOn, err := createUserInDB(testUtil.DB, exampleUser)

	assert.Nil(t, err)
	assert.Equal(t, exampleUser.ID, newID, "createUserInDB should return the created ID")
	assert.Equal(t, exampleUser.CreatedOn, createdOn, "createUserInDB should return the created ID")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestCreateUserInDBWhenErrorOccurs(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	setExpectationsForUserCreation(testUtil.Mock, exampleUser, generateArbitraryError())
	_, _, err := createUserInDB(testUtil.DB, exampleUser)

	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRetrieveUserFromDB(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	expected := models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "Frank",
		LastName:  "Zappa",
		Username:  "frankzappa",
		Email:     "frank@zappa.com",
		Password:  hashedExamplePassword,
		IsAdmin:   true,
		Salt:      dummySalt,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, expected.Username, nil)
	actual, err := retrieveUserFromDB(testUtil.DB, expected.Username)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "expected and actual users should match")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

// func TestRetrieveUserFromDBByID(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	expected := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, expected.ID, nil)
// 	actual, err := retrieveUserFromDBByID(testUtil.DB, expected.ID)
// 	assert.Nil(t, err)

// 	assert.Equal(t, expected, actual, "expected and actual users should match")
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

func TestPasswordMatches(t *testing.T) {
	t.Parallel()

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := &models.User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordMatches(examplePassword, exampleUser)
	assert.True(t, actual)
}

func TestPasswordMatchesFailsWhenPasswordsDoNotMatch(t *testing.T) {
	t.Parallel()

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := &models.User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordMatches("password", exampleUser)
	assert.False(t, actual)
}

func TestPasswordMatchesWithVeryLongPassword(t *testing.T) {
	t.Parallel()

	saltedPasswordHash, err := saltAndHashPassword(examplePassword, dummySalt)
	assert.Nil(t, err)
	exampleUser := &models.User{
		Password: saltedPasswordHash,
		Salt:     dummySalt,
	}

	actual := passwordMatches(examplePassword, exampleUser)
	assert.True(t, actual)
}

func TestUpdateUserInDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	examplePasswordChanged := false
	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "Frank",
		LastName:  "Zappa",
		Username:  "frankzappa",
		Email:     "frank@zappa.com",
		Password:  hashedExamplePassword,
		IsAdmin:   true,
		Salt:      dummySalt,
	}
	setExpectationsForUserUpdate(testUtil.Mock, exampleUser, examplePasswordChanged, nil)

	updatedOn, err := updateUserInDatabase(testUtil.DB, exampleUser, examplePasswordChanged)
	assert.Nil(t, err)
	assert.Equal(t, generateExampleTimeForTests(), updatedOn, "updateUserInDatabase should return valid updated_on time")

	ensureExpectationsWereMet(t, testUtil.Mock)
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

func TestLoginAttemptsHaveBeenExhausted(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUsername := "farts"
	setExpectationsForLoginAttemptExhaustionQuery(testUtil.Mock, exampleUsername, false, nil)

	exhausted, err := loginAttemptsHaveBeenExhausted(testUtil.DB, exampleUsername)
	assert.False(t, exhausted)
	assert.Nil(t, err)
}

func TestLoginAttemptsHaveBeenExhaustedWithErrorExecutingQuery(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUsername := "farts"
	setExpectationsForLoginAttemptExhaustionQuery(testUtil.Mock, exampleUsername, false, generateArbitraryError())

	exhausted, err := loginAttemptsHaveBeenExhausted(testUtil.DB, exampleUsername)
	assert.False(t, exhausted)
	assert.NotNil(t, err)
}

func TestCreateLoginAttemptRowInDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleUsername := "farts"
	successful := true
	setExpectationsForLoginAttemptCreationQuery(testUtil.Mock, exampleUsername, successful, nil)

	err := createLoginAttemptRowInDatabase(testUtil.DB, exampleUsername, successful)
	assert.Nil(t, err)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestUserCreationHandler(t *testing.T) {
	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleAdminInput := fmt.Sprintf(`
			{
				"first_name": "Frank",
				"last_name": "Zappa",
				"email": "frank@zappa.com",
				"username": "frankzappa",
				"password": "%s",
				"is_admin": true
			}
		`, examplePassword)

	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("CreateUser", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("already existent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("creating an admin user as a non-admin user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleAdminInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusForbidden)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		attachBadCookieToRequest(req)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("CreateUser", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestUserLoginHandler(t *testing.T) {
	// t.Parallel()
	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  hashedExamplePassword,
		Salt:      dummySalt,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("CreateLoginAttempt", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid login input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("CreateLoginAttempt", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("when login attempts have been exhausted", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusUnauthorized)
	})

	t.Run("with error checking login attempts", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("without matching user in database", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating a login attempt", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("CreateLoginAttempt", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusInternalServerError)

	})

	t.Run("with invalid password", func(*testing.T) {
		invalidInput := `
			{
				"username": "frankzappa",
				"password": "password"
			}
		`

		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("CreateLoginAttempt", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(invalidInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusUnauthorized)

	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("LoginAttemptsHaveBeenExhausted", mock.Anything, exampleUser.Username).
			Return(false, nil)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("CreateLoginAttempt", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		attachBadCookieToRequest(req)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})
}

func TestUserLogoutHandler(t *testing.T) {
	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariables(t)

		req, err := http.NewRequest(http.MethodPost, "/logout", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "logout handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariables(t)

		req, err := http.NewRequest(http.MethodPost, "/logout", nil)
		assert.Nil(t, err)
		attachBadCookieToRequest(req)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "logout handler should not attach a cookie when request is invalid")
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})
}

func TestUserDeletionHandler(t *testing.T) {
	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(true, nil)
		testUtil.MockDB.On("DeleteUser", mock.Anything, exampleID).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)

		cookie, err := buildCookieForRequest(t, testUtil.Store, true, true)
		assert.Nil(t, err)
		req.AddCookie(cookie)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(true, nil)
		testUtil.MockDB.On("DeleteUser", mock.Anything, exampleID).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)
		attachBadCookieToRequest(req)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("when deleting admin user as regular user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(true, nil)
		testUtil.MockDB.On("DeleteUser", mock.Anything, exampleID).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)
		cookie, err := buildCookieForRequest(t, testUtil.Store, true, false)
		assert.Nil(t, err)
		req.AddCookie(cookie)

		setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, true, nil)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusForbidden)
	})

	t.Run("with error deleting user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(true, nil)
		testUtil.MockDB.On("DeleteUser", mock.Anything, exampleID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)

		cookie, err := buildCookieForRequest(t, testUtil.Store, true, true)
		assert.Nil(t, err)
		req.AddCookie(cookie)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestUserForgottenPasswordHandler(t *testing.T) {
	exampleInput := fmt.Sprintf(`
		{
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &models.User{
		ID:        1,
		CreatedOn: generateExampleTimeForTests(),
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Username:  "frankzappa",
		Password:  examplePassword,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("PasswordResetTokenExistsForUserID", mock.Anything, mock.Anything).
			Return(false, nil)
		testUtil.MockDB.On("CreatePasswordResetToken", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving user from db", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with already existent password reset entry", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("PasswordResetTokenExistsForUserID", mock.Anything, mock.Anything).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating reset token", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("PasswordResetTokenExistsForUserID", mock.Anything, mock.Anything).
			Return(false, nil)
		testUtil.MockDB.On("CreatePasswordResetToken", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.DB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

// func TestPasswordResetValidationHandler(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleResetToken := "reset-token"
// 	req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
// 	assert.Nil(t, err)

// 	setExpectationsForPasswordResetEntryExistenceByResetToken(testUtil.Mock, exampleResetToken, true, nil)
// 	testUtil.Router.ServeHTTP(testUtil.Response, req)

// 	assertStatusCode(t, testUtil, http.StatusOK)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestPasswordResetValidationHandlerForNonexistentToken(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleResetToken := "reset-token"
// 	req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
// 	assert.Nil(t, err)

// 	setExpectationsForPasswordResetEntryExistenceByResetToken(testUtil.Mock, exampleResetToken, false, nil)
// 	testUtil.Router.ServeHTTP(testUtil.Response, req)

// 	assertStatusCode(t, testUtil, http.StatusNotFound)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandler(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"username": "captain_beefheart",
//  			"current_password": "%s"
//  		}
//  	`, examplePassword)

// 	examplePasswordChanged := false
// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	afterUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "captain_beefheart",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, nil)
// 	setExpectationsForUserUpdate(testUtil.Mock, afterUser, examplePasswordChanged, nil)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusOK)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWithInvalidInput(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleGarbageInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusBadRequest)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWhileAttemptingToChangePasswordToAnInvalidPassword(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"new_password": "passwordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
//  			"current_password": "%s"
//  		}
//  	`, examplePassword)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusBadRequest)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerForNonexistentUser(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"username": "captain_beefheart",
//  			"current_password": "%s"
//  		}
//  	`, examplePassword)

// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, sql.ErrNoRows)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusNotFound)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWithErrorRetrievingUser(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"username": "captain_beefheart",
//  			"current_password": "%s"
//  		}
//  	`, examplePassword)

// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, generateArbitraryError())

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusInternalServerError)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWhenPasswordDoesNotMatchExpectation(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"username": "captain_beefheart",
//  			"current_password": "%s"
//  		}
//  	`, fmt.Sprintf("%s!", examplePassword))

// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, nil)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusUnauthorized)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWithNewPassword(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleNewPassword := "P@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rd"
// 	exampleNewPasswordHashed := "$2a$13$xhhweT6OnsU7l6GyPGdin.YDANUGnFEu7xJQb7eU/zv4KBCiRwWbC"
// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"new_password": "%s",
//  			"current_password": "%s"
//  		}
//  	`, exampleNewPassword, examplePassword)

// 	examplePasswordChanged := true
// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	afterUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "captain_beefheart",
// 		Email:     "frank@zappa.com",
// 		Password:  exampleNewPasswordHashed,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, nil)
// 	setExpectationsForUserUpdateWithoutSpecifyingPassword(testUtil.Mock, afterUser, examplePasswordChanged, nil)

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusOK)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }

// func TestUserUpdateHandlerWithErrorUpdatingUser(t *testing.T) {
// 	t.Parallel()
// 	testUtil := setupTestVariables(t)

// 	exampleUserUpdateInput := fmt.Sprintf(`
//  		{
//  			"username": "captain_beefheart",
//  			"current_password": "%s"
//  		}
//  	`, examplePassword)

// 	examplePasswordChanged := false
// 	beforeUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "frankzappa",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	afterUser := &models.User{
// 		ID:        1,
// 		CreatedOn: generateExampleTimeForTests(),
// 		FirstName: "Frank",
// 		LastName:  "Zappa",
// 		Username:  "captain_beefheart",
// 		Email:     "frank@zappa.com",
// 		Password:  hashedExamplePassword,
// 		IsAdmin:   true,
// 		Salt:      dummySalt,
// 	}

// 	setExpectationsForUserRetrievalByID(testUtil.Mock, beforeUser.ID, nil)
// 	setExpectationsForUserUpdate(testUtil.Mock, afterUser, examplePasswordChanged, generateArbitraryError())

// 	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(exampleUserUpdateInput))
// 	assert.Nil(t, err)

// 	testUtil.Router.ServeHTTP(testUtil.Response, req)
// 	assertStatusCode(t, testUtil, http.StatusInternalServerError)
// 	ensureExpectationsWereMet(t, testUtil.Mock)
// }
