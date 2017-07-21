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
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, false)
		resp, err := createNewUser(newUserJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)))
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
	newUserJSON := createUserCreationBody(testUsername, testAwfulPassword, true)
	resp, err := createNewUser(newUserJSON)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a user with an invalid password should respond 400")

	// FIXME: this error response isn't super optimal
	expected := `{"status":400,"message":"Key: 'UserCreationInput.Password' Error:Field validation for 'Password' failed on the 'gte' tag"}`
	actual := turnResponseBodyIntoString(t, resp)
	assert.Equal(t, expected, actual, "response to invalid password should equal expectation")
}

func TestUserCreationWithInvalidCreationBody(t *testing.T) {
	t.Parallel()

	resp, err := createNewUser(exampleGarbageInput)
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
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, true)
		resp, err := createNewUser(newUserJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)))
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

func TestUserCreationForAlreadyExistentUsername(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_duplicate_user_creation"
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, false)
		resp, err := createNewUser(newUserJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a user that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testCreateUserAgain := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, validPassword, false)
		resp, err := createNewUser(newUserJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "creating a user that already exists should respond 400")
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)))
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
