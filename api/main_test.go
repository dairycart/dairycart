package api

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

////////////////////////////////////////////////////////
//                                                    //
//   This file doesn't actually test main.go, but     //
//      rather contains some small helper functions   //
//      that might be used by all the tests           //
//                                                    //
//                       .---.                        //
//                      /_____\                       //
//                     _HH.H.HH                       //
//      _          _-"" WHHHHHW""--__                 //
//      \\      _-"   __\VW=WV/__   /"".              //
//       \\  _-" \__--"  "-_-"   """    "_            //
//        \\/ PhH  _                      ""          //
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

func setupMockRequestsAndMux(db *sql.DB) (*httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	SetupAPIRoutes(m, db)
	return httptest.NewRecorder(), m
}

// this function is lame
func formatConstantQueryForSQLMock(query string) string {
	for _, x := range []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"} {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}
	return query
}
