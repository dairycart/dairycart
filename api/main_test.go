package api

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleSKU          = "example"
	exampleTimeString   = "2016-12-01 12:00:00.000000"
	exampleGarbageInput = `{"things": "stuff"}`
)

////////////////////////////////////////////////////////
//                                                    //
//    This file doesn't actually test main.go, but    //
//     rather contains some small helper functions    //
//     that might be used by all the tests            //
//                                                    //
////////////////////////////////////////////////////////

var arbitraryError error
var exampleTime time.Time
var exampleOlderTime time.Time
var exampleNewerTime time.Time

func init() {
	log.SetOutput(ioutil.Discard)
	arbitraryError = fmt.Errorf("arbitrary error")

	var timeParseErr error
	exampleOlderTime, timeParseErr = time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
	if timeParseErr != nil {
		log.Fatalf("error parsing time")
	}
	exampleTime = exampleOlderTime.Add(30 * (24 * time.Hour))
	exampleNewerTime = exampleTime.Add(30 * (24 * time.Hour))
}

func setupMockRequestsAndMux(db *sql.DB) (*httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	SetupAPIRoutes(m, db)
	return httptest.NewRecorder(), m
}

// sqlmock stuff TODO: give this an ASCII animal

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
