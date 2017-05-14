package main

import "time"

// ProductAttribute represents a products variant attributes. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductAttributes for that base_product, color and size.
type ProductAttribute struct {
	ID            int64        `json:"id"`
	Name          string       `json:"Name"`
	BaseProductID int64        `json:"base_product_id"` // note: I don't think this name is that descriptive
	BaseProduct   *BaseProduct `json:"base_product"`
	Active        bool         `json:"active"`
	CreatedAt     time.Time    `json:"created_at"`
	ArchivedAt    time.Time    `json:"archived_at"`
}
