package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	// local dependencies
	"github.com/dairycart/dairycart/api/storage/mock"
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleSKU          = "example"
	exampleGarbageInput = `{"things": "stuff"}`
)

func init() {
	log.SetOutput(ioutil.Discard)
}

///////////////////////////////////////////////////////
//                                                   //
//   These functions don't actually anything, but    //
//   rather contains some small helper functions     //
//   that might be used by all the tests.            //
//                                                   //
///////////////////////////////////////////////////////

// TODO: Rename much of these fields as well as this entire struct
type TestUtil struct {
	Response *httptest.ResponseRecorder
	Router   *chi.Mux
	PlainDB  *sql.DB
	Mock     sqlmock.Sqlmock
	MockDB   *dairymock.MockDB
	Store    *sessions.CookieStore
}

func generateExampleTimeForTests() time.Time {
	out, err := time.Parse("2006-01-02 03:04:00.000000", "2016-12-31 12:00:00.000000")
	if err != nil {
		log.Fatalf("error parsing time")
	}
	return out
}

func genereateDefaultQueryFilter() *models.QueryFilter {
	qf := &models.QueryFilter{
		Page:  1,
		Limit: 25,
	}
	return qf
}

func generateArbitraryError() error {
	return errors.New("pineapple on pizza")
}

func setupTestVariablesWithMock(t *testing.T) *TestUtil {
	t.Helper()
	mockDB, mock, _ := sqlmock.New()
	return &TestUtil{
		Response: httptest.NewRecorder(),
		Router:   chi.NewRouter(),
		PlainDB:  mockDB,
		Mock:     mock,
		MockDB:   &dairymock.MockDB{},
		Store:    sessions.NewCookieStore([]byte(os.Getenv("DAIRYSECRET"))),
	}
}

func ensureExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func buildCookieForRequest(t *testing.T, store *sessions.CookieStore, authorized bool, admin bool) (*http.Cookie, error) {
	t.Helper()
	session, err := store.New(&http.Request{}, dairycartCookieName)
	if err != nil {
		return nil, err
	}
	session.Values[sessionUserIDKeyName] = 666
	session.Values[sessionAuthorizedKeyName] = authorized
	session.Values[sessionAdminKeyName] = admin

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, store.Codecs...)
	assert.NoError(t, err)
	cookie := sessions.NewCookie(session.Name(), encoded, session.Options)

	return cookie, nil
}

func attachBadCookieToRequest(req *http.Request) {
	req.Header.Set("Cookie", fmt.Sprintf("%s=this is a bad cookie", dairycartCookieName))
}

func assertStatusCode(t *testing.T, testUtil *TestUtil, statusCode int) {
	t.Helper()
	assert.Equal(t, statusCode, testUtil.Response.Code, "status code should be %d", statusCode)
}

///////////////////////////////////////////////////////
//                                                   //
//        These functions actually test things       //
//                                                   //
///////////////////////////////////////////////////////

func TestParseRawFilterParams(t *testing.T) {
	t.Parallel()
	exampleUnixStartTime := int64(232747200)
	exampleUnixEndTime := int64(232747200 + 10000)

	exampleFilterStartTime := time.Unix(exampleUnixStartTime, 0)
	exampleFilterEndTime := time.Unix(exampleUnixEndTime, 0)

	testSuite := []struct {
		input          string
		expected       *models.QueryFilter
		failureMessage string
	}{
		{
			input:          "https://test.com/example",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no query params should parse to the default query filter",
		},
		{
			input:          "https://test.com/example?page=1&limit=25",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with query params set to the defaults should parse to the default query filter",
		},
		{
			input: "https://test.com/example?page=1&limit=500000",
			expected: &models.QueryFilter{
				Page:  1,
				Limit: 50,
			},
			failureMessage: "URL with limit param set to high should default to 50",
		},
		{
			input: "https://test.com/example?page=2&limit=40",
			expected: &models.QueryFilter{
				Page:  2,
				Limit: 40,
			},
			failureMessage: "URL with non-default page and limit params should parse correctly",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v", exampleUnixStartTime),
			expected: &models.QueryFilter{
				Page:         1,
				Limit:        25,
				UpdatedAfter: exampleFilterStartTime,
			},
			failureMessage: "URL with specified updated_after field should have a non-nil time value set for UpdatedAfter",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_before=%v", exampleUnixEndTime),
			expected: &models.QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_before field should have a non-nil time value set for UpdatedBefore",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &models.QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_after and updated_before fields should have a non-nil time value set for both UpdatedAfter and UpdatedBefore",
		},
		{
			input: fmt.Sprintf("https://test.com/example?page=2&limit=35&updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &models.QueryFilter{
				Page:          2,
				Limit:         35,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with all relevant filters should have a completely custom QueryFilter value",
		},
		{
			input: fmt.Sprintf("https://test.com/example?page=2&limit=35&created_after=%v&created_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &models.QueryFilter{
				Page:          2,
				Limit:         35,
				CreatedAfter:  exampleFilterStartTime,
				CreatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with all relevant filters should have a completely custom QueryFilter value",
		},
		{
			input:          "https://test.com/example?page=0",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with page set to zero should default to page 1",
		},
		{
			input:          fmt.Sprintf("https://test.com/example?rage=2&dimit=35&upgraded_after=%v&agitated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?page=two",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?limit=eleventy",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?updated_after=my_grandma_died",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?updated_before=my_grandma_lived",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?created_before=the_world_held_its_breath",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?created_after=the_world_exhaled",
			expected:       genereateDefaultQueryFilter(),
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
	}

	for _, test := range testSuite {
		earl, err := url.Parse(test.input)
		if err != nil {
			log.Fatal(err)
		}
		actual := parseRawFilterParams(earl.Query())
		assert.Equal(t, test.expected, actual, test.failureMessage)
	}

}

func TestRestrictedStringIsValid(t *testing.T) {
	testCases := []struct {
		Input        string
		ShouldPass   bool
		ErrorMessage string
	}{
		{
			Input:        "this_string_is_fine",
			ShouldPass:   true,
			ErrorMessage: "ordinary sku example should pass",
		},
		{
			ShouldPass:   false,
			ErrorMessage: "empty or uninitialized strings should fail",
		},
		{
			Input:        "this string has spaces",
			ShouldPass:   false,
			ErrorMessage: "database values should not have spaces",
		},
		{
			Input:        "this_entry_is_just_way_way_way_way_way_way_way_way_way_way_too_long",
			ShouldPass:   false,
			ErrorMessage: "nothing longer than fifty characters",
		},
		{
			Input:        "ⓖⓞⓞⓕⓨ ⓣⓔⓧⓣ ⓝⓞⓣ ⓐⓛⓛⓞⓦⓔⓓ",
			ShouldPass:   false,
			ErrorMessage: "goofy text should not be allowed",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.ShouldPass, restrictedStringIsValid(test.Input), test.ErrorMessage)
	}
}

func TestRespondThatRowDoesNotExist(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	respondThatRowDoesNotExist(req, w, "item", "something")

	actual := strings.TrimSpace(w.Body.String())
	expected := `{"status":404,"message":"The item you were looking for (identified by 'something') does not exist"}`

	assert.Equal(t, expected, actual, "response should indicate the row was not found")
	assert.Equal(t, http.StatusNotFound, w.Code, "status code should be 404")
}

func TestNotifyOfInvalidRequestBody(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	notifyOfInvalidRequestBody(w, errors.New("test"))

	actual := strings.TrimSpace(w.Body.String())
	expected := `{"status":400,"message":"test"}`

	assert.Equal(t, expected, actual, "response should indicate the request body was invalid")
	assert.Equal(t, http.StatusBadRequest, w.Code, "status code should be 404")
}

func TestNotifyOfInternalIssue(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()

	notifyOfInternalIssue(w, errors.New("test"), "do a thing")

	actual := strings.TrimSpace(w.Body.String())
	expected := `{"status":500,"message":"Unexpected internal error occurred"}`

	assert.Equal(t, expected, actual, "response should indicate their was an internal error")
	assert.Equal(t, http.StatusInternalServerError, w.Code, "status code should be 404")
}

func TestValidateRequestInput(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(fmt.Sprintf(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "%s"
		}
	`, examplePassword))

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual := &UserCreationInput{}
	err := validateRequestInput(req, actual)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestValidateRequestInputWithAwfulpassword(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(`
		{
			"first_name": "Frank",
			"last_name": "Zappa",
			"email": "frank@zappa.com",
			"username": "frankzappa",
			"password": "password"
		}
	`)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual := &UserCreationInput{}
	err := validateRequestInput(req, actual)

	assert.NotNil(t, err)
}

func TestValidateRequestInputWithGarbageInput(t *testing.T) {
	t.Parallel()

	exampleInput := strings.NewReader(exampleGarbageInput)
	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual := &UserCreationInput{}
	err := validateRequestInput(req, actual)

	assert.NotNil(t, err)
}

func TestValidateRequestInputWithCompletelyGarbageInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "http://example.com", nil)
	actual := &UserCreationInput{}
	err := validateRequestInput(req, actual)

	assert.NotNil(t, err)
}
