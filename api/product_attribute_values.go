package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	productAttributeValueRetrievalQuery = `SELECT * FROM product_attribute_values WHERE id = $1 AND archived_at IS NULL`
	productAttributeValueCreationQuery  = `INSERT INTO product_attribute_values ("product_attribute_id", "value") VALUES ($1, $2);`
)

// ProductAttributeValue represents a products variant attribute values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size, and six ProductAttributeValues,
// One for each color and one for each size.
type ProductAttributeValue struct {
	ID                 int64       `json:"id"`
	ProductAttributeID int64       `json:"product_attribute_id"`
	Value              string      `json:"value"`
	ProductsCreated    bool        `json:"products_created"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          pq.NullTime `json:"updated_at"`
	ArchivedAt         pq.NullTime `json:"-"`
}

func (pav ProductAttributeValue) generateScanArgs() []interface{} {
	return []interface{}{
		&pav.ID,
		&pav.ProductAttributeID,
		&pav.Value,
		&pav.ProductsCreated,
		&pav.CreatedAt,
		&pav.UpdatedAt,
		&pav.ArchivedAt,
	}
}

// createProductAttributeValue creates a ProductAttributeValue tied to a ProductAttribute
func createProductAttributeValue(db *sql.DB, pav *ProductAttributeValue) (*ProductAttributeValue, error) {
	_, err := db.Exec(productAttributeValueCreationQuery, pav.ProductAttributeID, pav.Value)
	return nil, err
}

// retrieveProductAttributeValue retrieves a ProductAttributeValue with a given ID from the database
func retrieveProductAttributeValue(db *sql.DB, id int64) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}

	err := db.QueryRow(productAttributeValueRetrievalQuery, id).Scan(pav.generateScanArgs()...)
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

func buildProductAttributeValueCreationHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	// productAttributeValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		attributeID := mux.Vars(req)["attribute_id"]

		// we can eat this error because Mux takes care of validating route params for us
		attributeIDInt, _ := strconv.ParseInt(attributeID, 10, 64)

		productAttribueExists, err := rowExistsInDB(db, "product_attributes", "id", attributeID)
		if err != nil || !productAttribueExists {
			respondThatRowDoesNotExist(req, res, "product attribute", "ID", attributeID)
			return
		}

		newProductAttributeValue, err := loadProductAttributeValueInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}
		newProductAttributeValue.ProductAttributeID = attributeIDInt

		// We don't want API consumers to be able to override this value
		newProductAttributeValue.ProductsCreated = false

		_, err = createProductAttributeValue(db, newProductAttributeValue)
		if err != nil {
			errorString := fmt.Sprintf("error inserting product into database: %v", err)
			log.Println(errorString)
			http.Error(res, errorString, http.StatusBadRequest)
			return
		}
	}
}
