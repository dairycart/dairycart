package models

import (
	"time"
)

// ProductOptionValue represents a Dairycart productoptionvalue
type ProductOptionValue struct {
	ID              uint64    `json:"id"`                // id
	ProductOptionID uint64    `json:"product_option_id"` // product_option_id
	Value           string    `json:"value"`             // value
	CreatedOn       time.Time `json:"created_on"`        // created_on
	UpdatedOn       NullTime  `json:"updated_on"`        // updated_on
	ArchivedOn      NullTime  `json:"archived_on"`       // archived_on
}
