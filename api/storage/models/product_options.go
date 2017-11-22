package models

import (
	"time"
)

// ProductOption represents a Dairycart productoption
type ProductOption struct {
	ID            uint64    `json:"id"`              // id
	Name          string    `json:"name"`            // name
	ProductRootID uint64    `json:"product_root_id"` // product_root_id
	CreatedOn     time.Time `json:"created_on"`      // created_on
	UpdatedOn     NullTime  `json:"updated_on"`      // updated_on
	ArchivedOn    NullTime  `json:"archived_on"`     // archived_on

	// useful for responses
	Values []ProductOptionValue `json:"values"`
}
