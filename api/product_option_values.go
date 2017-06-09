package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	productOptionValueExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_option_values WHERE id = $1 AND archived_at IS NULL)`
	productOptionValueRetrievalQuery = `SELECT * FROM product_option_values WHERE id = $1 AND archived_at IS NULL`
)

// ProductOptionValue represents a product's option values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size, and six ProductOptionValues,
// One for each color and one for each size.
type ProductOptionValue struct {
	ID              int64     `json:"id"`
	ProductOptionID int64     `json:"product_option_id"`
	Value           string    `json:"value"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       NullTime  `json:"updated_at,omitempty"`
	ArchivedAt      NullTime  `json:"archived_at,omitempty"`
}

func (pav *ProductOptionValue) generateScanArgs() []interface{} {
	return []interface{}{
		&pav.ID,
		&pav.ProductOptionID,
		&pav.Value,
		&pav.CreatedAt,
		&pav.UpdatedAt,
		&pav.ArchivedAt,
	}
}

// ProductOptionValueCreationInput is a struct to use for creating product option values
type ProductOptionValueCreationInput struct {
	ProductOptionID int64
	Value           string `json:"value"`
}

// ProductOptionValueUpdateInput is a struct to use for updating product option values
type ProductOptionValueUpdateInput struct {
	Value string `json:"value"`
}

func validateProductOptionValueUpdateInput(req *http.Request) (*ProductOptionValue, error) {
	i := &ProductOptionValueUpdateInput{}
	json.NewDecoder(req.Body).Decode(i)

	s := structs.New(i)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product option value body")
	}

	out := &ProductOptionValue{
		Value: i.Value,
	}

	return out, nil
}

// retrieveProductOptionValue retrieves a ProductOptionValue with a given ID from the database
func retrieveProductOptionValueFromDB(db *sql.DB, id int64) (*ProductOptionValue, error) {
	v := &ProductOptionValue{}
	err := db.QueryRow(productOptionValueRetrievalQuery, id).Scan(v.generateScanArgs()...)
	if err == sql.ErrNoRows {
		return v, errors.Wrap(err, "Error querying for product option values")
	}
	return v, err
}

// retrieveProductOptionValue retrieves a ProductOptionValue with a given product option ID from the database
func retrieveProductOptionValueForOptionFromDB(db *sql.DB, optionID int64) ([]*ProductOptionValue, error) {
	var values []*ProductOptionValue

	query := buildProductOptionValueRetrievalForOptionIDQuery(optionID)
	rows, err := db.Query(query, optionID)

	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	defer rows.Close()
	for rows.Next() {
		value := &ProductOptionValue{}
		_ = rows.Scan(value.generateScanArgs()...)
		values = append(values, value)
	}
	return values, nil
}

func validateProductOptionValueCreationInput(req *http.Request) (*ProductOptionValue, error) {
	i := &ProductOptionValueCreationInput{}
	err := json.NewDecoder(req.Body).Decode(i)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	s := structs.New(i)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product option value body")
	}

	v := &ProductOptionValue{
		Value: i.Value,
	}

	return v, err
}

func updateProductOptionValueInDB(db *sql.DB, v *ProductOptionValue) error {
	valueUpdateQuery, queryArgs := buildProductOptionValueUpdateQuery(v)
	err := db.QueryRow(valueUpdateQuery, queryArgs...).Scan(v.generateScanArgs()...)
	return err
}

func buildProductOptionValueUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductOptionValueUpdateHandler is a request handler that can update product option values
		reqVars := mux.Vars(req)
		optionValueID := reqVars["option_value_id"]
		// eating these errors because Mux should validate these for us.
		optionValueIDInt, _ := strconv.Atoi(optionValueID)

		// can't update an option value that doesn't exist!
		optionValueExists, err := rowExistsInDB(db, "product_option_values", "id", optionValueID)
		if err != nil || !optionValueExists {
			respondThatRowDoesNotExist(req, res, "product option value", optionValueID)
			return
		}

		updatedValueData, err := validateProductOptionValueUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingOptionValue, err := retrieveProductOptionValueFromDB(db, int64(optionValueIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option value from the database")
			return
		}
		existingOptionValue.Value = updatedValueData.Value

		err = updateProductOptionValueInDB(db, existingOptionValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option value in the database")
			return
		}

		json.NewEncoder(res).Encode(existingOptionValue)
	}
}

// createProductOptionValueInDB creates a ProductOptionValue tied to a ProductOption
func createProductOptionValueInDB(tx *sql.Tx, v *ProductOptionValue) (int64, error) {
	var newOptionValueID int64
	query, args := buildProductOptionValueCreationQuery(v)
	err := tx.QueryRow(query, args...).Scan(&newOptionValueID)
	return newOptionValueID, err
}

func optionValueAlreadyExistsForOption(db *sql.DB, optionID int64, value string) (bool, error) {
	var exists string

	query, args := buildProductOptionValueExistenceForOptionIDQuery(optionID, value)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

func buildProductOptionValueCreationHandler(db *sql.DB) http.HandlerFunc {
	// productOptionValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		optionID := mux.Vars(req)["option_id"]

		// we can eat this error because Mux takes care of validating route params for us
		optionIDInt, _ := strconv.ParseInt(optionID, 10, 64)

		// can't create values for a product option that doesn't exist
		productOptionValueExistsByID, err := rowExistsInDB(db, "product_options", "id", optionID)
		if err != nil || !productOptionValueExistsByID {
			respondThatRowDoesNotExist(req, res, "product option", optionID)
			return
		}

		newProductOptionValue, err := validateProductOptionValueCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		newProductOptionValue.ProductOptionID = optionIDInt

		// can't create a product option value that already exists
		productOptionValueExistsByValue, err := optionValueAlreadyExistsForOption(db, optionIDInt, newProductOptionValue.Value)
		if err != nil || productOptionValueExistsByValue {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option value `%s` already exists for option ID %s", newProductOptionValue.Value, optionID))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a transasction")
			return
		}

		newProductOptionValueID, err := createProductOptionValueInDB(tx, newProductOptionValue)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProductOptionValue.ID = newProductOptionValueID

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newProductOptionValue)
	}
}
