package models

import (
	"time"

	"github.com/lib/pq"
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

////////////////////////////////////////////////////////////////////////////////////////////////
//    ¸,ø¤º°º¤ø,¸¸,ø¤º°   Everything after this point is not borrowed.   °º¤ø,¸¸,ø¤º°º¤ø,¸    //
////////////////////////////////////////////////////////////////////////////////////////////////

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count uint64      `json:"count"`
	Limit uint8       `json:"limit"`
	Page  uint64      `json:"page"`
	Data  interface{} `json:"data"`
}

// ErrorResponse is a handy struct we can respond with in the event we have an error to report
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// QueryFilter represents a query filter
type QueryFilter struct {
	Page            uint64
	Limit           uint8
	CreatedAfter    time.Time
	CreatedBefore   time.Time
	UpdatedAfter    time.Time
	UpdatedBefore   time.Time
	IncludeArchived bool
}
