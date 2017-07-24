package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	productOptionsHeaders = `id,
		name,
		product_id,
		created_on,
		updated_on,
		archived_on
	`
	productOptionExistenceQuery                 = `SELECT EXISTS(SELECT 1 FROM product_options WHERE id = $1 AND archived_on IS NULL)`
	productOptionRetrievalQuery                 = `SELECT * FROM product_options WHERE id = $1`
	productOptionExistenceQueryForProductByName = `SELECT EXISTS(SELECT 1 FROM product_options WHERE name = $1 AND product_id = $2 and archived_on IS NULL)`
	productOptionDeletionQuery                  = `UPDATE product_options SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`
	productOptionValuesDeletionQueryByOptionID  = `UPDATE product_option_values SET archived_on = NOW() WHERE product_option_id = $1 AND archived_on IS NULL`
)

// ProductOption represents a products variant options. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size.
type ProductOption struct {
	DBRow
	Name      string               `json:"name"`
	ProductID uint64               `json:"product_id"`
	Values    []ProductOptionValue `json:"values"`
}

func (a *ProductOption) generateScanArgs() []interface{} {
	return []interface{}{
		&a.ID,
		&a.Name,
		&a.ProductID,
		&a.CreatedOn,
		&a.UpdatedOn,
		&a.ArchivedOn,
	}
}

func (a *ProductOption) generateScanArgsWithCount(count *uint64) []interface{} {
	scanArgs := []interface{}{count}
	optionScanArgs := a.generateScanArgs()
	return append(scanArgs, optionScanArgs...)
}

// ProductOptionsResponse is a product option response struct
type ProductOptionsResponse struct {
	ListResponse
	Data []ProductOption `json:"data"`
}

// ProductOptionUpdateInput is a struct to use for updating product options
type ProductOptionUpdateInput struct {
	Name string `json:"name"`
}

// ProductOptionCreationInput is a struct to use for creating product options
type ProductOptionCreationInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// FIXME: this function should be abstracted
func productOptionAlreadyExistsForProduct(db *sqlx.DB, in *ProductOptionCreationInput, productID string) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionExistenceQueryForProductByName, in.Name, productID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

// retrieveProductOptionFromDB retrieves a ProductOption with a given ID from the database
func retrieveProductOptionFromDB(db *sqlx.DB, id uint64) (*ProductOption, error) {
	option := &ProductOption{}
	scanArgs := option.generateScanArgs()
	err := db.QueryRow(productOptionRetrievalQuery, id).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return option, errors.Wrap(err, "Error querying for product")
	}
	return option, err
}

func getProductOptionsForProduct(db *sqlx.DB, productID uint64, queryFilter *QueryFilter) ([]ProductOption, uint64, error) {
	var options []ProductOption
	var count uint64

	query, args := buildProductOptionListQuery(productID, queryFilter)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Error encountered querying for product options")
	}

	defer rows.Close()
	for rows.Next() {
		var option ProductOption
		var queryCount uint64

		scanArgs := option.generateScanArgsWithCount(&queryCount)
		_ = rows.Scan(scanArgs...)

		count = queryCount
		optionValues, err := retrieveProductOptionValueForOptionFromDB(db, option.ID)
		if err != nil {
			return options, 0, errors.Wrap(err, "Error retrieving product option values for option")
		}
		option.Values = optionValues

		options = append(options, option)
	}
	return options, count, nil
}

func buildProductOptionListHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productID := chi.URLParam(req, "product_id")
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		productIDInt, _ := strconv.Atoi(productID)

		options, count, err := getProductOptionsForProduct(db, uint64(productIDInt), queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		optionsResponse := &ProductOptionsResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: count,
			},
			Data: options,
		}
		json.NewEncoder(res).Encode(optionsResponse)
	}
}

func updateProductOptionInDB(db *sqlx.DB, a *ProductOption) error {
	optionUpdateQuery, queryArgs := buildProductOptionUpdateQuery(a)
	err := db.QueryRow(optionUpdateQuery, queryArgs...).Scan(a.generateScanArgs()...)
	return err
}

func buildProductOptionUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductOptionUpdateHandler is a request handler that can update product options
		optionID := chi.URLParam(req, "option_id")
		// eating this error because Chi should validate this for us.
		optionIDInt, _ := strconv.Atoi(optionID)

		// can't update an option that doesn't exist!
		optionExists, err := rowExistsInDB(db, productOptionExistenceQuery, optionID)
		if err != nil || !optionExists {
			respondThatRowDoesNotExist(req, res, "product option", optionID)
			return
		}

		updatedOptionData := &ProductOptionUpdateInput{}
		err = validateRequestInput(req, updatedOptionData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingOption, err := retrieveProductOptionFromDB(db, uint64(optionIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from the database")
			return
		}
		existingOption.Name = updatedOptionData.Name

		err = updateProductOptionInDB(db, existingOption)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option in the database")
			return
		}

		json.NewEncoder(res).Encode(existingOption)
	}
}

func createProductOptionInDB(tx *sql.Tx, a *ProductOption, productID uint64) (*ProductOption, error) {
	var newOptionID int64
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	query, queryArgs := buildProductOptionCreationQuery(a, productID)
	err := tx.QueryRow(query, queryArgs...).Scan(&newOptionID)

	a.ID = uint64(newOptionID)
	a.ProductID = productID
	return a, err
}

func createProductOptionAndValuesInDBFromInput(tx *sql.Tx, in *ProductOptionCreationInput, productID uint64) (*ProductOption, error) {
	newProductOption := &ProductOption{Name: in.Name}
	newProductOption, err := createProductOptionInDB(tx, newProductOption, productID)
	if err != nil {
		return nil, err
	}

	for _, value := range in.Values {
		newOptionValue := ProductOptionValue{
			ProductOptionID: newProductOption.ID,
			Value:           value,
		}
		newOptionValueID, optionCreatedOn, err := createProductOptionValueInDB(tx, &newOptionValue)
		if err != nil {
			return nil, err
		}
		newOptionValue.ID = newOptionValueID
		newOptionValue.CreatedOn = optionCreatedOn
		newProductOption.Values = append(newProductOption.Values, newOptionValue)
	}

	return newProductOption, nil
}

func buildProductOptionCreationHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductOptionCreationHandler is a request handler that can create product options
	return func(res http.ResponseWriter, req *http.Request) {
		productID := chi.URLParam(req, "product_id")
		// eating this error because Chi should validate this for us.
		i, _ := strconv.Atoi(productID)
		productIDInt := uint64(i)

		// can't create an option for a product that doesn't exist!
		productExists, err := rowExistsInDB(db, productExistenceQuery, productID)
		if err != nil || !productExists {
			respondThatRowDoesNotExist(req, res, "product", productID)
			return
		}

		newOptionData := &ProductOptionCreationInput{}
		err = validateRequestInput(req, newOptionData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create an option that already exists!
		optionExists, err := productOptionAlreadyExistsForProduct(db, newOptionData, productID)
		if err != nil || optionExists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option with the name '%s' already exists", newOptionData.Name))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		newProductOption, err := createProductOptionAndValuesInDBFromInput(tx, newOptionData, productIDInt)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "create product option in the database")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newProductOption)
	}
}

func archiveProductOption(db *sqlx.Tx, optionID uint64) error {
	_, err := db.Exec(productOptionDeletionQuery, optionID)
	return err
}

func archiveProductOptionValuesForOption(db *sqlx.Tx, optionID uint64) error {
	_, err := db.Exec(productOptionValuesDeletionQueryByOptionID, optionID)
	return err
}

func buildProductOptionDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductOptionDeletionHandler is a request handler that can delete product options
		optionID := chi.URLParam(req, "option_id")
		// eating this error because Chi should validate this for us.
		optionIDInt, _ := strconv.Atoi(optionID)

		// can't delete an option that doesn't exist!
		optionExists, err := rowExistsInDB(db, productOptionExistenceQuery, optionID)
		if err != nil || !optionExists {
			respondThatRowDoesNotExist(req, res, "product option", optionID)
			return
		}

		tx, err := db.Beginx()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		err = archiveProductOptionValuesForOption(tx, uint64(optionIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "archiving product option values")
			return
		}

		err = archiveProductOption(tx, uint64(optionIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "archiving product options")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
