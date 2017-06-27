package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types are a thing
	DefaultLimitString = "25"
	// MaxLimit is the maximum number of objects Dairycart will return in a response
	MaxLimit = 50

	dataValidationPattern = `^[a-zA-Z\-_]{1,50}$`
	timeLayout            = "2006-01-02T15:04:05.000000Z"
)

// Modified from code borrowed from http://stackoverflow.com/questions/32825640/custom-marshaltext-for-golang-sql-null-types
// NullTime is a json.Marshal-able pq.NullTime.
type NullTime struct {
	pq.NullTime
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (nt NullTime) MarshalText() ([]byte, error) {
	if nt.Valid {
		return []byte(nt.Time.Format(timeLayout)), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullTime
func (nt *NullTime) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		nt.Time = time.Time{}
		nt.Valid = true
		return nil
	}

	t, _ := time.Parse(timeLayout, string(text))
	nt.Time = t
	nt.Valid = true
	return nil
}

// NullString is a json.Marshal-able sql.NullString.
type NullString struct {
	sql.NullString
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (ns NullString) MarshalText() ([]byte, error) {
	if ns.Valid {
		return []byte(ns.String), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullString
func (ns *NullString) UnmarshalText(text []byte) (err error) {
	ns.String = string(text)
	ns.Valid = true
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
//    ¸,ø¤º°º¤ø,¸¸,ø¤º°   Everything after this point is not borrowed.   °º¤ø,¸¸,ø¤º°º¤ø,¸    //
////////////////////////////////////////////////////////////////////////////////////////////////

// DBRow is meant to represent the base columns that every database table should have
type DBRow struct {
	ID         uint64    `json:"id"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  NullTime  `json:"updated_on,omitempty"`
	ArchivedOn NullTime  `json:"archived_on,omitempty"`
}

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count uint64 `json:"count"`
	Limit uint8  `json:"limit"`
	Page  uint64 `json:"page"`
}

// ErrorResponse is a handy struct we can respond with in the event we have an error to report
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// QueryFilter represents a query filter
type QueryFilter struct {
	Page          uint64
	Limit         uint8
	CreatedAfter  time.Time
	CreatedBefore time.Time
	UpdatedAfter  time.Time
	UpdatedBefore time.Time
}

func parseRawFilterParams(rawFilterParams url.Values) *QueryFilter {
	qf := &QueryFilter{
		Page:  1,
		Limit: 25,
	}

	page := rawFilterParams["page"]
	if len(page) == 1 {
		i, err := strconv.ParseUint(page[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `Page`, err)
		} else {
			qf.Page = uint64(math.Max(float64(i), 1))
		}
	}

	limit := rawFilterParams["limit"]
	if len(limit) == 1 {
		i, err := strconv.ParseFloat(limit[0], 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `Limit`, err)
		} else {
			qf.Limit = uint8(math.Min(i, MaxLimit))
		}
	}

	updatedAfter := rawFilterParams["updated_after"]
	if len(updatedAfter) == 1 {
		i, err := strconv.ParseUint(updatedAfter[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `UpdatedAfter`, err)
		} else {
			qf.UpdatedAfter = time.Unix(int64(i), 0)
		}
	}

	updatedBefore := rawFilterParams["updated_before"]
	if len(updatedBefore) == 1 {
		i, err := strconv.ParseUint(updatedBefore[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `UpdatedBefore`, err)
		} else {
			qf.UpdatedBefore = time.Unix(int64(i), 0)
		}
	}

	createdAfter := rawFilterParams["created_after"]
	if len(createdAfter) == 1 {
		i, err := strconv.ParseUint(createdAfter[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `CreatedAfter`, err)
		} else {
			qf.CreatedAfter = time.Unix(int64(i), 0)
		}
	}

	createdBefore := rawFilterParams["created_before"]
	if len(createdBefore) == 1 {
		i, err := strconv.ParseUint(createdBefore[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `CreatedBefore`, err)
		} else {
			qf.CreatedBefore = time.Unix(int64(i), 0)
		}
	}

	return qf
}

func dataValueIsValid(input string) bool {
	// This is a rather simple function, but is sort of strictly meant to
	// ensure that certain values (like skus, option values, option names)
	// aren't allowed to have crazy values in the database
	dataValidator := regexp.MustCompile(dataValidationPattern)
	return dataValidator.MatchString(input)
}

func getRowCount(db *sqlx.DB, table string, queryFilter *QueryFilter) (uint64, error) {
	var count uint64
	query := buildCountQuery(table, queryFilter)
	err := db.Get(&count, query)
	return count, err
}

func retrieveListOfRowsFromDB(db *sqlx.DB, query string, args []interface{}, rows interface{}) error {
	return db.Select(rows, query, args...)
}

// rowExistsInDB will return whether or not a product/option/etc with a given identifier exists in the database
func rowExistsInDB(db *sqlx.DB, query string, identifier string) (bool, error) {
	var exists string

	err := db.QueryRow(query, identifier).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		errStr := err.Error()
		noop(errStr)
		return false, err
	}

	return exists == "true", err
}

func respondThatRowDoesNotExist(req *http.Request, res http.ResponseWriter, itemType, id string) {
	itemTypeToIdentifierMap := map[string]string{
		"product option":       "id",
		"product option value": "id",
		"product progenitor":   "id",
		"product":              "sku",
		"discount":             "id",
	}

	// in case we forget one, default to ID
	identifier := itemTypeToIdentifierMap[itemType]
	if _, ok := itemTypeToIdentifierMap[itemType]; !ok {
		identifier = "identified by"
	}

	log.Printf("informing user that the %s they were looking for (%s %s) does not exist", itemType, identifier, id)
	res.WriteHeader(http.StatusNotFound)
	errRes := &ErrorResponse{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("The %s you were looking for (%s `%s`) does not exist", itemType, identifier, id),
	}
	json.NewEncoder(res).Encode(errRes)
}

func notifyOfInvalidRequestBody(res http.ResponseWriter, err error) {
	log.Printf("Encountered this error decoding a request body: %v", err)
	res.WriteHeader(http.StatusBadRequest)
	errRes := &ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: err.Error(),
	}
	json.NewEncoder(res).Encode(errRes)
}

func notifyOfInternalIssue(res http.ResponseWriter, err error, attemptedTask string) {
	log.Println(fmt.Sprintf("Encountered this error trying to %s: %v", attemptedTask, err))
	res.WriteHeader(http.StatusInternalServerError)
	errRes := &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "Unexpected internal error occurred",
	}
	json.NewEncoder(res).Encode(errRes)
}
