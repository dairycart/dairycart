package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	discountsTableColumns = `
		id,
		name,
		discount_type,
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

func retrieveDiscountFromDB(db *sqlx.DB, discountID string) (models.Discount, error) {
	var d models.Discount
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

		var discounts []models.Discount
		query, args := buildDiscountListQuery(queryFilter)
		err = retrieveListOfRowsFromDB(db, query, args, &discounts)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discounts from the database")
			return
		}

		discountsResponse := &ListResponse{
			Page:  queryFilter.Page,
			Limit: queryFilter.Limit,
			Count: count,
			Data:  discounts,
		}
		json.NewEncoder(res).Encode(discountsResponse)
	}
}

func createDiscountInDB(db *sqlx.DB, in *models.Discount) (uint64, time.Time, error) {
	var createdID uint64
	var createdOn time.Time
	discountCreationQuery, queryArgs := buildDiscountCreationQuery(in)
	err := db.QueryRow(discountCreationQuery, queryArgs...).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

func buildDiscountCreationHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountCreationHandler is a request handler that creates a Discount from user input
	return func(res http.ResponseWriter, req *http.Request) {
		newDiscount := &models.Discount{}
		err := validateRequestInput(req, newDiscount)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		newDiscount.ID, newDiscount.CreatedOn, err = createDiscountInDB(db, newDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert discount into database")
			return
		}

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

func updateDiscountInDatabase(db *sqlx.DB, up *models.Discount) (time.Time, error) {
	var updatedTime time.Time
	discountUpdateQuery, queryArgs := buildDiscountUpdateQuery(up)
	err := db.QueryRow(discountUpdateQuery, queryArgs...).Scan(&updatedTime)
	return updatedTime, err
}

func buildDiscountUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountUpdateHandler is a request handler that can update discounts
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := chi.URLParam(req, "discount_id")

		updatedDiscount := &models.Discount{}
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

		use := func(...interface{}) {}

		updatedOn, err := updateDiscountInDatabase(db, updatedDiscount)
		if err != nil {
			errStr := err.Error()
			use(errStr)
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}
		updatedDiscount.UpdatedOn = models.NullTime{NullTime: pq.NullTime{Time: updatedOn, Valid: true}}

		json.NewEncoder(res).Encode(updatedDiscount)
	}
}
