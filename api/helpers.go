package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg/orm"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types ¯\_(ツ)_/¯
	DefaultLimitString = "25"
)

////////////////////////////////////////////////////////////////////////////////////////////////
//       ¸,ø¤º°º¤ø,¸¸,ø¤º°       Begin ~stolen~ borrowed structs.      °º¤ø,¸¸,ø¤º°º¤ø,¸      //
////////////////////////////////////////////////////////////////////////////////////////////////

// borrowed from https://www.reddit.com/r/golang/comments/3ibxdt/null_time_value_in_sql_results/

// NullTime is a time field which can be null
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

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

// genericActiveFilter limits the rows to only active products. This function should be used
// for single-table requests only.
func genericActiveFilter(req *http.Request, q *orm.Query) {
	if req.URL.Query().Get("include_archived") != "true" {
		q.Where("archived_at is null")
	}
}

// genericListActiveFilter limits the rows to only active products from multiple tables.
// This function should be used for multi-table requests only.
func genericListActiveFilter(req *http.Request, tables []string, q *orm.Query) {
	if req.URL.Query().Get("include_archived") != "true" {
		for _, tableName := range tables {
			q.Where(fmt.Sprintf("%s.archived_at is null", tableName))
		}
	}
}

// GenericListQueryHandler aims to generalize retrieving a list of paginated
// models from the database as much as possible without being dangerous.
func genericListQueryHandler(req *http.Request, input *orm.Query) (*orm.Pager, error) {
	ensureValidLimitInRequest(req, input)

	pager := orm.NewPager(req.URL.Query(), DefaultLimit)
	pager.Paginate(input)

	err := input.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		return pager, err
	}

	return pager, nil
}

// EnsureValidLimitInRequest determines the requested limits, and supplies a default if the request doesn't contain one
func ensureValidLimitInRequest(req *http.Request, q *orm.Query) {
	requestedLimit := req.URL.Query().Get("limit")

	if requestedLimit == "" {
		req.URL.Query().Set("limit", DefaultLimitString)
	}
}
