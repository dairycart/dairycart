package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
)

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

func buildProductsFromOptions(input *models.ProductCreationInput, createdOptions []models.ProductOption) (toCreate []*models.Product) {
	// lovingly borrowed from:
	//     https://stackoverflow.com/questions/29002724/implement-ruby-style-cartesian-product-in-go
	// NextIndex sets ix to the lexicographically next value,
	// such that for each i > 0, 0 <= ix[i] < lens(i).
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
	var optionData [][]optionPlaceholder
	for _, o := range createdOptions {
		var newOptions []optionPlaceholder
		for _, v := range o.Values {
			summary := fmt.Sprintf("%s: %s", o.Name, v.Value)
			ph := optionPlaceholder{
				ID:            v.ID,
				Summary:       summary,
				Value:         v.Value,
				OriginalValue: v,
			}
			newOptions = append(newOptions, ph)
		}
		optionData = append(optionData, newOptions)
	}

	for ix := make([]int, len(optionData)); ix[0] < len(optionData[0]); next(ix, optionData) {
		var skuPrefixParts, optionSummaryParts []string
		var originalValues []models.ProductOptionValue
		for j, k := range ix {
			optionSummaryParts = append(optionSummaryParts, optionData[j][k].Summary)
			skuPrefixParts = append(skuPrefixParts, strings.ToLower(optionData[j][k].Value))
			originalValues = append(originalValues, optionData[j][k].OriginalValue)
		}

		productTemplate := newProductFromCreationInput(input)
		productTemplate.OptionSummary = strings.Join(optionSummaryParts, ", ")
		productTemplate.SKU = fmt.Sprintf("%s_%s", input.SKU, strings.Join(skuPrefixParts, "_"))
		productTemplate.ApplicableOptionValues = originalValues
		toCreate = append(toCreate, productTemplate)

	}
	return toCreate
}

func buildProductOptionListHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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

func buildProductOptionUpdateHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	// ProductOptionUpdateHandler is a request handler that can update product options
	return func(res http.ResponseWriter, req *http.Request) {
		optionIDStr := chi.URLParam(req, "option_id")
		// eating this error because the router should have ensured this is an integer
		optionID, _ := strconv.ParseUint(optionIDStr, 10, 64)

		updatedOptionData := &models.ProductOption{}
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

		mergo.MergeWithOverwrite(existingOption, updatedOptionData)

		updatedOn, err := client.UpdateProductOption(db, existingOption)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product option in the database")
			return
		}
		existingOption.UpdatedOn = &models.Dairytime{Time: updatedOn}

		values, err := client.GetProductOptionValuesForOption(db, existingOption.ID)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product option from the database")
			return
		}
		existingOption.Values = values

		json.NewEncoder(res).Encode(existingOption)
	}
}

func createProductOptionAndValuesInDBFromInput(tx *sql.Tx, in models.ProductOptionCreationInput, productRootID uint64, client database.Storer) (models.ProductOption, error) {
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
		newOptionValue.ID, newOptionValue.CreatedOn, err = client.CreateProductOptionValue(tx, &newOptionValue)
		if err != nil {
			return models.ProductOption{}, err
		}
		newProductOption.Values = append(newProductOption.Values, newOptionValue)
	}

	return *newProductOption, nil
}

func buildProductOptionCreationHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
	// ProductOptionCreationHandler is a request handler that can create product options
	return func(res http.ResponseWriter, req *http.Request) {
		productRootIDStr := chi.URLParam(req, "product_root_id")
		// eating this error because the router should have ensured this is an integer
		productRootID, _ := strconv.ParseUint(productRootIDStr, 10, 64)

		newOptionData := models.ProductOptionCreationInput{}
		err := validateRequestInput(req, &newOptionData)
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

func buildProductOptionDeletionHandler(db *sql.DB, client database.Storer) http.HandlerFunc {
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
		existingOption.ArchivedOn = &models.Dairytime{Time: archivedOn}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "close out transaction")
			return
		}

		json.NewEncoder(res).Encode(existingOption)
	}
}
