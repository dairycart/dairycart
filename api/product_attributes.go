package api

import (
	"time"

	"github.com/lib/pq"
)

const (
	productAttributeCreationQuery = `INSERT INTO product_attributes ("name", "product_progenitor_id") VALUES ($1, $2) RETURNING *;`
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
