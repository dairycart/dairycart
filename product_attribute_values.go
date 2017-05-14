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
