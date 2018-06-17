package dairytest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	existentID     = uint64(1)
	nonexistentID  = uint64(999999999)
	existentSKU    = "t-shirt-small-red"
	nonexistentSKU = "nonexistent"

	exampleGarbageInput           = `{"testing_garbage_input": true}`
	expectedBadRequestResponse    = `Invalid input provided in request body`
	expectedInternalErrorResponse = `Unexpected internal error occurred`
)

func assertStatusCode(t *testing.T, resp *http.Response, statusCode int) {
	t.Helper()
	if resp.StatusCode != statusCode {
		assert.Equal(t, statusCode, resp.StatusCode, "status code should be %d", statusCode)
		t.FailNow()
	}
}
