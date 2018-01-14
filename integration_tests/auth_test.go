package dairytest

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
)

func compareUsers(t *testing.T, expected, actual models.User) {
	t.Helper()
	assert.Equal(t, expected.FirstName, actual.FirstName, "expected and actual FirstName should be equal")
	assert.Equal(t, expected.LastName, actual.LastName, "expected and actual LastName should be equal")
	assert.Equal(t, expected.Email, actual.Email, "expected and actual Email should be equal")
	assert.Equal(t, expected.IsAdmin, actual.IsAdmin, "expected and actual IsAdmin should be equal")
}

func TestUserCreationRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		testUsername := "test_user_creation"
		userShouldBeAdmin := false

		expected := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, expected)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		var actual models.User
		unmarshalBody(_t, resp, &actual)
		compareUsers(_t, expected, actual)
	})

	t.Run("with invalid input", func(_t *testing.T) {
		resp, err := createNewUser(exampleGarbageInput, false)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: expectedBadRequestResponse,
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("with invalid password", func(_t *testing.T) {
		testUsername := "test_user_creation_with_invalid_password"
		userShouldBeAdmin := false

		example := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  "invalid",
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, example)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "password must be at least 64 characters",
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("with duplicate username", func(_t *testing.T) {
		testUsername := "test_duplicate_user_creation"
		userShouldBeAdmin := false

		example := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, example)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		// create the user again to trigger the error
		resp, err = createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusBadRequest, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "username already taken",
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("creating admin user as admin", func(_t *testing.T) {
		testUsername := "test_admin_user_creation"
		userShouldBeAdmin := true

		expected := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, expected)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		var actual models.User
		unmarshalBody(_t, resp, &actual)
		compareUsers(_t, expected, actual)
	})

	t.Run("creating admin user as regular user", func(_t *testing.T) {
		testUsername := "test_admin_user_creation_without_creds"
		userShouldBeAdmin := true

		example := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, example)
		resp, err := createNewUser(exampleInput, !userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusForbidden, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "User is not authorized to create admin users",
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestUserDeletionRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		testUsername := "test_user_deletion"
		userShouldBeAdmin := true

		expected := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, expected)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		var createdUser models.User
		unmarshalBody(_t, resp, &createdUser)

		resp, err = deleteUser(createdUser.ID, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		var actual models.User
		unmarshalBody(_t, resp, &actual)
		assert.NotNil(t, actual.ArchivedOn)
	})

	t.Run("for nonexistent user", func(_t *testing.T) {
		resp, err := deleteUser(nonexistentID, true)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The user you were looking for (username '%d') does not exist", nonexistentID),
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("as non-admin user", func(_t *testing.T) {
		testUsername := "test_user_deletion_as_non_admin_user"
		userShouldBeAdmin := false

		example := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, example)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		var createdUser models.User
		unmarshalBody(_t, resp, &createdUser)

		resp, err = deleteUser(createdUser.ID, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusForbidden, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "User is not authorized to delete users",
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestUserLoginRoute(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		testUsername := "test_user_login"
		userShouldBeAdmin := false

		expected := models.UserCreationInput{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, expected)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		testUserCookie := resp.Cookies()[0]
		resp, err = logoutUser(testUserCookie)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		resp, err = loginUser(expected.Username, expected.Password)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)
	})

	t.Run("for nonexistent user", func(_t *testing.T) {
		testUsername := "nonexistent_user"
		resp, err := loginUser(testUsername, validPassword)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusNotFound, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("The user you were looking for (username '%s') does not exist", testUsername),
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})

	t.Run("with invalid input", func(_t *testing.T) {
		url := buildVersionlessPath("login")
		body := strings.NewReader(exampleGarbageInput)
		req, err := http.NewRequest(http.MethodPost, url, body)
		assert.Nil(t, err)

		resp, err := requester.Do(req)
		assert.Nil(t, err)
		assertStatusCode(_t, resp, http.StatusBadRequest)
	})

	t.Run("with invalid password", func(_t *testing.T) {
		testUsername := "test_user_login_with_invalid_password"
		userShouldBeAdmin := false

		example := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, example)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		testUserCookie := resp.Cookies()[0]
		resp, err = logoutUser(testUserCookie)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)

		resp, err = loginUser(example.Username, "password")
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusUnauthorized, resp.StatusCode)

		expected := models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid email and/or password",
		}
		var actual models.ErrorResponse
		unmarshalBody(_t, resp, &actual)
		assert.Equal(_t, expected, actual)
	})
}

func TestUserLogout(t *testing.T) {
	t.Parallel()

	t.Run("normal usage", func(_t *testing.T) {
		testUsername := "test_user_logout"
		userShouldBeAdmin := false

		expected := models.User{
			FirstName: "Frank",
			LastName:  "Zappa",
			Email:     "frank@zappa.com",
			Username:  testUsername,
			Password:  validPassword,
			IsAdmin:   userShouldBeAdmin,
		}
		exampleInput := createJSONBody(t, expected)
		resp, err := createNewUser(exampleInput, userShouldBeAdmin)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusCreated, resp.StatusCode)

		testUserCookie := resp.Cookies()[0]
		resp, err = logoutUser(testUserCookie)
		assert.NoError(_t, err)
		assert.Equal(_t, http.StatusOK, resp.StatusCode)
	})
}
