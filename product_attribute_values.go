package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
	Active             bool              `json:"active"`
	CreatedAt          time.Time         `json:"created_at"`
	ArchivedAt         time.Time         `json:"archived_at"`
}

// CreateProductAttributeValue creates a ProductAttributeValue tied to a ProductAttribute
func CreateProductAttributeValue(pav *ProductAttributeValue) (*ProductAttributeValue, error) {
	err := db.Insert(pav)
	return pav, err
}

// RetrieveProductAttributeValue retrieves a ProductAttributeValue with a given ID from the database
func RetrieveProductAttributeValue(id int64) (*ProductAttributeValue, error) {
	pav := &ProductAttributeValue{}
	productAttributeValue := db.Model(pav).
		Where("id = ?", id).
		Where("product_attribute_value.archived_at is null")

	err := productAttributeValue.Select()
	return pav, err
}

// ProductAttributeValueCreationHandler is a product creation handler
func ProductAttributeValueCreationHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	providedAttributeID := vars["attribute_id"]

	attributeID, err := strconv.ParseInt(providedAttributeID, 10, 64)
	if err != nil {
		errorString := fmt.Sprintf("Error encountered parsing base product ID: %v", err)
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

	newProductAttributeValue := &ProductAttributeValue{}
	bodyIsInvalid := ensureRequestBodyValidity(res, req, newProductAttributeValue)
	if bodyIsInvalid {
		return
	}
	newProductAttributeValue.ProductAttributeID = attributeID

	_, err = CreateProductAttributeValue(newProductAttributeValue)
	if err != nil {
		errorString := fmt.Sprintf("error inserting product into database: %v", err)
		log.Println(errorString)
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}
}
