package models

import (
	"time"
)

// ProductOptionValue represents a Dairycart product option value
type ProductOptionValue struct {
	ID              uint64     `json:"id"`                // id
	ProductOptionID uint64     `json:"product_option_id"` // product_option_id
	Value           string     `json:"value"`             // value
	CreatedOn       time.Time  `json:"created_on"`        // created_on
	UpdatedOn       *Dairytime `json:"updated_on"`        // updated_on
	ArchivedOn      *Dairytime `json:"archived_on"`       // archived_on
}

// ProductOptionValueCreationInput is a struct to use for creating ProductOptionValues
type ProductOptionValueCreationInput struct {
	Value string `json:"value,omitempty"` // value
}

// ProductOptionValueUpdateInput is a struct to use for updating ProductOptionValues
type ProductOptionValueUpdateInput struct {
	ProductOptionID uint64 `json:"product_option_id,omitempty"` // product_option_id
	Value           string `json:"value,omitempty"`             // value
}

type ProductOptionValueListResponse struct {
	ListResponse
	ProductOptionValues []ProductOptionValue `json:"product_option_values"`
}
