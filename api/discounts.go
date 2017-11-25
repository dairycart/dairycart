package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
)

func buildDiscountRetrievalHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// DiscountRetrievalHandler is a request handler that returns a single Discount
	return func(res http.ResponseWriter, req *http.Request) {
		discountIDStr := chi.URLParam(req, "discount_id")
		// eating this error because the router should have ensured this is an integer
		discountID, _ := strconv.ParseUint(discountIDStr, 10, 64)

		discount, err := client.GetDiscount(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		json.NewEncoder(res).Encode(discount)
	}
}

func buildDiscountListRetrievalHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// DiscountListRetrievalHandler is a request handler that returns a list of Discounts
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)

		count, err := client.GetDiscountCount(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of discounts from the database")
			return
		}

		discounts, err := client.GetDiscountList(db, queryFilter)
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

func buildDiscountCreationHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// DiscountCreationHandler is a request handler that creates a Discount from user input
	return func(res http.ResponseWriter, req *http.Request) {
		newDiscount := &models.Discount{}
		err := validateRequestInput(req, newDiscount)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		newDiscount.ID, newDiscount.CreatedOn, err = client.CreateDiscount(db, newDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert discount into database")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newDiscount)
	}
}

func buildDiscountDeletionHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		discountIDStr := chi.URLParam(req, "discount_id")
		// eating this error because the router should have ensured this is an integer
		discountID, _ := strconv.ParseUint(discountIDStr, 10, 64)

		discount, err := client.GetDiscount(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		archivedOn, err := client.DeleteDiscount(db, discountID)
		if err != nil {
			notifyOfInternalIssue(res, err, "archive discount in database")
			return
		}
		discount.ArchivedOn = models.NullTime{NullTime: pq.NullTime{Time: archivedOn, Valid: true}}

		json.NewEncoder(res).Encode(discount)
	}
}

func buildDiscountUpdateHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// DiscountUpdateHandler is a request handler that can update discounts
	return func(res http.ResponseWriter, req *http.Request) {
		discountIDStr := chi.URLParam(req, "discount_id")
		// eating this error because the router should have ensured this is an integer
		discountID, _ := strconv.ParseUint(discountIDStr, 10, 64)

		updatedDiscount := &models.Discount{}
		err := validateRequestInput(req, updatedDiscount)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingDiscount, err := client.GetDiscount(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(updatedDiscount, &existingDiscount)

		updatedOn, err := client.UpdateDiscount(db, updatedDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}
		updatedDiscount.UpdatedOn = models.NullTime{NullTime: pq.NullTime{Time: updatedOn, Valid: true}}

		json.NewEncoder(res).Encode(updatedDiscount)
	}
}
