package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

// UnmarshalText is a function which unmarshals a NullFloat64
func (nf *NullFloat64) UnmarshalText(text []byte) (err error) {
	s := string(text)
	nf.NullFloat64.Float64, err = strconv.ParseFloat(s, 64)
	if err != nil {
		nf.NullFloat64.Float64 = 0
	}
	// returning nil because we've ensured that Float64 is set to at least zero.
	return nil
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

// rowExistsInDB will return whether or not a product/attribute/etc with a given identifier exists in the database
func rowExistsInDB(db *sql.DB, table, identifier, id string) (bool, error) {
	var exists string

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1 and archived_at is null);`, table, identifier)
	err := db.QueryRow(query, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, errors.Wrap(err, "Error querying for row")
	}

	return exists == "true", err
}

func respondThatRowDoesNotExist(req *http.Request, res http.ResponseWriter, itemType, identifierType, identifier string) {
	log.Printf("informing user that the %s they were looking for (%s %s) does not exist", itemType, identifierType, identifier)
	http.Error(res, fmt.Sprintf("No %s with the %s '%s' found", itemType, identifierType, identifier), http.StatusNotFound)
}

func notifyOfInvalidRequestBody(res http.ResponseWriter, err error) {
	log.Printf("Encountered this error decoding a request body: %v", err)
	http.Error(res, "Invalid request body", http.StatusBadRequest)
}

func notifyOfInternalIssue(res http.ResponseWriter, err error, attemptedTask string) {
	log.Println(fmt.Sprintf("Encountered this error trying to %s: %v", attemptedTask, err))
	http.Error(res, "Unexpected internal error", http.StatusInternalServerError)
}
