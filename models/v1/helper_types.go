package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
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

// Dairytime is a custom time pointer struct that should interface well with Postgres's time type and allow for easily nullable time
type Dairytime struct {
	time.Time
}

// Scan implements the Scanner interface.
func (dt *Dairytime) Scan(value interface{}) error {
	if t, ok := value.(time.Time); !ok {
		return errors.New("value is not a time.Time")
	} else {
		dt.Time = t
	}
	return nil
}

// Value implements the driver Valuer interface.
func (dt Dairytime) Value() (driver.Value, error) {
	return dt.Time, nil
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (dt Dairytime) MarshalText() ([]byte, error) {
	if dt.Time.IsZero() {
		return nil, nil
	}
	return []byte(dt.Time.Format(timeLayout)), nil
}

// UnmarshalText is a function which unmarshals a NullTime
func (dt *Dairytime) UnmarshalText(text []byte) (err error) {
	if text == nil {
		return nil
	}
	if len(text) == 0 {
		return nil
	}

	dt.Time, _ = time.Parse(timeLayout, string(text))
	return nil
}

var _ fmt.Stringer = (*Dairytime)(nil)

func (dt *Dairytime) String() string {
	if dt == nil {
		return "nil"
	}
	return dt.Time.Format(timeLayout)
}

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count uint64 `json:"count"`
	Limit uint8  `json:"limit"`
	Page  uint64 `json:"page"`
}

var _ = error(new(ErrorResponse))

// ErrorResponse is a handy struct we can respond with in the event we have an error to report
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
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
