package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	productOptionValueRetrievalForOptionIDQuery = `SELECT id, product_option_id, value, created_on, updated_on, archived_on FROM product_option_values WHERE product_option_id = $1 AND archived_on IS NULL`
)

// retrieveProductOptionValuesForOptionFromDB retrieves a list of ProductOptionValue with a given product_option_id from the database
func retrieveProductOptionValuesForOptionFromDB(db *sqlx.DB, optionID uint64) ([]models.ProductOptionValue, error) {
	var values []models.ProductOptionValue
	err := db.Select(&values, productOptionValueRetrievalForOptionIDQuery, optionID)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	return values, nil
}

func buildProductOptionValueUpdateHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductOptionValueUpdateHandler is a request handler that can update product option values
	return func(res http.ResponseWriter, req *http.Request) {
		optionValueIDStr := chi.URLParam(req, "option_value_id")
		// we can eat this error because Mux takes care of validating route params for us
		optionValueID, _ := strconv.ParseUint(optionValueIDStr, 10, 64)

		updatedValueData := &models.ProductOptionValue{}
		err := validateRequestInput(req, updatedValueData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't update an option value that doesn't exist!
		existingOptionValue, err := client.GetProductOptionValue(db, optionValueID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product option value", optionValueIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}
		existingOptionValue.Value = updatedValueData.Value

		updatedOn, err := client.UpdateProductOptionValue(db, existingOptionValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option value in the database")
			return
		}
		existingOptionValue.UpdatedOn = models.NullTime{NullTime: pq.NullTime{Time: updatedOn, Valid: true}}

		json.NewEncoder(res).Encode(existingOptionValue)
	}
}

func buildProductOptionValueCreationHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// productOptionValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		optionIDStr := chi.URLParam(req, "option_id")
		// we can eat this error because Mux takes care of validating route params for us
		optionID, _ := strconv.ParseUint(optionIDStr, 10, 64)

		newValue := &models.ProductOptionValue{}
		err := validateRequestInput(req, newValue)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create values for a product option that doesn't exist
		productOptionExists, err := client.ProductOptionExists(db, optionID)
		if err == sql.ErrNoRows || !productOptionExists {
			respondThatRowDoesNotExist(req, res, "product option", optionIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}
		newValue.ProductOptionID = optionID

		// can't create a product option value that already exists
		valueExists, err := client.ProductOptionValueForOptionIDExists(db, optionID, newValue.Value)
		if valueExists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option value '%s' already exists for option ID %d", newValue.Value, optionID))
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a transaction")
			return
		}

		newValue.ID, newValue.CreatedOn, err = client.CreateProductOptionValue(tx, newValue)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newValue)
	}
}

func buildProductOptionValueDeletionHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductOptionValueDeletionHandler is a request handler that can delete product option values
	return func(res http.ResponseWriter, req *http.Request) {
		optionValueIDStr := chi.URLParam(req, "option_value_id")
		// we can eat this error because Mux takes care of validating route params for us
		optionValueID, _ := strconv.ParseUint(optionValueIDStr, 10, 64)

		// can't delete an option value that doesn't exist!
		optionValue, err := client.GetProductOptionValue(db, optionValueID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product option value", optionValueIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discount from database")
			return
		}

		archivedOn, err := client.DeleteProductOptionValue(db, optionValueID)
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}
		optionValue.ArchivedOn = models.NullTime{NullTime: pq.NullTime{Time: archivedOn, Valid: true}}

		json.NewEncoder(res).Encode(optionValue)
	}
}
