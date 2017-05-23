package api

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondThatRowDoesNotExist(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	respondThatRowDoesNotExist(req, w, "item", "field", "something")

	assert.Equal(t, w.Body.String(), "No item with the field 'something' found\n", "response should indicate the row was not found")
	assert.Equal(t, w.Code, 404, "status code should be 404")
}

func TestNotifyOfInvalidRequestBody(t *testing.T) {
	w := httptest.NewRecorder()
	notifyOfInvalidRequestBody(w, errors.New("test"))

	assert.Equal(t, w.Body.String(), "Invalid request body\n", "response should indicate the request body was invalid")
	assert.Equal(t, w.Code, 400, "status code should be 404")
}

func TestNotifyOfInternalIssue(t *testing.T) {
	w := httptest.NewRecorder()

	notifyOfInternalIssue(w, errors.New("test"), "do a thing")

	assert.Equal(t, w.Body.String(), "Unexpected internal error\n", "response should indicate their was an internal error")
	assert.Equal(t, w.Code, 500, "status code should be 404")
}
