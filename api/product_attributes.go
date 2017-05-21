package api

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/lib/pq"
)

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID         int64       `json:"id"`
	Name       string      `json:"Name"`
	ProductID  int64       `json:"product_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  pq.NullTime `json:"updated_at"`
	ArchivedAt pq.NullTime `json:"archived_at"`
}

// productAttributeExists checks for the existence of a given ProductAttribute in the database
func productAttributeExists(db *pg.DB, id int64) bool {
	count, err := db.Model(&ProductAttribute{}).Where("id = ?", id).Where("archived_at is null").Count()
	if err != nil {
		log.Printf("error occurred querying for product_attribute: %v\n", err)
	}

	return count == 1
}
