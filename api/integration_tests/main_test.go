package dairytest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/minify"
	jsonMinify "github.com/tdewolff/minify/json"
)

const (
	// we can't reliably predict what the `updated_on` or `archived_on` columns could possibly equal,
	// so we strip them out of the body because we're bad at programming. The (sort of) plus side to
	// this is that we ensure our timestamps have a particular format (because if they didn't, this
	// function, and as a consequence, the tests, would fail spectacularly).
	// Note that this pattern needs to be run as ungreedy because of the possiblity of prefix and or
	// suffixed commas
	idReplacementPattern                = `(?U)(,?)"(id|product_option_id|product_root_id)":\s?\d+,`
	productTimeReplacementPattern       = `(?U)(,?)"(available_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`
	timeFieldReplacementPattern         = `(?U)(,?)"(created_on|updated_on|archived_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`
	discountTimeFieldReplacementPattern = `(?U)(,?)"(starts_on|expires_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?(Z|\+\d{2}:\d{2}))?"(,?)`

	existentID     = uint64(1)
	nonexistentID  = uint64(999999999)
	existentSKU    = "t-shirt-small-red"
	nonexistentSKU = "nonexistent"

	exampleGarbageInput           = `{"testing_garbage_input": true}`
	expectedBadRequestResponse    = `Invalid input provided in request body`
	expectedInternalErrorResponse = `Unexpected internal error occurred`
)

type listResponse struct {
	Count uint64        `json:"count"`
	Limit uint8         `json:"limit"`
	Page  uint64        `json:"page"`
	Data  []interface{} `json:"data"`
}

type subtest struct {
	Message string
	Test    func(t *testing.T)
}

func cleanAPIResponseBody(body string) string {
	idRegex := regexp.MustCompile(idReplacementPattern)
	allRowsTimeRegex := regexp.MustCompile(timeFieldReplacementPattern)
	productTimeRegex := regexp.MustCompile(productTimeReplacementPattern)
	discountTimeRegex := regexp.MustCompile(discountTimeFieldReplacementPattern)
	out := allRowsTimeRegex.ReplaceAllString(body, "")
	out = productTimeRegex.ReplaceAllString(out, "")
	out = discountTimeRegex.ReplaceAllString(out, "")
	out = idRegex.ReplaceAllString(out, "")
	return out
}

func replaceTimeStringsForTests(body string) string {
	genericTimeRegex := regexp.MustCompile(timeFieldReplacementPattern)
	productTimeRegex := regexp.MustCompile(productTimeReplacementPattern)
	discountTimeRegex := regexp.MustCompile(discountTimeFieldReplacementPattern)
	out := strings.TrimSpace(genericTimeRegex.ReplaceAllString(body, ""))
	out = productTimeRegex.ReplaceAllString(out, "")
	out = discountTimeRegex.ReplaceAllString(out, "")
	return out
}

func turnResponseBodyIntoString(t *testing.T, res *http.Response) string {
	t.Helper()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return strings.TrimSpace(string(bodyBytes))
}

func minifyJSON(t *testing.T, jsonBody string) string {
	t.Helper()
	jsonMinifier := minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
	minified, err := jsonMinifier.String("application/json", jsonBody)
	require.Nil(t, err)
	return minified
}

func retrieveProductIDFromResponseBody(t *testing.T, body string) uint64 {
	t.Helper()
	idContainer := struct {
		Products []struct {
			ID uint64 `json:"id"`
		} `json:"products"`
	}{}
	err := json.Unmarshal([]byte(body), &idContainer)
	assert.NoError(t, err)
	assert.NotEmpty(t, idContainer)
	assert.NotZero(t, idContainer.Products[0].ID, fmt.Sprintf("ID should not be zero, body is:\n%s", body))

	return idContainer.Products[0].ID
}

func retrieveIDFromResponseBody(t *testing.T, body string) uint64 {
	t.Helper()
	idContainer := struct {
		ID uint64 `json:"id"`
	}{}
	err := json.Unmarshal([]byte(body), &idContainer)
	assert.NoError(t, err)
	assert.NotEmpty(t, idContainer)
	assert.NotZero(t, idContainer.ID, fmt.Sprintf("ID should not be zero, body is:\n%s", body))

	return idContainer.ID
}

func retrieveProductRootIDFromResponseBody(t *testing.T, body string) uint64 {
	t.Helper()
	idContainer := struct {
		ID uint64 `json:"product_root_id"`
	}{}
	err := json.Unmarshal([]byte(body), &idContainer)
	assert.NoError(t, err)
	assert.NotEmpty(t, idContainer)
	assert.NotZero(t, idContainer.ID, fmt.Sprintf("ID should not be zero, body is:\n%s", body))

	return idContainer.ID
}

func parseResponseIntoStruct(t *testing.T, body string) listResponse {
	t.Helper()
	lr := listResponse{}
	err := json.Unmarshal([]byte(body), &lr)
	assert.NoError(t, err)
	assert.NotEmpty(t, lr)

	return lr
}

// func runSubtestSuite(t *testing.T, testFuncs map[string]func(t *testing.T)) {
func runSubtestSuite(t *testing.T, tests []subtest) {
	t.Helper()
	testPassed := true
	for _, test := range tests {
		if !testPassed {
			t.FailNow()
		}
		testPassed = t.Run(test.Message, test.Test)
	}
}

func assertStatusCode(t *testing.T, resp *http.Response, statusCode int) {
	t.Helper()
	assert.Equal(t, statusCode, resp.StatusCode, "status code should be %d", statusCode)
}
