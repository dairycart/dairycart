package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const (
	examplePassword = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := generateSalt()
	assert.Nil(t, err)
	assert.Equal(t, saltSize, len(salt), fmt.Sprintf("Generated salt should be %d bytes large", saltSize))
}

func TestSaltAndHashPassword(t *testing.T) {
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
