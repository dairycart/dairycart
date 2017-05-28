package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func (pav ProductAttributeValue) generateScanArgs() []interface{} {
	return []interface{}{
		&pav.ID,
		&pav.ProductAttributeID,
		&pav.Value,
		&pav.CreatedAt,
		&pav.UpdatedAt,
		&pav.ArchivedAt,
	}
}

// retrieveProductAttributeValue retrieves a ProductAttributeValue with a given ID from the database
func retrieveProductAttributeValue(db *sql.DB, id int64) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}
	query := buildProductAttributeValueRetrievalQuery(id)
	err := db.QueryRow(query, id).Scan(pav.generateScanArgs()...)
	if err == sql.ErrNoRows {
		return pav, errors.Wrap(err, "Error querying for product attribute values")
	}

	return pav, err
}

func loadProductAttributeValueInput(req *http.Request) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(pav)

	return pav, err
}

// createProductAttributeValueInDB creates a ProductAttributeValue tied to a ProductAttribute
func createProductAttributeValueInDB(db *sql.DB, pav *ProductAttributeValue) (*ProductAttributeValue, error) {
	query, args := buildProductAttributeValueCreationQuery(pav)
	err := db.QueryRow(query, args...).Scan(pav.generateScanArgs()...)
	return pav, err
}

func buildProductAttributeValueCreationHandler(db *sql.DB) http.HandlerFunc {
	// productAttributeValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		attributeID := mux.Vars(req)["attribute_id"]

		// we can eat this error because Mux takes care of validating route params for us
		attributeIDInt, _ := strconv.ParseInt(attributeID, 10, 64)

		productAttributeExists, err := rowExistsInDB(db, "product_attributes", "id", attributeID)
		if err != nil || !productAttributeExists {
			respondThatRowDoesNotExist(req, res, "product_attribute", attributeID)
			return
		}

		newProductAttributeValue, err := loadProductAttributeValueInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		newProductAttributeValue.ProductAttributeID = attributeIDInt

		newProductAttributeValue, err = createProductAttributeValueInDB(db, newProductAttributeValue)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}

		json.NewEncoder(res).Encode(newProductAttributeValue)
	}
}
