package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	productOptionValueExistenceQuery                 = `SELECT EXISTS(SELECT 1 FROM product_option_values WHERE id = $1 AND archived_on IS NULL)`
	productOptionValueExistenceForOptionIDQuery      = `SELECT EXISTS(SELECT 1 FROM product_option_values WHERE product_option_id = $1 AND value = $2 AND archived_on IS NULL)`
	productOptionValueRetrievalQuery                 = `SELECT id, product_option_id, value, created_on, updated_on, archived_on FROM product_option_values WHERE id = $1`
	productOptionValueRetrievalForOptionIDQuery      = `SELECT id, product_option_id, value, created_on, updated_on, archived_on FROM product_option_values WHERE product_option_id = $1 AND archived_on IS NULL`
	productOptionValueDeletionQuery                  = `UPDATE product_option_values SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`
	productVariantBridgeDeletionQueryByProductID     = `UPDATE product_variant_bridge SET archived_on = NOW() WHERE product_id = $1 AND archived_on IS NULL`
)

// ProductOptionValue represents a product's option values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size, and six ProductOptionValues,
// One for each color and one for each size.
type ProductOptionValue struct {
	DBRow
	ProductOptionID uint64 `json:"product_option_id"`
	Value           string `json:"value"`
}

func createBridgeEntryForProductValues(tx *sql.Tx, productID uint64, ids []uint64) error {
	query, queryArgs := buildProductVariantBridgeCreationQuery(productID, ids)
	_, err := tx.Exec(query, queryArgs...)
	return err
}

// retrieveProductOptionValue retrieves a ProductOptionValue with a given ID from the database
func retrieveProductOptionValueFromDB(db *sqlx.DB, id uint64) (*ProductOptionValue, error) {
	v := &ProductOptionValue{}
	err := db.QueryRowx(productOptionValueRetrievalQuery, id).StructScan(v)
	if err == sql.ErrNoRows {
		return v, errors.Wrap(err, "Error querying for product option values")
	}
	return v, err
}

// retrieveProductOptionValuesForOptionFromDB retrieves a list of ProductOptionValue with a given product_option_id from the database
func retrieveProductOptionValuesForOptionFromDB(db *sqlx.DB, optionID uint64) ([]ProductOptionValue, error) {
	var values []ProductOptionValue
	err := db.Select(&values, productOptionValueRetrievalForOptionIDQuery, optionID)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	return values, nil
}

func updateProductOptionValueInDB(db *sqlx.DB, v *ProductOptionValue) (time.Time, error) {
	var updatedOn time.Time
	valueUpdateQuery, queryArgs := buildProductOptionValueUpdateQuery(v)
	err := db.QueryRow(valueUpdateQuery, queryArgs...).Scan(&updatedOn)
	return updatedOn, err
}

func buildProductOptionValueUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductOptionValueUpdateHandler is a request handler that can update product option values
		optionValueID := chi.URLParam(req, "option_value_id")
		// eating these errors because Mux should validate these for us.
		optionValueIDInt, _ := strconv.Atoi(optionValueID)

		// can't update an option value that doesn't exist!
		optionValueExists, err := rowExistsInDB(db, productOptionValueExistenceQuery, optionValueID)
		if err != nil || !optionValueExists {
			respondThatRowDoesNotExist(req, res, "product option value", optionValueID)
			return
		}

		updatedValueData := &ProductOptionValue{}
		err = validateRequestInput(req, updatedValueData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingOptionValue, err := retrieveProductOptionValueFromDB(db, uint64(optionValueIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option value from the database")
			return
		}
		existingOptionValue.Value = updatedValueData.Value

		updatedOn, err := updateProductOptionValueInDB(db, existingOptionValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option value in the database")
			return
		}
		existingOptionValue.UpdatedOn = NullTime{pq.NullTime{Time: updatedOn, Valid: true}}

		json.NewEncoder(res).Encode(existingOptionValue)
	}
}

// createProductOptionValueInDB creates a ProductOptionValue tied to a ProductOption
func createProductOptionValueInDB(tx *sql.Tx, v *ProductOptionValue) (uint64, time.Time, error) {
	var newOptionValueID uint64
	var createdOn time.Time
	query, args := buildProductOptionValueCreationQuery(v)
	err := tx.QueryRow(query, args...).Scan(&newOptionValueID, &createdOn)
	return newOptionValueID, createdOn, err
}

func optionValueAlreadyExistsForOption(db *sqlx.DB, optionID int64, value string) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionValueExistenceForOptionIDQuery, optionID, value).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

func buildProductOptionValueCreationHandler(db *sqlx.DB) http.HandlerFunc {
	// productOptionValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		optionID := chi.URLParam(req, "option_id")

		// we can eat this error because Mux takes care of validating route params for us
		optionIDInt, _ := strconv.ParseInt(optionID, 10, 64)

		// can't create values for a product option that doesn't exist
		productOptionExistsByID, err := rowExistsInDB(db, productOptionExistenceQuery, optionID)
		if err != nil || !productOptionExistsByID {
			respondThatRowDoesNotExist(req, res, "product option", optionID)
			return
		}

		newProductOptionValue := &ProductOptionValue{}
		err = validateRequestInput(req, newProductOptionValue)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		newProductOptionValue.ProductOptionID = uint64(optionIDInt)

		// can't create a product option value that already exists
		productOptionValueExistsByValue, err := optionValueAlreadyExistsForOption(db, optionIDInt, newProductOptionValue.Value)
		if err != nil || productOptionValueExistsByValue {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option value '%s' already exists for option ID %s", newProductOptionValue.Value, optionID))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a transaction")
			return
		}

		newProductOptionValue.ID, newProductOptionValue.CreatedOn, err = createProductOptionValueInDB(tx, newProductOptionValue)
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
		json.NewEncoder(res).Encode(newProductOptionValue)
	}
}

func archiveProductOptionValue(db *sqlx.DB, id uint64) error {
	_, err := db.Exec(productOptionValueDeletionQuery, id)
	return err
}

func buildProductOptionValueDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductOptionValueDeletionHandler is a request handler that can delete product option values
		optionValueID := chi.URLParam(req, "option_value_id")
		// eating these errors because Mux should validate these for us.
		optionValueIDInt, _ := strconv.Atoi(optionValueID)

		// can't delete an option value that doesn't exist!
		optionValueExists, err := rowExistsInDB(db, productOptionValueExistenceQuery, optionValueID)
		if err != nil || !optionValueExists {
			respondThatRowDoesNotExist(req, res, "product option value", optionValueID)
			return
		}

		err = archiveProductOptionValue(db, uint64(optionValueIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
