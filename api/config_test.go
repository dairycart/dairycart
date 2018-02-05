package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupCookieStorage(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()
		cs, err := setupCookieStorage("arbitrarily long secret for testing purposes")
		assert.NoError(_t, err)
		assert.NotNil(_t, cs)
	})

	t.Run("with short secret", func(_t *testing.T) {
		_t.Parallel()
		cs, err := setupCookieStorage("lol")
		assert.Error(_t, err)
		assert.Nil(_t, cs)
	})
}
