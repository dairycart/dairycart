package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	discountsTableColumns = `
		id,
		name,
		type,
		amount,
		starts_on,
		expires_on,
		requires_code,
		code,
		limited_use,
		number_of_uses,
		login_required,
		created_on,
		updated_on,
		archived_on
	`

	discountRetrievalQuery = `SELECT * FROM discounts WHERE id = $1`
	discountExistenceQuery = `SELECT EXISTS(SELECT 1 FROM discounts WHERE id = $1 AND archived_on IS NULL)`
	discountDeletionQuery  = `UPDATE discounts SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`
)

// Discount represents pricing changes that apply temporarily to products
type Discount struct {
	DBRow
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     NullTime  `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code,omitempty"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses,omitempty"`
	LoginRequired bool      `json:"login_required"`
}

// DiscountCreationInput represents user input for creating new discounts
type DiscountCreationInput struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     NullTime  `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses"`
	LoginRequired bool      `json:"login_required"`
}

// DiscountsResponse is a discount response struct
type DiscountsResponse struct {
	ListResponse
	Data []Discount `json:"data"`
}

func (d *Discount) discountTypeIsValid() bool {
	// Because Go doesn't have typed enums (https://github.com/golang/go/issues/19814),
	// this is my only real line of defense against a user attempting to load an invalid
	// discount type into the database. It's lame, type enums aren't, here's hoping.
	return d.Type == "percentage" || d.Type == "flat_amount"
}

func retrieveDiscountFromDB(db *sqlx.DB, discountID string) (Discount, error) {
	var d Discount
	err := db.Get(&d, discountRetrievalQuery, discountID)
	return d, err
}

func buildDiscountRetrievalHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountRetrievalHandler is a request handler that returns a single Discount
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := chi.URLParam(req, "discount_id")

		discount, err := retrieveDiscountFromDB(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		json.NewEncoder(res).Encode(discount)
	}
}

func buildDiscountListRetrievalHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountListRetrievalHandler is a request handler that returns a list of Discounts
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := getRowCount(db, "discounts", queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of discounts from the database")
			return
		}

		var discounts []Discount
		query, args := buildDiscountListQuery(queryFilter)
		err = retrieveListOfRowsFromDB(db, query, args, &discounts)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discounts from the database")
			return
		}

		discountsResponse := &DiscountsResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: count,
			},
			Data: discounts,
		}
		json.NewEncoder(res).Encode(discountsResponse)
	}
}

func createDiscountInDB(db *sqlx.DB, in *Discount) (uint64, time.Time, error) {
	var createdID uint64
	var createdOn time.Time
	discountCreationQuery, queryArgs := buildDiscountCreationQuery(in)
	err := db.QueryRow(discountCreationQuery, queryArgs...).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

func buildDiscountCreationHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountCreationHandler is a request handler that creates a Discount from user input
	return func(res http.ResponseWriter, req *http.Request) {
		newDiscount := &Discount{}
		err := validateRequestInput(req, newDiscount)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		id, createdOn, err := createDiscountInDB(db, newDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert discount into database")
			return
		}
		newDiscount.ID = id
		newDiscount.CreatedOn = createdOn

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newDiscount)
	}
}

func archiveDiscount(db *sqlx.DB, discountID string) error {
	_, err := db.Exec(discountDeletionQuery, discountID)
	return err
}

func buildDiscountDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := chi.URLParam(req, "discount_id")

		// can't delete a discount that doesn't exist!
		exists, err := rowExistsInDB(db, discountExistenceQuery, discountID)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		}

		err = archiveDiscount(db, discountID)
		if err != nil {
			notifyOfInternalIssue(res, err, "archive discount in database")
			return
		}

		io.WriteString(res, fmt.Sprintf("Successfully archived discount `%s`", discountID))
	}
}

func updateDiscountInDatabase(db *sqlx.DB, up *Discount) (time.Time, error) {
	var updatedTime time.Time
	discountUpdateQuery, queryArgs := buildDiscountUpdateQuery(up)
	err := db.QueryRow(discountUpdateQuery, queryArgs...).Scan(&updatedTime)
	return updatedTime, err
}

func buildDiscountUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountUpdateHandler is a request handler that can update discounts
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := chi.URLParam(req, "discount_id")

		updatedDiscount := &Discount{}
		err := validateRequestInput(req, updatedDiscount)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingDiscount, err := retrieveDiscountFromDB(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(updatedDiscount, &existingDiscount)

		updatedOn, err := updateDiscountInDatabase(db, updatedDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}
		updatedDiscount.UpdatedOn = NullTime{pq.NullTime{Time: updatedOn, Valid: true}}

		json.NewEncoder(res).Encode(updatedDiscount)
	}
}
