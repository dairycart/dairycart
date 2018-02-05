package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPlugin(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with empty plugin path", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with empty symbol name", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error opening plugin", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error looking up symbol", func(_t *testing.T) {
		_t.Parallel()

	})
}

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

func TestSetConfigDefaults(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestLoadServerConfig(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error reading config file", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestBuildServerConfig(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error building database configuration", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading image storage", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error setting up cookie storage", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestBuildDatabaseFromConfig(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with missing database key", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with empty connection key", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with missing plugin key", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with empty plugin key path", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading plugin", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestLoadDatabasePlugin(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading plugin", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading symbol", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestBuildImageStorerFromConfig(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with missing storage key", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with missing plugin key", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with empty plugin key path", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading plugin", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestLoadImageStoragePlugin(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading plugin", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error loading symbol", func(_t *testing.T) {
		_t.Parallel()

	})
}

func TestInitializeServerComponents(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	t.Run("normal use case", func(_t *testing.T) {
		_t.Parallel()
	})

	t.Run("with error initializing image storage", func(_t *testing.T) {
		_t.Parallel()

	})

	t.Run("with error migrating database", func(_t *testing.T) {
		_t.Parallel()

	})
}
