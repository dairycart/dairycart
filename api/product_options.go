package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	productOptionsHeaders = `id,
		name,
		product_root_id,
		created_on,
		updated_on,
		archived_on
	`
	productOptionExistenceQuery                 = `SELECT EXISTS(SELECT 1 FROM product_options WHERE id = $1 AND archived_on IS NULL)`
	productOptionRetrievalQuery                 = `SELECT * FROM product_options WHERE id = $1`
	productOptionExistenceQueryForProductByName = `SELECT EXISTS(SELECT 1 FROM product_options WHERE name = $1 AND product_root_id = $2 and archived_on IS NULL)`
	productOptionDeletionQuery                  = `UPDATE product_options SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`
	productOptionValuesDeletionQueryByOptionID  = `UPDATE product_option_values SET archived_on = NOW() WHERE product_option_id = $1 AND archived_on IS NULL`
)

// ProductOption represents a products variant options. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size.
type ProductOption struct {
	DBRow
	Name          string               `json:"name"`
	ProductRootID uint64               `json:"product_root_id"`
	Values        []ProductOptionValue `json:"values"`
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

type simpleProductOption struct {
	OptionSummary string
	SKUPostfix    string
}

func generateCartesianProductForOptions(inputOptions []ProductOptionCreationInput) []simpleProductOption {
	/*
		Some notes about this function:

		It's probably hilariously expensive to run, like O(n^(log(n)³)) or some other equally absurd thing
		I based this off a stackoverflow post and didn't go to college. I've tried to use anonymous structs where
		I could so I don't have data structures floating around that exist solely for this function, and
		also tried to name things as clearly as possible. But it still kind of just _feels_ messy, so forgive me,
		Rob Pike. I have taken your beautiful language and violated it with my garbage brain
	*/

	// lovingly borrowed from:
	//     https://stackoverflow.com/questions/29002724/implement-ruby-style-cartesian-product-in-go
	// NextIndex sets ix to the lexicographically next value,
	// such that for each i>0, 0 <= ix[i] < lens(i).
	nextIndex := func(ix []int, sl [][]struct{ Summary, Value string }) {
		for j := len(ix) - 1; j >= 0; j-- {
			ix[j]++
			if j == 0 || ix[j] < len(sl[j]) {
				return
			}
			ix[j] = 0
		}
	}

	// meat & potatoes starts here
	optionStrings := [][]struct{ Summary, Value string }{}
	for _, o := range inputOptions {
		newOptions := []struct{ Summary, Value string }{}
		for _, value := range o.Values {
			newOptions = append(newOptions, struct{ Summary, Value string }{Summary: fmt.Sprintf("%s: %s", o.Name, value), Value: value})
		}
		optionStrings = append(optionStrings, newOptions)
	}

	output := []simpleProductOption{}
	for ix := make([]int, len(optionStrings)); ix[0] < len(optionStrings[0]); nextIndex(ix, optionStrings) {
		var skuPrefixParts []string
		var optionSummaryParts []string
		for j, k := range ix {
			optionSummaryParts = append(optionSummaryParts, optionStrings[j][k].Summary)
			skuPrefixParts = append(skuPrefixParts, strings.ToLower(optionStrings[j][k].Value))
		}
		output = append(output, simpleProductOption{
			OptionSummary: strings.Join(optionSummaryParts, ", "),
			SKUPostfix:    strings.Join(skuPrefixParts, "_"),
		})
	}

	return output
}

// FIXME: this function should be abstracted
func productOptionAlreadyExistsForProduct(db *sqlx.DB, in *ProductOptionCreationInput, productRootID string) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionExistenceQueryForProductByName, in.Name, productRootID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

// retrieveProductOptionFromDB retrieves a ProductOption with a given ID from the database
func retrieveProductOptionFromDB(db *sqlx.DB, id uint64) (*ProductOption, error) {
	option := &ProductOption{}
	err := db.QueryRowx(productOptionRetrievalQuery, id).StructScan(option)
	if err == sql.ErrNoRows {
		return option, errors.Wrap(err, "Error querying for product")
	}
	return option, err
}

func getProductOptionsForProduct(db *sqlx.DB, productRootID uint64, queryFilter *QueryFilter) ([]ProductOption, error) {
	var options []ProductOption

	query, args := buildProductOptionListQuery(productRootID, queryFilter)
	err := db.Select(&options, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for product options")
	}

	for _, option := range options {
		optionValues, err := retrieveProductOptionValuesForOptionFromDB(db, option.ID)
		if err != nil {
			return options, errors.Wrap(err, "Error retrieving product option values for option")
		}
		option.Values = optionValues
	}
	return options, nil
}

func buildProductOptionListHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productRootID := chi.URLParam(req, "product_root_id")
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		productRootIDInt, _ := strconv.Atoi(productRootID)

		// FIXME: this will return the count of all options, not the options for a given product root
		count, err := getRowCount(db, "product_options", queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of product options from the database")
			return
		}

		options, err := getProductOptionsForProduct(db, uint64(productRootIDInt), queryFilter)
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

func updateProductOptionInDB(db *sqlx.DB, a *ProductOption) (time.Time, error) {
	var updatedOn time.Time
	optionUpdateQuery, queryArgs := buildProductOptionUpdateQuery(a)
	err := db.QueryRow(optionUpdateQuery, queryArgs...).Scan(&updatedOn)
	return updatedOn, err
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

		noop := func(...interface{}) {
			return
		}

		existingOption, err := retrieveProductOptionFromDB(db, uint64(optionIDInt))
		if err != nil {
			errStr := err.Error()
			noop(errStr)
			notifyOfInternalIssue(res, err, "retrieve product option from the database")
			return
		}
		existingOption.Name = updatedOptionData.Name

		optionUpdatedOn, err := updateProductOptionInDB(db, existingOption)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option in the database")
			return
		}
		existingOption.UpdatedOn = NullTime{pq.NullTime{Time: optionUpdatedOn, Valid: true}}

		json.NewEncoder(res).Encode(existingOption)
	}
}

func createProductOptionInDB(tx *sql.Tx, o *ProductOption, productRootID uint64) (uint64, time.Time, error) {
	var newOptionID uint64
	var createdOn time.Time
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	query, queryArgs := buildProductOptionCreationQuery(o, productRootID)
	err := tx.QueryRow(query, queryArgs...).Scan(&newOptionID, &createdOn)

	return newOptionID, createdOn, err
}

func createProductOptionAndValuesInDBFromInput(tx *sql.Tx, in *ProductOptionCreationInput, productRootID uint64) (*ProductOption, error) {
	newProductOption := &ProductOption{Name: in.Name, ProductRootID: productRootID}
	newProductOptionID, newProductOptionCreatedOn, err := createProductOptionInDB(tx, newProductOption, productRootID)
	if err != nil {
		return nil, err
	}
	newProductOption.ID = newProductOptionID
	newProductOption.CreatedOn = newProductOptionCreatedOn

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
		productRootID := chi.URLParam(req, "product_root_id")
		// eating this error because Chi should validate this for us.
		i, _ := strconv.Atoi(productRootID)
		productRootIDInt := uint64(i)

		// can't create an option for a product that doesn't exist!
		productRootExists, err := rowExistsInDB(db, productRootExistenceQuery, productRootID)
		if err != nil || !productRootExists {
			respondThatRowDoesNotExist(req, res, "product", productRootID)
			return
		}

		newOptionData := &ProductOptionCreationInput{}
		err = validateRequestInput(req, newOptionData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create an option that already exists!
		optionExists, err := productOptionAlreadyExistsForProduct(db, newOptionData, productRootID)
		if err != nil || optionExists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option with the name '%s' already exists", newOptionData.Name))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		newProductOption, err := createProductOptionAndValuesInDBFromInput(tx, newOptionData, productRootIDInt)
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
