package main

import (
	"log"
	"net/http"

	"github.com/go-pg/pg/orm"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types
	DefaultLimitString = "25"
)

// ListResponse is a generic list response struct
type ListResponse struct {
	Count int           `json:"count"`
	Limit int           `json:"limit"`
	Page  int           `json:"page"`
	Data  []interface{} `json:"data"`
}

// EnsureValidLimitInRequest determines the requested limits, and supplies a default if the request doesn't contain one
func ensureValidLimitInRequest(req *http.Request, q *orm.Query) {
	requestedLimit := req.URL.Query().Get("limit")

	if requestedLimit == "" {
		req.URL.Query().Set("limit", DefaultLimitString)
	}
}

// GenericActiveFilter limits the rows to only active products
func genericActiveFilter(req *http.Request, q *orm.Query) {
	if req.URL.Query().Get("include_archived") != "true" {
		q.Where("archived_at is null")
	}
}

// GenericListQueryHandler aims to generalize retrieving a list of paginated
// models from the database as much as possible without being dangerous.
func genericListQueryHandler(req *http.Request, input *orm.Query, activeRowFilterFunc func(req *http.Request, q *orm.Query)) (*orm.Pager, error) {
	ensureValidLimitInRequest(req, input)
	activeRowFilterFunc(req, input)

	pager := orm.NewPager(req.URL.Query(), DefaultLimit)
	pager.Paginate(input)

	err := input.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
		return pager, err
	}

	return pager, nil
}
