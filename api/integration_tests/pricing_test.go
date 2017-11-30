package dairytest

import (
	"testing"
	//
	// "github.com/stretchr/testify/assert"
)

func TestDiscountRetrievalRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {

	})

	t.Run("for nonexistent discount", func(*testing.T) {

	})

}

func TestDiscountListRoute(t *testing.T) {
	// t.Parallel()

	t.Run("default filter", func(*testing.T) {

	})

	t.Run("custom filter", func(*testing.T) {

	})
}

func TestDiscountCreationRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create discount
		// delete discount
	})

	t.Run("with invalid input", func(*testing.T) {

	})
}

func TestDiscountDeletionRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create discount
		// delete discount
	})

	t.Run("for nonexistent discount", func(*testing.T) {

	})
}

func TestDiscountUpdateRoute(t *testing.T) {
	// t.Parallel()

	t.Run("normal usage", func(*testing.T) {
		// create discount
		// update discount
	})

	t.Run("with invalid input", func(*testing.T) {

	})

	t.Run("for nonexistent discount", func(*testing.T) {

	})

}
