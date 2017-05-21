package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types are a thing
	DefaultLimitString = "25"
)

////////////////////////////////////////////////////////////////////////////////////////////////
//       ¸,ø¤º°º¤ø,¸¸,ø¤º°       Begin ~stolen~ borrowed structs.      °º¤ø,¸¸,ø¤º°º¤ø,¸      //
////////////////////////////////////////////////////////////////////////////////////////////////

// borrowed from http://stackoverflow.com/questions/32825640/custom-marshaltext-for-golang-sql-null-types

// There's not really a great solution for these two stinkers here. Because []byte is what's expected, passing
// nil results in an empty string. The original has []byte("null"), which I think is actually worse. At least
// an empty string is falsy in most languages. ¯\_(ツ)_/¯

// NullFloat64 is a json.Marshal-able 64-bit float.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (nf NullFloat64) MarshalText() ([]byte, error) {
	if nf.Valid {
		nfv := nf.Float64
		return []byte(strconv.FormatFloat(nfv, 'f', -1, 64)), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullFloat64 so that gorilla/schema can parse it
func (nf *NullFloat64) UnmarshalText(text []byte) (err error) {
	nf.NullFloat64.Float64, err = strconv.ParseFloat(string(text), 64)
	return err
}

// This isn't borrowed, but rather inferred from stuff I borrowed above

// NullString is a json.Marshal-able String.
type NullString struct {
	sql.NullString
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (ns NullString) MarshalText() ([]byte, error) {
	if ns.Valid {
		nsv := ns.String
		return []byte(nsv), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullString so that gorilla/schema can parse it
func (ns *NullString) UnmarshalText(text []byte) (err error) {
	ns.String = string(text)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
//        ¸,ø¤º°º¤ø,¸¸,ø¤º°       End ~stolen~ borrowed structs.       °º¤ø,¸¸,ø¤º°º¤ø,¸      //
////////////////////////////////////////////////////////////////////////////////////////////////

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count int `json:"count"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

// respondToInvalidRequest takes an error, a string format, and a response object, and
// writes an error response when Dairycart determines that a user is at fault or provides
// information that would otherwise cause an error. Things like providing an inadequate sku
// are handled by this function.
func respondToInvalidRequest(err error, errorFormat string, res http.ResponseWriter) {
	errorString := fmt.Sprintf("%s: %v", errorFormat, err)
	log.Println(errorString)
	http.Error(res, errorFormat, http.StatusBadRequest)
}

// informOfServerIssue takes an error, a string format, and a response object, and
// writes an error response when during the course of normal operation, Dairycart
// experiences something out of the user's hands. Things like database query errors
// are handled by this function.
func informOfServerIssue(err error, errorFormat string, res http.ResponseWriter) {
	errorString := fmt.Sprintf("%s: %v", errorFormat, err)
	log.Println(errorString)
	http.Error(res, errorFormat, http.StatusInternalServerError)
}

// ensureRequestBodyValidity takes a request object and checks that:
// 		1) the body contained therein is not nil and
//		2) the body decodes without errors into a particular struct
// It then returns whether or not either of those conditions was untrue.
func ensureRequestBodyValidity(res http.ResponseWriter, req *http.Request, object interface{}) bool {
	if req.Body == nil {
		http.Error(res, "Please send a request body", http.StatusBadRequest)
		return true
	}

	err := json.NewDecoder(req.Body).Decode(object)
	if err != nil {
		respondToInvalidRequest(err, "Invalid request body", res)
	}
	return err != nil
}
