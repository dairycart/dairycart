package main

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	dummySalt = []byte("farts")
	userTableHeaders = strings.Split(usersTableHeaders, ", ")
	exampleUserData = []driver.Value{
		1, "Frank", "Zappa", "frankzappa", "frank@zappa.com", hashedExamplePassword, dummySalt, true, nil, generateExampleTimeForTests(), nil, nil,
	}
}

func setExpectationsForUserExistence(mock sqlmock.Sqlmock, email string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(userExistenceQuery)).
		WithArgs(email).
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

func TestValidateSessionCookieMiddleware(t *testing.T) {
	t.Parallel()

	handlerWasCalled := false
	exampleHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
	}

	testUtil := setupTestVariables(t)

	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(t, err)

	session, err := testUtil.Store.Get(req, dairycartCookieName)
	assert.Nil(t, err)
	session.Values["authenticated"] = true
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

	req, err := http.NewRequest("GET", "", nil)
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

func TestValidateUserCreationInput(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword))

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual, err := validateUserCreationInput(req)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
}

func TestValidateUserCreationInputWithAwfulpassword(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"password": "password"
		}
	`)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateUserCreationInput(req)

	assert.NotNil(t, err)
}

func TestValidateUserCreationInputWithGarbageInput(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(exampleGarbageInput)
	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateUserCreationInput(req)

	assert.NotNil(t, err)
}

func TestValidateUserCreationInputWithCompletelyGarbageInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateUserCreationInput(req)

	assert.NotNil(t, err)
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

	setExpectationsForUserRetrieval(testUtil.Mock, expected.Email, nil)
	actual, err := retrieveUserFromDB(testUtil.DB, expected.Email)
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

func TestValidateLoginInput(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword))

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual, err := validateLoginInput(req)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
}

func TestValidateLoginInputWithCompletelyGarbageInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := validateLoginInput(req)

	assert.NotNil(t, err)
}

func TestArchiveUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForUserDeletion(testUtil.Mock, 1, nil)

	err := archiveUser(testUtil.DB, 1)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, testUtil.Mock)
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
		Password:  examplePassword,
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 201, testUtil.Response.Code, "status code should be 201")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 400, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserCreationHandlerForAlreadyExistentUserEmail(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword)

	setExpectationsForUserExistence(testUtil.Mock, "frank@zappa.com", true, nil)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 400, testUtil.Response.Code, "status code should be 400")
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
		Password:  examplePassword,
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, nil)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 201, testUtil.Response.Code, "status code should be 201")
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
		Password:  "password",
	}

	setExpectationsForUserExistence(testUtil.Mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(testUtil.Mock, exampleUser, arbitraryError)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 500, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"email": "frank@zappa.com",
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
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Email, nil)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler should attach a cookie when request is valid")
	assert.Equal(t, 200, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithInvalidLoginInput(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.NotContains(t, testUtil.Response.HeaderMap, "Set-Cookie", "login handler shouldn't attach a cookie when request is invalid")
	assert.Equal(t, 400, testUtil.Response.Code, "status code should be 400")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithErrorRetrievingUserFromDatabase(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := fmt.Sprintf(`
		{
			"email": "frank@zappa.com",
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
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Email, arbitraryError)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 500, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLoginHandlerWithInvalidPassword(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleInput := `
		{
			"email": "frank@zappa.com",
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
		Password:  examplePassword,
	}

	setExpectationsForUserRetrieval(testUtil.Mock, exampleUser.Email, nil)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 401, testUtil.Response.Code, "status code should be 401")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserLogoutHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	req, err := http.NewRequest("POST", "/logout", nil)
	assert.Nil(t, err)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "logout handler should attach a cookie when request is valid")
	assert.Equal(t, 200, testUtil.Response.Code, "status code should be 200")
}

func TestUserDeletionHandler(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest("DELETE", buildRoute("user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, true, nil)
	setExpectationsForUserDeletion(testUtil.Mock, exampleID, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 200, testUtil.Response.Code, "status code should be 200")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserDeletionHandlerForNonexistentUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest("DELETE", buildRoute("user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, false, nil)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 404, testUtil.Response.Code, "status code should be 404")
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestUserDeletionHandlerWithArbitraryErrorWhenDeletingUser(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	exampleID := uint64(1)
	exampleIDString := strconv.Itoa(int(exampleID))
	req, err := http.NewRequest("DELETE", buildRoute("user", exampleIDString), nil)
	assert.Nil(t, err)

	setExpectationsForUserExistenceByID(testUtil.Mock, exampleIDString, true, nil)
	setExpectationsForUserDeletion(testUtil.Mock, exampleID, arbitraryError)
	testUtil.Router.ServeHTTP(testUtil.Response, req)

	assert.Equal(t, 500, testUtil.Response.Code, "status code should be 500")
	ensureExpectationsWereMet(t, testUtil.Mock)
}
