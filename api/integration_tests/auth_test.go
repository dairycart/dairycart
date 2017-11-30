package dairytest

import (
	"testing"
	//
	// "github.com/stretchr/testify/assert"
)

func TestUserCreationRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create user
		// delete user
	})

	t.Run("with invalid input", func(*testing.T) {

	})

	t.Run("with invalid password", func(*testing.T) {

	})

	t.Run("with duplicate username", func(*testing.T) {

	})

	t.Run("creating admin user as admin", func(*testing.T) {

	})

	t.Run("creating admin user as regular user", func(*testing.T) {

	})
}

func TestUserDeletionRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create user
		// delete user
	})

	t.Run("for nonexistent user", func(*testing.T) {

	})

	t.Run("as non-admin user", func(*testing.T) {

	})

	t.Run("delete admin user", func(*testing.T) {
		// create admin user
		// delete admin user
	})
}

func TestUserLoginRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create user
		// delete user
	})

	t.Run("for nonexistent user", func(*testing.T) {

	})

	t.Run("with invalid input", func(*testing.T) {

	})

	t.Run("with invalid password", func(*testing.T) {

	})
}

func TestUserLogout(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create user
		// delete user
	})
}
