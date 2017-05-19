package api

import (
	"log"
	"time"
)

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID            int64        `json:"id"`
	Name          string       `json:"Name"`
	BaseProductID int64        `json:"base_product_id"` // note: I don't think this name is that descriptive
	BaseProduct   *BaseProduct `json:"base_product"`
	CreatedAt     time.Time    `json:"created_at"`
	ArchivedAt    time.Time    `json:"archived_at"`
}

// productAttributeExists checks for the existence of a given ProductAttribute in the database
func productAttributeExists(db Database, id int64) bool {
	count, err := db.Model(&ProductAttribute{}).Where("id = ?", id).Where("archived_at is null").Count()
	if err != nil {
		log.Printf("error occurred querying for product_attribute: %v\n", err)
	}

	return count == 1
}
