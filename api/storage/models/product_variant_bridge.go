package models

import (
	"time"
)

// ProductVariantBridge represents a Dairycart productvariantbridge
type ProductVariantBridge struct {
	ID                   uint64    `json:"id"`                      // id
	ProductID            uint64    `json:"product_id"`              // product_id
	ProductOptionValueID uint64    `json:"product_option_value_id"` // product_option_value_id
	CreatedOn            time.Time `json:"created_on"`              // created_on
	ArchivedOn           NullTime  `json:"archived_on"`             // archived_on
}
