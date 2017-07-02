package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleSKU          = "example"
	exampleTimeString   = "2016-12-01 12:00:00.000000"
	exampleGarbageInput = `{"things": "stuff"}`
)

///////////////////////////////////////////////////////
//                                                   //
//   These functions don't actually test main.go,    //
//   but rather contains some small helper           //
//   functions that might be used by all the tests   //
//                                                   //
///////////////////////////////////////////////////////

var arbitraryError error
var exampleTime time.Time
var exampleOlderTime time.Time
var exampleNewerTime time.Time

func init() {
	log.SetOutput(ioutil.Discard)
	arbitraryError = errors.New("pineapple on pizza")

	var timeParseErr error
	exampleOlderTime, timeParseErr = time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
	if timeParseErr != nil {
		log.Fatalf("error parsing time")
	}
	exampleTime = exampleOlderTime.Add(30 * (24 * time.Hour))
	exampleNewerTime = exampleTime.Add(30 * (24 * time.Hour))
}

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
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func argsToDriverValues(args []interface{}) []driver.Value {
	rv := []driver.Value{}
	for _, x := range args {
		rv = append(rv, x)
	}
	return rv
}
