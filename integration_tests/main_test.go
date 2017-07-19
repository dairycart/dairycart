package dairytest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify"

	jsonMinify "github.com/tdewolff/minify/json"
)

const (
	// we can't reliably predict what the `updated_on` or `archived_on` columns could possibly equal,
	// so we strip them out of the body becuase we're bad at programming. The (sort of) plus side to
	// this is that we ensure our timestamps have a particular format (because if they didn't, this
	// function, and as a consequence, the tests, would fail spectacularly).
	// Note that this pattern needs to be run as ungreedy because of the possiblity of prefix and or
	// suffixed commas
	idReplacementPattern          = `(?U)(,?)"(id)":\s?\d+,`
	timeFieldReplacementPattern   = `(?U)(,?)"(created_on|updated_on|archived_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`
	productTimeReplacementPattern = `(?U)(,?)"(available_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`

	existentID     = "1"
	nonexistentID  = "999999999"
	existentSKU    = "t-shirt"
	nonexistentSKU = "nonexistent"

	exampleGarbageInput           = `{"testing_garbage_input": true}`
	expectedBadRequestResponse    = `{"status":400,"message":"Invalid input provided in request body"}`
	expectedInternalErrorResponse = `{"status":500,"message":"Unexpected internal error occurred"}`
)

func loadExpectedResponse(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("expected_responses/%s/%s.json", folder, filename))
	assert.Nil(t, err)
	assert.NotEmpty(t, bodyBytes, "example response file requested is empty and should not be")
	return strings.TrimSpace(string(bodyBytes))
}

func loadExampleInput(t *testing.T, folder string, filename string) string {
	bodyBytes, err := ioutil.ReadFile(fmt.Sprintf("example_inputs/%s/%s.json", folder, filename))
	assert.Nil(t, err)
	assert.NotEmpty(t, bodyBytes, "example input file requested is empty and should not be")
	return strings.TrimSpace(string(bodyBytes))
}

func cleanAPIResponseBody(body string) string {
	idRegex := regexp.MustCompile(idReplacementPattern)
	productTimeRegex := regexp.MustCompile(productTimeReplacementPattern)
	allRowsTimeRegex := regexp.MustCompile(timeFieldReplacementPattern)
	out := allRowsTimeRegex.ReplaceAllString(body, "")
	out = productTimeRegex.ReplaceAllString(out, "")
	out = idRegex.ReplaceAllString(out, "")
	return out
}

func replaceTimeStringsForTests(body string) string {
	re := regexp.MustCompile(timeFieldReplacementPattern)
	return strings.TrimSpace(re.ReplaceAllString(body, ""))
}

func turnResponseBodyIntoString(t *testing.T, res *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	return strings.TrimSpace(string(bodyBytes))
}

func minifyJSON(t *testing.T, jsonBody string) string {
	jsonMinifier := minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
	minified, err := jsonMinifier.String("application/json", jsonBody)
	assert.Nil(t, err)
	return minified
}

func retrieveIDFromResponseBody(body string, t *testing.T) uint64 {
	idContainer := struct {
		ID uint64 `json:"id"`
	}{}
	err := json.Unmarshal([]byte(body), &idContainer)
	assert.Nil(t, err)
	assert.NotEmpty(t, idContainer)
	assert.NotEqual(t, 0, idContainer.ID)

	return idContainer.ID
}
