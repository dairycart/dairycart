package main

import (
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

const (
	examplePassword = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

// Here's where we test the awful parts of the code (which ones, lol)

type arbitraryPanicker struct {
	FatalWasCalled bool
}

func (a *arbitraryPanicker) Fatal(...interface{}) {
	a.FatalWasCalled = true
}

func TestPanicUponSaltGenerationError(t *testing.T) {
	t.Parallel()
	a := &arbitraryPanicker{}
	panicUponSaltGenerationError(arbitraryError, a)
	assert.True(t, a.FatalWasCalled)
}

type arbitraryChecker struct {
	Returning bool
}

func (a *arbitraryChecker) IsErroneous(err error) bool {
	return a.Returning
}

func TestCheckIfBcryptErrorIsValid(t *testing.T) {
	t.Parallel()
	ch := &arbitraryChecker{
		Returning: true,
	}
	actual := checkIfBcryptErrorIsValid(arbitraryError, ch)
	assert.True(t, actual)
}

// Begin normal testing things

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

func TestCreateUserFromInput(t *testing.T) {
	t.Parallel()
	exampleUserInput := &UserCreationInput{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  "password",
		IsAdmin:   true,
	}
	expected := &User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		IsAdmin:   true,
	}

	actual, err := createUserFromInput(exampleUserInput, defaultChecker{})
	assert.Nil(t, err)

	assert.Equal(t, expected.FirstName, actual.FirstName, "FirstName fields should match")
	assert.Equal(t, expected.LastName, actual.LastName, "LastName fields should match")
	assert.Equal(t, expected.Email, actual.Email, "Email fields should match")
	assert.Equal(t, expected.IsAdmin, actual.IsAdmin, "IsAdmin fields should match")
	assert.NotEqual(t, expected.Password, actual.Password, "Generated User password should not have the same password as the user input")
	assert.Equal(t, saltSize, len(actual.Salt), fmt.Sprintf("Generated salt should be %d bytes large", saltSize))
}

func TestCreateUserFromInputWithError(t *testing.T) {
	t.Parallel()
	exampleUserInput := &UserCreationInput{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		IsAdmin:   true,
	}

	ch := &arbitraryChecker{
		Returning: true,
	}
	actual, _ := createUserFromInput(exampleUserInput, ch)
	assert.Nil(t, actual)
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
	assert.Nil(t, bcrypt.CompareHashAndPassword(actual, saltedPass))
}

func TestValidateUserCreationInputWithValidInput(t *testing.T) {
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
		Password:  []byte("password"),
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
		Password:  []byte("password"),
	}

	db, mock := setupDBForTest(t)
	setExpectationsForUserCreation(mock, exampleUser, arbitraryError)
	_, err := createUserInDB(db, exampleUser)

	assert.NotNil(t, err)
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
		Password:  []byte("password"),
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserExistence(mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(mock, exampleUser, nil)

	req, err := http.NewRequest("POST", "/v1/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code, "status code should be 201")
	ensureExpectationsWereMet(t, mock)
}

func TestUserCreationHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	req, err := http.NewRequest("POST", "/v1/user", strings.NewReader(exampleGarbageInput))
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

	req, err := http.NewRequest("POST", "/v1/user", strings.NewReader(exampleInput))
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
		Password:  []byte("password"),
	}

	db, mock := setupDBForTest(t)
	res, router := setupMockRequestsAndMux(db)

	setExpectationsForUserExistence(mock, exampleUser.Email, false, nil)
	setExpectationsForUserCreation(mock, exampleUser, arbitraryError)

	req, err := http.NewRequest("POST", "/v1/user", strings.NewReader(exampleInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
