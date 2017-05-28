package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/gorilla/mux"
)

const (
	exampleSKU        = "example"
	exampleTimeString = "2017-01-01 12:00:00.000000"
)

////////////////////////////////////////////////////////
//                                                    //
//    This file doesn't actually test main.go, but    //
//     rather contains some small helper functions    //
//     that might be used by all the tests            //
//                                                    //
//                       .---.                        //
//                      /_____\                       //
//                     _HH.H.HH                       //
//      _          _-"" WHHHHHW""--__                 //
//      \\      _-"   __\VW=WV/__   /"".              //
//       \\  _-" \__--"  "-_-"   """    "_            //
//        \\/      _                      ""          //
//         \\----_/_|     ___      /"\  T""\====-     //
//          \\ /"-._     |%|H|    (   "\|) | /  .:)   //
//           \/     /    |-+-|     \    |_ J .:::-'   //
//           /     /     |H|%|  _-' '-._  " )/;"      //
//          /     / \    __    (  \ \   \   "         //
//         /     /\/ '. /  \   \ \ \ _- \             //
//         "'-._/  \/  \    "-_ \ -"" _- \            //
//        _,'\\  \  \/  )      "-, -""    \           //
//     _,'_- _ \\ \  \,'          \ \_\_\  \          //
//   ,'    _-    \_\  \            \ \_\_\  \         //
//   \_ _-   _- _,' \  \            \ """"   )        //
//    C\_ _- _,'     \  "--------.   L_""""_/         //
//     " \/-'         "-_________|     '"-Y           //
////////////////////////////////////////////////////////

var arbitraryError error
var exampleTime time.Time

func init() {
	log.SetOutput(ioutil.Discard)
	arbitraryError = fmt.Errorf("arbitrary error")

	var timeParseErr error
	exampleTime, timeParseErr = time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
	if timeParseErr != nil {
		log.Fatalf("error parsing time")
	}
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
