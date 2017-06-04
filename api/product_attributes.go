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
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID                  int64                    `json:"id"`
	Name                string                   `json:"name"`
	ProductProgenitorID int64                    `json:"product_progenitor_id"`
	Values              []*ProductAttributeValue `json:"values"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           pq.NullTime              `json:"-"`
	ArchivedAt          pq.NullTime              `json:"-"`
}

func (a *ProductAttribute) generateScanArgs() []interface{} {
	return []interface{}{
		&a.ID,
		&a.Name,
		&a.ProductProgenitorID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.ArchivedAt,
	}
}

func (a *ProductAttribute) generateScanArgsWithCount(count *uint64) []interface{} {
	scanArgs := []interface{}{count}
	attributeScanArgs := a.generateScanArgs()
	return append(scanArgs, attributeScanArgs...)
}

// ProductAttributesResponse is a product attribute response struct
type ProductAttributesResponse struct {
	ListResponse
	Data []ProductAttribute `json:"data"`
}

// ProductAttributeUpdateInput is a struct to use for updating product attributes
type ProductAttributeUpdateInput struct {
	Name string `json:"name"`
}

// ProductAttributeCreationInput is a struct to use for creating product attributes
type ProductAttributeCreationInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func productAttributeAlreadyExistsForProgenitor(db *sql.DB, in *ProductAttributeCreationInput, progenitorID string) (bool, error) {
	var exists string

	query := buildProductAttributeExistenceQueryForProductByName(in.Name, progenitorID)
	err := db.QueryRow(query, in.Name, progenitorID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

// retrieveProductAttributeFromDB retrieves a ProductAttribute with a given ID from the database
func retrieveProductAttributeFromDB(db *sql.DB, id int64) (*ProductAttribute, error) {
	attribute := &ProductAttribute{}
	scanArgs := attribute.generateScanArgs()
	query := buildProductAttributeRetrievalQuery(id)
	err := db.QueryRow(query, id).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return attribute, errors.Wrap(err, "Error querying for product")
	}
	return attribute, err
}

func getProductAttributesForProgenitor(db *sql.DB, progenitorID string, queryFilter *QueryFilter) ([]ProductAttribute, uint64, error) {
	var attributes []ProductAttribute
	var count uint64

	query := buildProductAttributeListQueryWithCount(progenitorID, queryFilter)
	rows, err := db.Query(query, progenitorID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Error encountered querying for product attributes")
	}

	defer rows.Close()
	for rows.Next() {
		var attribute ProductAttribute
		var queryCount uint64

		scanArgs := attribute.generateScanArgsWithCount(&queryCount)
		_ = rows.Scan(scanArgs...)

		count = queryCount
		attributeValues, err := retrieveProductAttributeValueForAttributeFromDB(db, attribute.ID)
		if err != nil {
			return attributes, 0, errors.Wrap(err, "Error retrieving product attribute values for attribute")
		}
		attribute.Values = attributeValues

		attributes = append(attributes, attribute)
	}
	return attributes, count, nil
}

func buildProductAttributeListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		progenitorID := mux.Vars(req)["progenitor_id"]
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)

		attributes, count, err := getProductAttributesForProgenitor(db, progenitorID, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		attributesResponse := &ProductAttributesResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: count,
			},
			Data: attributes,
		}
		json.NewEncoder(res).Encode(attributesResponse)
	}
}

func validateProductAttributeUpdateInput(req *http.Request) (*ProductAttributeUpdateInput, error) {
	i := &ProductAttributeUpdateInput{}
	json.NewDecoder(req.Body).Decode(i)

	s := structs.New(i)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product attribute body")
	}

	return i, nil
}

func updateProductAttributeInDB(db *sql.DB, a *ProductAttribute) error {
	attributeUpdateQuery, queryArgs := buildProductAttributeUpdateQuery(a)
	err := db.QueryRow(attributeUpdateQuery, queryArgs...).Scan(a.generateScanArgs()...)
	return err
}

func buildProductAttributeUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductAttributeUpdateHandler is a request handler that can update product attributes
		reqVars := mux.Vars(req)
		attributeID := reqVars["attribute_id"]
		// eating this error because Mux should validate this for us.
		attributeIDInt, _ := strconv.Atoi(attributeID)

		// can't update an attribute that doesn't exist!
		attributeExists, err := rowExistsInDB(db, "product_attributes", "id", attributeID)
		if err != nil || !attributeExists {
			respondThatRowDoesNotExist(req, res, "product attribute", attributeID)
			return
		}

		updatedAttributeData, err := validateProductAttributeUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingAttribute, err := retrieveProductAttributeFromDB(db, int64(attributeIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product attribute from the database")
			return
		}
		existingAttribute.Name = updatedAttributeData.Name

		err = updateProductAttributeInDB(db, existingAttribute)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product attribute in the database")
			return
		}

		json.NewEncoder(res).Encode(existingAttribute)

	}
}

func validateProductAttributeCreationInput(req *http.Request) (*ProductAttributeCreationInput, error) {
	i := &ProductAttributeCreationInput{}
	json.NewDecoder(req.Body).Decode(i)

	s := structs.New(i)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product attribute body")
	}

	return i, nil
}

func createProductAttributeInDB(tx *sql.Tx, a *ProductAttribute) (*ProductAttribute, error) {
	var newAttributeID int64
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	query, queryArgs := buildProductAttributeCreationQuery(a)
	err := tx.QueryRow(query, queryArgs...).Scan(&newAttributeID)

	a.ID = newAttributeID
	return a, err
}

func createProductAttributeAndValuesInDBFromInput(tx *sql.Tx, in *ProductAttributeCreationInput, progenitorID int64) (*ProductAttribute, error) {
	newProductAttribute := &ProductAttribute{
		Name:                in.Name,
		ProductProgenitorID: progenitorID,
	}

	newProductAttribute, err := createProductAttributeInDB(tx, newProductAttribute)
	if err != nil {
		return nil, err
	}

	for _, value := range in.Values {
		newAttributeValue := &ProductAttributeValue{
			ProductAttributeID: newProductAttribute.ID,
			Value:              value,
		}
		newAttributeValueID, err := createProductAttributeValueInDB(tx, newAttributeValue)
		if err != nil {
			return nil, err
		}
		newAttributeValue.ID = newAttributeValueID
		newProductAttribute.Values = append(newProductAttribute.Values, newAttributeValue)
	}

	return newProductAttribute, nil
}

func buildProductAttributeCreationHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductAttributeCreationHandler is a request handler that can create product attributes
		progenitorID := mux.Vars(req)["progenitor_id"]
		// eating this error because Mux should validate this for us.
		progenitorIDInt, _ := strconv.Atoi(progenitorID)

		// can't create an attribute for a product progenitor that doesn't exist!
		progenitorExists, err := rowExistsInDB(db, "product_progenitors", "id", progenitorID)
		if err != nil || !progenitorExists {
			respondThatRowDoesNotExist(req, res, "product progenitor", progenitorID)
			return
		}

		newAttributeData, err := validateProductAttributeCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create an attribute that already exist!
		attributeExists, err := productAttributeAlreadyExistsForProgenitor(db, newAttributeData, progenitorID)
		if err != nil || attributeExists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product attribute with the name `%s` already exists", newAttributeData.Name))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			notifyOfInternalIssue(res, err, "starting a new transaction")
			return
		}

		newProductAttribute, err := createProductAttributeAndValuesInDBFromInput(tx, newAttributeData, int64(progenitorIDInt))
		if err != nil {
			tx.Rollback()
			notifyOfInternalIssue(res, err, "create product attribute in the database")
			return
		}

		err = tx.Commit()
		if err != nil {
			notifyOfInternalIssue(res, err, "closing out transaction")
			return
		}

		json.NewEncoder(res).Encode(newProductAttribute)
	}
}
