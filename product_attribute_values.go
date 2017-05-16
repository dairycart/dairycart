package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// ProductAttributeValue represents a products variant attribute values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size, and six ProductAttributeValues,
// One for each color and one for each size.
type ProductAttributeValue struct {
	ID                 int64             `json:"id"`
	ProductAttributeID int64             `json:"product_attribute_id"`
	ProductAttribute   *ProductAttribute `json:"product_attribute"`
	Value              string            `json:"value"`
	ProductsCreated    bool              `json:"products_created"`
	CreatedAt          time.Time         `json:"created_at"`
	ArchivedAt         time.Time         `json:"archived_at"`
}

// CreateProductAttributeValue creates a ProductAttributeValue tied to a ProductAttribute
func CreateProductAttributeValue(db *pg.DB, pav *ProductAttributeValue) (*ProductAttributeValue, error) {
	err := db.Insert(pav)
	return pav, err
}

// RetrieveProductAttributeValue retrieves a ProductAttributeValue with a given ID from the database
func RetrieveProductAttributeValue(db *pg.DB, id int64) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}
	productAttributeValue := db.Model(pav).
		Where("id = ?", id).
		Where("product_attribute_value.archived_at is null")

	err := productAttributeValue.Select()
	return pav, err
}

// ProductAttributeValueExists checks for the existence of a given ProductAttributeValue in the database
func ProductAttributeValueExists(db *pg.DB, id int64) (bool, error) {
	count, err := db.Model(&ProductAttributeValue{}).Where("id = ?", id).Where("archived_at is null").Count()

	return count == 1, err
}

func buildProductAttributeValueCreationHandler(db *pg.DB) func(res http.ResponseWriter, req *http.Request) {
	// ProductAttributeValueCreationHandler is a product creation handler
	return func(res http.ResponseWriter, req *http.Request) {
		providedAttributeID := mux.Vars(req)["attribute_id"]

		// we can eat this error because Mux takes care of validating route params for us
		attributeID, _ := strconv.ParseInt(providedAttributeID, 10, 64)

		productAttribueExists := ProductAttributeExists(db, attributeID)
		if !productAttribueExists {
			respondToInvalidRequest(nil, fmt.Sprintf("No matching product attribute for ID: %d", attributeID), res)
			return
		}

		newProductAttributeValue := &ProductAttributeValue{}
		bodyIsInvalid := ensureRequestBodyValidity(res, req, newProductAttributeValue)
		if bodyIsInvalid {
			return
		}
		newProductAttributeValue.ProductAttributeID = attributeID

		// We don't want API consumers to be able to override this value
		newProductAttributeValue.ProductsCreated = false

		_, err := CreateProductAttributeValue(db, newProductAttributeValue)
		if err != nil {
			errorString := fmt.Sprintf("error inserting product into database: %v", err)
			log.Println(errorString)
			http.Error(res, errorString, http.StatusBadRequest)
			return
		}
	}
}
