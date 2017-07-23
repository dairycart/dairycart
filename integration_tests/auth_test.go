package dairytest

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createUserCreationBody(username string, password string, admin bool) string {
	output := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "%s",
			"password": "%s",
			"is_admin": %v
		}
	`, username, password, admin)
	return output
}

func TestUserCreation(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_user_creation"
	userShouldBeAdmin := false
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), true)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a user that exists should respond 200")
	}

	subtests := []subtest{
		subtest{
			Message: "create user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestUserCreationWithInvalidPassword(t *testing.T) {
	t.Parallel()

	testUsername := "test_bad_password"
	testAwfulPassword := "password"
	userShouldBeAdmin := false
	newUserJSON := createUserCreationBody(testUsername, testAwfulPassword, userShouldBeAdmin)
	resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a user with an invalid password should respond 400")

	// FIXME: this error response isn't super optimal
	expected := `{"status":400,"message":"Key: 'UserCreationInput.Password' Error:Field validation for 'Password' failed on the 'gte' tag"}`
	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expected, actual, "response to invalid password should equal expectation")
}

func TestUserCreationWithInvalidCreationBody(t *testing.T) {
	t.Parallel()

	userShouldBeAdmin := false
	resp, err := createNewUser(exampleGarbageInput, userShouldBeAdmin)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a user with an invalid password should respond 400")

	expected := `{"status":400,"message":"Invalid input provided in request body"}`
	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expected, actual, "response to invalid password should equal expectation")
}

func TestAdminUserCreation(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_admin_user_creation"
	userShouldBeAdmin := true
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), true)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a user that exists should respond 200")
	}

	subtests := []subtest{
		subtest{
			Message: "create admin user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created admin user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestAdminUserCreationFailsWithoutAdminCredentials(t *testing.T) {
	t.Parallel()
	testUsername := "test_admin_user_creation_without_creds"
	newUserJSON := createUserCreationBody(testUsername, validPassword, true)
	resp, err := createNewUser(newUserJSON, false)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "creating an admin user without valid credentials should respond 403")
}

func TestUserCreationForAlreadyExistentUsername(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_duplicate_user_creation"
	userShouldBeAdmin := false
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testCreateUserAgain := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a user that already exists should respond 400")
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), true)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a user that exists should respond 200")
	}

	subtests := []subtest{
		subtest{
			Message: "create user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "create user again",
			Test:    testCreateUserAgain,
		},
		subtest{
			Message: "delete created user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestUserDeletion(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_user_deletion"
	userShouldBeAdmin := false
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		createdUserID = retrieveIDFromResponseBody(turnResponseBodyIntoString(t, resp), t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), true)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a user that exists should respond 200")
		body := turnResponseBodyIntoString(t, resp)
		assert.Empty(t, body)
	}

	subtests := []subtest{
		subtest{
			Message: "create user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestUserDeletionForNonexistentUser(t *testing.T) {
	t.Parallel()

	resp, err := deleteUser(nonexistentID, true)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "trying to delete a user that doesn't exist should respond 404")

	expected := `{"status":404,"message":"The user you were looking for (username '999999999') does not exist"}`
	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expected, actual, "anticipated response body should match")
}

func TestUserDeletionAsRegularUser(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_user_deletion_as_regular_user"
	userShouldBeAdmin := false
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		createdUserID = retrieveIDFromResponseBody(turnResponseBodyIntoString(t, resp), t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), userShouldBeAdmin)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "trying to delete an admin user as a regular user should respond 403")

		expected := `{"status":403,"message":"User is not authorized to delete users"}`
		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, expected, actual, "anticipated response body should match")
	}

	subtests := []subtest{
		subtest{
			Message: "create user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestAdminUserDeletion(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_admin_user_deletion"
	userShouldBeAdmin := true
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		createdUserID = retrieveIDFromResponseBody(turnResponseBodyIntoString(t, resp), t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), true)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete an admin user that exists should respond 200")
		body := turnResponseBodyIntoString(t, resp)
		assert.Empty(t, body)
	}

	subtests := []subtest{
		subtest{
			Message: "create admin user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created admin user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestAdminUserDeletionAsRegularUser(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_admin_user_deletion_as_non_admin"
	userShouldBeAdmin := true
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, userShouldBeAdmin)
		resp, err := createNewUser(newUserJSON, userShouldBeAdmin)
		assert.Nil(t, err)
		createdUserID = retrieveIDFromResponseBody(turnResponseBodyIntoString(t, resp), t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)), false)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "trying to delete an admin user as a regular user should respond 403")

		expected := `{"status":403,"message":"User is not authorized to delete users"}`
		actual := turnResponseBodyIntoString(t, resp)
		assert.Equal(t, expected, actual, "anticipated response body should match")
	}

	subtests := []subtest{
		subtest{
			Message: "create admin user",
			Test:    testCreateUser,
		},
		subtest{
			Message: "delete created admin user",
			Test:    testDeleteUser,
		},
	}
	runSubtestSuite(t, subtests)
}
