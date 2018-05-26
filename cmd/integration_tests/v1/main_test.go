package dairytest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func turnResponseBodyIntoString(t *testing.T, res *http.Response) string {
	t.Helper()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return strings.TrimSpace(string(bodyBytes))
}

func assertStatusCode(t *testing.T, resp *http.Response, statusCode int) {
	t.Helper()
	if resp.StatusCode != statusCode {
		assert.Equal(t, statusCode, resp.StatusCode, "status code should be %d", statusCode)
		t.FailNow()
	}
}
