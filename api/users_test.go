package main

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
	userTableHeaders = []string{"id", "first_name", "last_name", "email", "password", "salt", "is_admin", "created_on", "updated_on", "archived_on"}
	exampleUserData = []driver.Value{
		1, "Frank", "Zappa", "frank@zappa.com", hashedExamplePassword, dummySalt, true, exampleTime, nil, nil,
	}
}

func setExpectationsForUserExistence(mock sqlmock.Sqlmock, email string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	mock.ExpectQuery(formatQueryForSQLMock(userExistenceQuery)).
		WithArgs(email).
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
	query := formatQueryForSQLMock(userRetrievalQuery)
	mock.ExpectQuery(query).
		WithArgs(email).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForUserDeletion(mock sqlmock.Sqlmock, email string, err error) {
	mock.ExpectExec(formatQueryForSQLMock(userDeletionQuery)).
		WithArgs(email).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(err)
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
	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	setExpectationsForUserCreation(mock, exampleUser, nil)
	newID, err := createUserInDB(db, exampleUser)

	assert.Nil(t, err)
	assert.Equal(t, exampleUser.ID, newID, "createProductInDB should return the created ID")
	ensureExpectationsWereMet(t, mock)
}

func TestCreateUserInDBWhenErrorOccurs(t *testing.T) {
	t.Parallel()
	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	setExpectationsForUserCreation(mock, exampleUser, arbitraryError)
	_, err := createUserInDB(db, exampleUser)

	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveUserFromDB(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)

	expected := User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  hashedExamplePassword,
		Salt:      dummySalt,
		IsAdmin:   true,
	}

	setExpectationsForUserRetrieval(mock, expected.Email, nil)
	actual, err := retrieveUserFromDB(db, expected.Email)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "expected and actual discounts should match")
	ensureExpectationsWereMet(t, mock)
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
	exampleEmail := "frank@zappa.com"
	db, mock := setupDBForTest(t)
	setExpectationsForUserDeletion(mock, exampleEmail, nil)

	err := archiveUser(db, exampleEmail)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestUserCreationHandler(t *testing.T) {
	t.Parallel()
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
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserExistence(mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(mock, exampleUser, nil)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestUserCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestUserCreationHandlerForAlreadyExistentUserEmail(t *testing.T) {
	t.Parallel()
	exampleInput := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword)
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserExistence(mock, "frank@zappa.com", true, nil)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestUserCreationHandlerWhenErrorEncounteredInsertingIntoDB(t *testing.T) {
	t.Parallel()
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
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  "password",
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserExistence(mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(mock, exampleUser, arbitraryError)

	req, err := http.NewRequest("POST", "/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestUserLoginHandler(t *testing.T) {
	t.Parallel()
	exampleInput := fmt.Sprintf(`
		{
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserRetrieval(mock, exampleUser.Email, nil)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestUserLoginHandlerWithInvalidLoginInput(t *testing.T) {
	t.Parallel()

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestUserLoginHandlerWithErrorRetrievingUserFromDatabase(t *testing.T) {
	t.Parallel()
	exampleInput := fmt.Sprintf(`
		{
			"email": "frank@zappa.com",
			"password": "%s"
		}
	`, examplePassword)

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserRetrieval(mock, exampleUser.Email, arbitraryError)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestUserLoginHandlerWithInvalidPassword(t *testing.T) {
	t.Parallel()
	exampleInput := `
		{
			"email": "frank@zappa.com",
			"password": "password"
		}
	`

	exampleUser := &User{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: exampleTime,
		},
		FirstName: "Frank",
		LastName:  "Zappa",
		Email:     "frank@zappa.com",
		Password:  examplePassword,
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserRetrieval(mock, exampleUser.Email, nil)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 401, res.Code, "status code should be 401")
	ensureExpectationsWereMet(t, mock)
}
