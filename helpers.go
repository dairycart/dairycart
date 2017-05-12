package main

import (
	"net/http"
	"strconv"

	"github.com/go-pg/pg/orm"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form
	DefaultLimitString = "25"
)

// LimitRequestedQuery determines the requested limits, and supplies a default if the request doesn't contain one
func LimitRequestedQuery(req *http.Request, q *orm.Query) int {
	var actualLimit int
	requestedLimit := req.URL.Query().Get("limit")
	actualLimit, err := strconv.Atoi(requestedLimit)

	if requestedLimit == "" || err != nil {
		actualLimit = DefaultLimit
		req.URL.Query().Set("limit", DefaultLimitString)
	}

	q.Limit(actualLimit)

	// We return actualLimit so it can be used in the response object.
	return actualLimit
}
