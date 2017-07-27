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
	// so we strip them out of the body because we're bad at programming. The (sort of) plus side to
	// this is that we ensure our timestamps have a particular format (because if they didn't, this
	// function, and as a consequence, the tests, would fail spectacularly).
	// Note that this pattern needs to be run as ungreedy because of the possiblity of prefix and or
	// suffixed commas
	idReplacementPattern                = `(?U)(,?)"(id)":\s?\d+,`
	productTimeReplacementPattern       = `(?U)(,?)"(available_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`
	timeFieldReplacementPattern         = `(?U)(,?)"(created_on|updated_on|archived_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?Z)?"(,?)`
	discountTimeFieldReplacementPattern = `(?U)(,?)"(starts_on|expires_on)":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?(Z|\+\d{2}:\d{2}))?"(,?)`

	existentID     = "1"
	nonexistentID  = "999999999"
	existentSKU    = "t-shirt-small-red"
	nonexistentSKU = "nonexistent"

	exampleGarbageInput           = `{"testing_garbage_input": true}`
	expectedBadRequestResponse    = `{"status":400,"message":"Invalid input provided in request body"}`
	expectedInternalErrorResponse = `{"status":500,"message":"Unexpected internal error occurred"}`
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
	assert.NotZero(t, idContainer.ID, fmt.Sprintf("ID should not be zero, body is:\n%s", body))

	return idContainer.ID
}

func retrieveProductRootIDFromResponseBody(body string, t *testing.T) uint64 {
	idContainer := struct {
		ID uint64 `json:"product_root_id"`
	}{}
	err := json.Unmarshal([]byte(body), &idContainer)
	assert.Nil(t, err)
	assert.NotEmpty(t, idContainer)
	assert.NotZero(t, idContainer.ID, fmt.Sprintf("ID should not be zero, body is:\n%s", body))

	return idContainer.ID
}

func parseResponseIntoStruct(body string, t *testing.T) listResponse {
	lr := listResponse{}
	err := json.Unmarshal([]byte(body), &lr)
	assert.Nil(t, err)
	assert.NotEmpty(t, lr)

	return lr
}

// func runSubtestSuite(t *testing.T, testFuncs map[string]func(t *testing.T)) {
func runSubtestSuite(t *testing.T, tests []subtest) {
	testPassed := true
	for _, test := range tests {
		if !testPassed {
			t.FailNow()
		}
		testPassed = t.Run(test.Message, test.Test)
	}
}
