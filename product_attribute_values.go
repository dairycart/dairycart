package main

import "time"

// ProductAttributeValue represents a products variant attribute values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size, and six ProductAttributeValues,
// One for each color and one for each size.
type ProductAttributeValue struct {
	ID                 int64             `json:"id"`
	ProductAttributeID int64             `json:"product_attribute_id"`
	ProductAttribute   *ProductAttribute `json:"product_attribute"`
	Value              string            `json:"value"`
	Active             bool              `json:"active"`
	CreatedAt          time.Time         `json:"created_at"`
	ArchivedAt         time.Time         `json:"archived_at"`
}

// CreateProductAttributeValue creates a ProductAttributeValue tied to a ProductAttribute
func CreateProductAttributeValue(productAttributeID int64, pav *ProductAttributeValue) (*ProductAttributeValue, error) {
	pav.ProductAttributeID = productAttributeID
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
