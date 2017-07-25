package main

import (
	"database/sql"
	"database/sql/driver"
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

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

const (
	exampleSKU               = "example"
	exampleTimeString        = "2016-12-01 12:00:00.000000"
	exampleGarbageInput      = `{"things": "stuff"}`
	exampleMarshalTimeString = "2016-12-31T12:00:00.000000Z"
)

var (
	exampleFilterStartTime time.Time
	exampleFilterEndTime   time.Time
	defaultQueryFilter     *QueryFilter
	customQueryFilter      *QueryFilter

	arbitraryError   error
	exampleOlderTime time.Time
)

func init() {
	log.SetOutput(ioutil.Discard)
	arbitraryError = errors.New("pineapple on pizza")

	var timeParseErr error
	exampleOlderTime, timeParseErr = time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
	if timeParseErr != nil {
		log.Fatalf("error parsing time")
	}

	defaultQueryFilter = &QueryFilter{
		Page:  1,
		Limit: 25,
	}

	customQueryFilter = &QueryFilter{
		Page:         2,
		Limit:        35,
		CreatedAfter: generateExampleTimeForTests(),
	}
}

///////////////////////////////////////////////////////
//                                                   //
//   These functions don't actually anything, but    //
//   rather contains some small helper functions     //
//   that might be used by all the tests.            //
//                                                   //
///////////////////////////////////////////////////////

type TestUtil struct {
	Response *httptest.ResponseRecorder
	Router   *chi.Mux
	DB       *sqlx.DB
	Mock     sqlmock.Sqlmock
	Store    *sessions.CookieStore
}

func generateExampleTimeForTests() time.Time {
	t, err := time.Parse("2006-01-02 03:04:00.000000", "2016-12-31 12:00:00.000000")
	if err != nil {
		log.Fatalf("error parsing time")
	}
	return t
}

func setExpectationsForRowCount(mock sqlmock.Sqlmock, table string, queryFilter *QueryFilter, count uint64, err error) {
	exampleRows := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(formatQueryForSQLMock(buildCountQuery(table, queryFilter))).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setupTestVariables(t *testing.T) *TestUtil {
	mockDB, mock, err := sqlmock.New()
	db := sqlx.NewDb(mockDB, "postgres")
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	assert.Nil(t, err)

	secret := os.Getenv("DAIRYSECRET")
	if len(secret) < 32 {
		log.Fatalf("Something is up with your app secret: `%s`", secret)
	}
	store := sessions.NewCookieStore([]byte(secret))

	router := chi.NewRouter()
	SetupAPIRoutes(router, db, store)

	return &TestUtil{
		Response: httptest.NewRecorder(),
		Router:   router,
		DB:       db,
		Mock:     mock,
		Store:    store,
	}
}

func formatQueryForSQLMock(query string) string {
	for _, x := range []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"} {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}
	return query
}

func ensureExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func argsToDriverValues(args []interface{}) []driver.Value {
	rv := []driver.Value{}
	for _, x := range args {
		rv = append(rv, x)
	}
	return rv
}

func buildCookieForRequest(store *sessions.CookieStore, authorized bool, admin bool) (*http.Cookie, error) {
	session, err := store.New(&http.Request{}, dairycartCookieName)
	if err != nil {
		return nil, err
	}
	session.Values[sessionUserIDKeyName] = 666
	session.Values[sessionAuthorizedKeyName] = authorized
	session.Values[sessionAdminKeyName] = admin

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, store.Codecs...)
	assert.Nil(t, err)
	cookie := sessions.NewCookie(session.Name(), encoded, session.Options)

	return cookie, nil
}

func attachBadCookieToRequest(req *http.Request) {
	req.Header.Set("Cookie", fmt.Sprintf("%s=this is a bad cookie", dairycartCookieName))
}

///////////////////////////////////////////////////////
//                                                   //
//        These functions actually test things       //
//                                                   //
///////////////////////////////////////////////////////

func TestNullStringMarshalTextReturnsNilIfStringIsInvalid(t *testing.T) {
	t.Parallel()
	example := NullString{sql.NullString{String: "test", Valid: false}}
	alsoNil, err := example.MarshalText()
	assert.Nil(t, err)
	assert.Nil(t, alsoNil)
}

func TestNullTimeMarshalText(t *testing.T) {
	t.Parallel()
	expected := []byte(exampleMarshalTimeString)
	example := NullTime{pq.NullTime{Time: generateExampleTimeForTests(), Valid: true}}
	actual, err := example.MarshalText()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "Marshaled time string should marshal correctly")
}

func TestNullTimeUnmarshalText(t *testing.T) {
	t.Parallel()
	example := []byte(exampleMarshalTimeString)
	nt := NullTime{}
	err := nt.UnmarshalText(example)
	assert.Nil(t, err)
}

func TestParseRawFilterParams(t *testing.T) {
	t.Parallel()
	exampleUnixStartTime := int64(232747200)
	exampleUnixEndTime := int64(232747200 + 10000)

	exampleFilterStartTime := time.Unix(exampleUnixStartTime, 0)
	exampleFilterEndTime := time.Unix(exampleUnixEndTime, 0)

	testSuite := []struct {
		input          string
		expected       *QueryFilter
		failureMessage string
	}{
		{
			input:          "https://test.com/example",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no query params should parse to the default query filter",
		},
		{
			input:          "https://test.com/example?page=1&limit=25",
			expected:       defaultQueryFilter,
			failureMessage: "URL with query params set to the defaults should parse to the default query filter",
		},
		{
			input: "https://test.com/example?page=1&limit=500000",
			expected: &QueryFilter{
				Page:  1,
				Limit: 50,
			},
			failureMessage: "URL with limit param set to high should default to 50",
		},
		{
			input: "https://test.com/example?page=2&limit=40",
			expected: &QueryFilter{
				Page:  2,
				Limit: 40,
			},
			failureMessage: "URL with non-default page and limit params should parse correctly",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v", exampleUnixStartTime),
			expected: &QueryFilter{
				Page:         1,
				Limit:        25,
				UpdatedAfter: exampleFilterStartTime,
			},
			failureMessage: "URL with specified updated_after field should have a non-nil time value set for UpdatedAfter",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_before=%v", exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_before field should have a non-nil time value set for UpdatedBefore",
		},
		{
			input: fmt.Sprintf("https://test.com/example?updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          1,
				Limit:         25,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with specified updated_after and updated_before fields should have a non-nil time value set for both UpdatedAfter and UpdatedBefore",
		},
		{
			input: fmt.Sprintf("https://test.com/example?page=2&limit=35&updated_after=%v&updated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          2,
				Limit:         35,
				UpdatedAfter:  exampleFilterStartTime,
				UpdatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with all relevant filters should have a completely custom QueryFilter value",
		},
		{
			input: fmt.Sprintf("https://test.com/example?page=2&limit=35&created_after=%v&created_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected: &QueryFilter{
				Page:          2,
				Limit:         35,
				CreatedAfter:  exampleFilterStartTime,
				CreatedBefore: exampleFilterEndTime,
			},
			failureMessage: "URL with all relevant filters should have a completely custom QueryFilter value",
		},
		{
			input:          "https://test.com/example?page=0",
			expected:       defaultQueryFilter,
			failureMessage: "URL with page set to zero should default to page 1",
		},
		{
			input:          fmt.Sprintf("https://test.com/example?rage=2&dimit=35&upgraded_after=%v&agitated_before=%v", exampleUnixStartTime, exampleUnixEndTime),
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?page=two",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?limit=eleventy",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?updated_after=my_grandma_died",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?updated_before=my_grandma_lived",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?created_before=the_world_held_its_breath",
			expected:       defaultQueryFilter,
			failureMessage: "URL with no relevant values should parsee to the default query filter",
		},
		{
			input:          "https://test.com/example?created_after=the_world_exhaled",
			expected:       defaultQueryFilter,
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

func TestRowExistsInDBWhenDBThrowsError(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, true, sql.ErrNoRows)
	exists, err := rowExistsInDB(testUtil.DB, skuExistenceQuery, exampleSKU)

	assert.Nil(t, err)
	assert.False(t, exists)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRowExistsInDBForExistingRow(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, true, nil)
	exists, err := rowExistsInDB(testUtil.DB, skuExistenceQuery, exampleSKU)

	assert.Nil(t, err)
	assert.True(t, exists)
	ensureExpectationsWereMet(t, testUtil.Mock)
}

func TestRowExistsInDBForNonexistentRow(t *testing.T) {
	t.Parallel()
	testUtil := setupTestVariables(t)

	setExpectationsForProductExistence(testUtil.Mock, exampleSKU, false, nil)
	exists, err := rowExistsInDB(testUtil.DB, skuExistenceQuery, exampleSKU)

	assert.Nil(t, err)
	assert.False(t, exists)
	ensureExpectationsWereMet(t, testUtil.Mock)
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

	assert.Nil(t, err)
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
