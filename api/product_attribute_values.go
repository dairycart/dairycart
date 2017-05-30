package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ProductAttributeValue represents a products variant attribute values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size, and six ProductAttributeValues,
// One for each color and one for each size.
type ProductAttributeValue struct {
	ID                 int64       `json:"id"`
	ProductAttributeID int64       `json:"product_attribute_id"`
	Value              string      `json:"value"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          pq.NullTime `json:"updated_at"`
	ArchivedAt         pq.NullTime `json:"-"`
}

func (pav *ProductAttributeValue) generateScanArgs() []interface{} {
	return []interface{}{
		&pav.ID,
		&pav.ProductAttributeID,
		&pav.Value,
		&pav.CreatedAt,
		&pav.UpdatedAt,
		&pav.ArchivedAt,
	}
}

// ProductAttributeValueUpdateInput is a struct to use for updating product attribute values
type ProductAttributeValueUpdateInput struct {
	Value string `json:"value"`
}

func validateProductAttributeValueUpdateInput(req *http.Request) (*ProductAttributeValueUpdateInput, error) {
	i := &ProductAttributeValueUpdateInput{}
	json.NewDecoder(req.Body).Decode(i)

	s := structs.New(i)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product attribute body")
	}

	return i, nil
}

// retrieveProductAttributeValue retrieves a ProductAttributeValue with a given ID from the database
func retrieveProductAttributeValueFromDB(db *sql.DB, id int64) (*ProductAttributeValue, error) {
	v := &ProductAttributeValue{}
	query := buildProductAttributeValueRetrievalQuery(id)
	err := db.QueryRow(query, id).Scan(v.generateScanArgs()...)
	if err == sql.ErrNoRows {
		return v, errors.Wrap(err, "Error querying for product attribute values")
	}
	return v, err
}

func loadProductAttributeValueInput(req *http.Request) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(pav)

	s := structs.New(pav)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for product attribute value body")
	}

	return pav, err
}

func updateProductAttributeValueInDB(db *sql.DB, v *ProductAttributeValue) error {
	valueUpdateQuery, queryArgs := buildProductAttributeValueUpdateQuery(v)
	err := db.QueryRow(valueUpdateQuery, queryArgs...).Scan(v.generateScanArgs()...)
	return err
}

func buildProductAttributeValueUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductAttributeValueUpdateHandler is a request handler that can update product attribute values
		reqVars := mux.Vars(req)
		attributeID := reqVars["attribute_id"]
		attributeValueID := reqVars["attribute_value_id"]
		// eating these errors because Mux should validate these for us.
		attributeValueIDInt, _ := strconv.Atoi(attributeValueID)

		// can't update an attribute that doesn't exist!
		attributeExists, err := rowExistsInDB(db, "product_attributes", "id", attributeID)
		if err != nil || !attributeExists {
			respondThatRowDoesNotExist(req, res, "product attribute", attributeID)
			return
		}

		// can't update an attribute value that doesn't exist!
		attributeValueExists, err := rowExistsInDB(db, "product_attribute_values", "id", attributeValueID)
		if err != nil || !attributeValueExists {
			respondThatRowDoesNotExist(req, res, "product attribute value", attributeValueID)
			return
		}

		updatedValueData, err := validateProductAttributeValueUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingAttributeValue, err := retrieveProductAttributeValueFromDB(db, int64(attributeValueIDInt))
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve product attribute from the database")
			return
		}
		existingAttributeValue.Value = updatedValueData.Value

		err = updateProductAttributeValueInDB(db, existingAttributeValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product attribute in the database")
			return
		}

		json.NewEncoder(res).Encode(existingAttributeValue)

	}
}

// createProductAttributeValueInDB creates a ProductAttributeValue tied to a ProductAttribute
func createProductAttributeValueInDB(db *sql.DB, v *ProductAttributeValue) (int64, error) {
	var newAttributeValueID int64
	query, args := buildProductAttributeValueCreationQuery(v)
	err := db.QueryRow(query, args...).Scan(&newAttributeValueID)
	return newAttributeValueID, err
}

func buildProductAttributeValueCreationHandler(db *sql.DB) http.HandlerFunc {
	// productAttributeValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		attributeID := mux.Vars(req)["attribute_id"]

		// we can eat this error because Mux takes care of validating route params for us
		attributeIDInt, _ := strconv.ParseInt(attributeID, 10, 64)

		productAttributeExists, err := rowExistsInDB(db, "product_attributes", "id", attributeID)
		if err != nil || !productAttributeExists {
			respondThatRowDoesNotExist(req, res, "product attribute", attributeID)
			return
		}

		newProductAttributeValue, err := loadProductAttributeValueInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		newProductAttributeValue.ProductAttributeID = attributeIDInt

		newProductAttributeValueID, err := createProductAttributeValueInDB(db, newProductAttributeValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProductAttributeValue.ID = newProductAttributeValueID

		json.NewEncoder(res).Encode(newProductAttributeValue)
	}
}
