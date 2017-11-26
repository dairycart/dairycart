package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
)

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
	IDs            []uint64
	OptionSummary  string
	SKUPostfix     string
	OriginalValues []models.ProductOptionValue
}

type optionPlaceholder struct {
	ID            uint64
	Summary       string
	Value         string
	OriginalValue models.ProductOptionValue
}

// FIXME: don't use pointers here
func generateCartesianProductForOptions(inputOptions []models.ProductOption) []simpleProductOption {
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
	next := func(ix []int, sl [][]optionPlaceholder) {
		for j := len(ix) - 1; j >= 0; j-- {
			ix[j]++
			if j == 0 || ix[j] < len(sl[j]) {
				return
			}
			ix[j] = 0
		}
	}

	// meat & potatoes starts here
	optionData := [][]optionPlaceholder{}
	for _, o := range inputOptions {
		newOptions := []optionPlaceholder{}
		for _, v := range o.Values {
			ph := optionPlaceholder{
				ID:            v.ID,
				Summary:       fmt.Sprintf("%s: %s", o.Name, v.Value),
				Value:         v.Value,
				OriginalValue: v,
			}
			newOptions = append(newOptions, ph)
		}
		optionData = append(optionData, newOptions)
	}

	output := []simpleProductOption{}
	for ix := make([]int, len(optionData)); ix[0] < len(optionData[0]); next(ix, optionData) {
		var ids []uint64
		var skuPrefixParts []string
		var optionSummaryParts []string
		var originalValues []models.ProductOptionValue
		for j, k := range ix {
			ids = append(ids, optionData[j][k].ID)
			optionSummaryParts = append(optionSummaryParts, optionData[j][k].Summary)
			skuPrefixParts = append(skuPrefixParts, strings.ToLower(optionData[j][k].Value))
			originalValues = append(originalValues, optionData[j][k].OriginalValue)
		}
		output = append(output, simpleProductOption{
			IDs:            ids,
			OptionSummary:  strings.Join(optionSummaryParts, ", "),
			SKUPostfix:     strings.Join(skuPrefixParts, "_"),
			OriginalValues: originalValues,
		})
	}

	return output
}

func buildProductOptionListHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		// eating this error because the router should have ensured this is an integer
		productRootID, _ := strconv.ParseUint(productRootIDStr, 10, 64)
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)

		// FIXME: this will return the count of all options, not the options for a given product root
		count, err := client.GetProductOptionCount(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of product options from the database")
			return
		}

		options, err := client.GetProductOptionsByProductRootID(db, productRootID)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		optionsResponse := &ListResponse{
			Page:  queryFilter.Page,
			Limit: queryFilter.Limit,
			Count: count,
			Data:  options,
		}
		json.NewEncoder(res).Encode(optionsResponse)
	}
}

func buildProductOptionUpdateHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductOptionUpdateHandler is a request handler that can update product options
	return func(res http.ResponseWriter, req *http.Request) {
		optionIDStr := chi.URLParam(req, "option_id")
		// eating this error because the router should have ensured this is an integer
		optionID, _ := strconv.ParseUint(optionIDStr, 10, 64)

		updatedOptionData := &ProductOptionUpdateInput{}
		err := validateRequestInput(req, updatedOptionData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingOption, err := client.GetProductOption(db, optionID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product option", optionIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from database")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(updatedOptionData, existingOption)

		updatedOn, err := client.UpdateProductOption(db, existingOption)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option in the database")
			return
		}
		existingOption.UpdatedOn = models.NullTime{NullTime: pq.NullTime{Time: updatedOn, Valid: true}}

		values, err := client.GetProductOptionValuesForOption(db, existingOption.ID)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from the database")
			return
		}
		existingOption.Values = values

		json.NewEncoder(res).Encode(existingOption)
	}
}

func createProductOptionAndValuesInDBFromInput(tx *sql.Tx, in *ProductOptionCreationInput, productRootID uint64, client storage.Storer) (models.ProductOption, error) {
	var err error
	newProductOption := &models.ProductOption{Name: in.Name, ProductRootID: productRootID}
	newProductOption.ID, newProductOption.CreatedOn, err = client.CreateProductOption(tx, newProductOption)
	if err != nil {
		return models.ProductOption{}, err
	}

	for _, value := range in.Values {
		newOptionValue := models.ProductOptionValue{
			ProductOptionID: newProductOption.ID,
			Value:           value,
		}
		newOptionValueID, optionCreatedOn, err := client.CreateProductOptionValue(tx, &newOptionValue)
		if err != nil {
			return models.ProductOption{}, err
		}
		newOptionValue.ID = newOptionValueID
		newOptionValue.CreatedOn = optionCreatedOn
		newProductOption.Values = append(newProductOption.Values, newOptionValue)
	}

	return *newProductOption, nil
}

func buildProductOptionCreationHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductOptionCreationHandler is a request handler that can create product options
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		// eating this error because the router should have ensured this is an integer
		productRootID, _ := strconv.ParseUint(productRootIDStr, 10, 64)

		newOptionData := &ProductOptionCreationInput{}
		err := validateRequestInput(req, newOptionData)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create an option for a product that doesn't exist!
		productRootExists, err := client.ProductRootExists(db, productRootID)
		if err == sql.ErrNoRows || !productRootExists {
			respondThatRowDoesNotExist(req, res, "product root", productRootIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product root from database")
			return
		}

		// can't create an option that already exists!
		optionExists, err := client.ProductOptionWithNameExistsForProductRoot(db, newOptionData.Name, productRootID)
		if optionExists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product option with the name '%s' already exists", newOptionData.Name))
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from database")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		newProductOption, err := createProductOptionAndValuesInDBFromInput(tx, newOptionData, productRootID, client)
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "create product option in the database")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newProductOption)
	}
}

func buildProductOptionDeletionHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// ProductOptionDeletionHandler is a request handler that can delete product options
	return func(res http.ResponseWriter, req *http.Request) {
		optionIDStr := chi.URLParam(req, "option_id")
		// eating this error because the router should have ensured this is an integer
		optionID, _ := strconv.ParseUint(optionIDStr, 10, 64)

		// can't delete an option that doesn't exist!
		existingOption, err := client.GetProductOption(db, optionID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "product option", optionIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from database")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		_, err = client.ArchiveProductOptionValuesForOption(tx, optionID)
		if err != nil {
			notifyOfInternalIssue(res, err, "archiving product option values")
			return
		}

		archivedOn, err := client.DeleteProductOption(tx, optionID)
		if err != nil {
			notifyOfInternalIssue(res, err, "archiving product options")
			return
		}
		existingOption.ArchivedOn = models.NullTime{NullTime: pq.NullTime{Time: archivedOn, Valid: true}}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		json.NewEncoder(res).Encode(existingOption)
	}
}
