package main

import (
	"database/sql"
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
)

var dummySalt []byte

const (
	examplePassword       = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
	hashedExamplePassword = "$2a$13$hsflIwHM55jooxaTmYahhOO8LdfI.utMBjpHe5Fr311W4PpRxqyXm"
)

func init() {
	os.Setenv("DAIRYSECRET", "do-not-use-secrets-like-this-plz")
	dummySalt = []byte("farts")
}

func TestValidateSessionCookieMiddleware(t *testing.T) {
	t.Parallel()

	handlerWasCalled := false
	exampleHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
	}

	testUtil := setupTestVariablesWithMock(t)

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

	testUtil := setupTestVariablesWithMock(t)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("already existent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("creating an admin user as a non-admin user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserWithUsernameExists", mock.Anything, exampleUser.Username).
			Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleAdminInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusForbidden)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/logout", nil)
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)

		assert.Contains(t, testUtil.Response.HeaderMap, "Set-Cookie", "logout handler should attach a cookie when request is valid")
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid cookie", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, buildRoute("v1", "user", exampleIDString), nil)
		assert.Nil(t, err)
		cookie, err := buildCookieForRequest(t, testUtil.Store, true, false)
		assert.Nil(t, err)
		req.AddCookie(cookie)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusForbidden)
	})

	t.Run("with error deleting user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("UserExists", mock.Anything, exampleID).
			Return(true, nil)
		testUtil.MockDB.On("DeleteUser", mock.Anything, exampleID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		testUtil.MockDB.On("PasswordResetTokenForUserIDExists", mock.Anything, mock.Anything).
			Return(false, nil)
		testUtil.MockDB.On("CreatePasswordResetToken", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving user from db", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with already existent password reset entry", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("PasswordResetTokenForUserIDExists", mock.Anything, mock.Anything).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating reset token", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUserByUsername", mock.Anything, exampleUser.Username).
			Return(exampleUser, nil)
		testUtil.MockDB.On("PasswordResetTokenForUserIDExists", mock.Anything, mock.Anything).
			Return(false, nil)
		testUtil.MockDB.On("CreatePasswordResetToken", mock.Anything, mock.Anything).
			Return(exampleUser.ID, generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, "/password_reset", strings.NewReader(exampleInput))
		assert.Nil(t, err)
		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestPasswordResetValidationHandler(t *testing.T) {
	exampleResetToken := "reset-token"
	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("PasswordResetTokenWithTokenExists", mock.Anything, exampleResetToken).
			Return(true, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent token", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("PasswordResetTokenWithTokenExists", mock.Anything, exampleResetToken).
			Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodHead, fmt.Sprintf("/password_reset/%s", exampleResetToken), nil)
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})
}

func TestUserUpdateHandler(t *testing.T) {
	exampleUserUpdateInput := fmt.Sprintf(`
 		{
 			"username": "captain_beefheart",
 			"current_password": "%s"
 		}
 	`, examplePassword)

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

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, nil)
		testUtil.MockDB.On("UpdateUser", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleGarbageInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with invalid new password", func(*testing.T) {
		exampleInvalidUserUpdateInput := fmt.Sprintf(`
			{
				"new_password": "passwordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				"current_password": "%s"
			}
		`, examplePassword)

		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleInvalidUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("when password does not match expectation", func(*testing.T) {
		exampleInvalidUserUpdateInput := fmt.Sprintf(`
			{
				"username": "captain_beefheart",
				"current_password": "%s"
			}
		`, fmt.Sprintf("%s!", examplePassword))

		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, nil)
		testUtil.MockDB.On("UpdateUser", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleInvalidUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusUnauthorized)
	})

	t.Run("optimal conditions", func(*testing.T) {
		exampleNewPassword := "P@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rd"
		// exampleNewPasswordHashed := "$2a$13$xhhweT6OnsU7l6GyPGdin.YDANUGnFEu7xJQb7eU/zv4KBCiRwWbC"
		exampleUserUpdateInput := fmt.Sprintf(`
			{
				"new_password": "%s",
				"current_password": "%s"
			}
		`, exampleNewPassword, examplePassword)

		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, nil)
		testUtil.MockDB.On("UpdateUser", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error updating user", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetUser", mock.Anything, exampleUser.ID).
			Return(exampleUser, nil)
		testUtil.MockDB.On("UpdateUser", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%d", exampleUser.ID), strings.NewReader(exampleUserUpdateInput))
		assert.Nil(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
