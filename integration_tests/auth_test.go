package dairytest

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createUserCreationBody(username string, admin bool) string {
	output := fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "%s",
			"password": "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd",
			"is_admin": %v
		}
	`, username, admin)
	return output
}

func TestUserCreation(t *testing.T) {
	t.Parallel()

	var createdUserID uint64
	testUsername := "test_admin_creation"
	testCreateUser := func(t *testing.T) {
		newUserJSON := createUserCreationBody(testUsername, false)
		resp, err := createNewUser(newUserJSON)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "creating a discount that doesn't exist should respond 201")

		body := turnResponseBodyIntoString(t, resp)
		createdUserID = retrieveIDFromResponseBody(body, t)
	}

	testDeleteUser := func(t *testing.T) {
		resp, err := deleteUser(strconv.Itoa(int(createdUserID)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "trying to delete a discount that exists should respond 200")
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
