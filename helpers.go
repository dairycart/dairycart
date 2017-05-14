package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg/orm"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types ¯\_(ツ)_/¯
	DefaultLimitString = "25"
)

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
	http.Error(res, errorString, http.StatusBadRequest)
}

// informOfServerIssue takes an error, a string format, and a response object, and
// writes an error response when during the course of normal operation, Dairycart
// experiences something out of the user's hands. Things like database query errors
// are handled by this function.
func informOfServerIssue(err error, errorFormat string, res http.ResponseWriter) {
	errorString := fmt.Sprintf("%s: %v", errorFormat, err)
	log.Println(errorString)
	http.Error(res, errorString, http.StatusInternalServerError)
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

// EnsureValidLimitInRequest determines the requested limits, and supplies a default if the request doesn't contain one
func ensureValidLimitInRequest(req *http.Request, q *orm.Query) {
	requestedLimit := req.URL.Query().Get("limit")

	if requestedLimit == "" {
		req.URL.Query().Set("limit", DefaultLimitString)
	}
}
