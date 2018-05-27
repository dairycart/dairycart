package models

import (
	"time"
)

// Discount represents a Dairycart discount
type Discount struct {
	ID            uint64     `json:"id"`             // id
	Name          string     `json:"name"`           // name
	DiscountType  string     `json:"discount_type"`  // discount_type
	Amount        float64    `json:"amount"`         // amount
	ExpiresOn     *Dairytime `json:"expires_on"`     // expires_on
	RequiresCode  bool       `json:"requires_code"`  // requires_code
	Code          string     `json:"code"`           // code
	LimitedUse    bool       `json:"limited_use"`    // limited_use
	NumberOfUses  uint64     `json:"number_of_uses"` // number_of_uses
	LoginRequired bool       `json:"login_required"` // login_required
	StartsOn      time.Time  `json:"starts_on"`      // starts_on
	CreatedOn     time.Time  `json:"created_on"`     // created_on
	UpdatedOn     *Dairytime `json:"updated_on"`     // updated_on
	ArchivedOn    *Dairytime `json:"archived_on"`    // archived_on
}

// DiscountCreationInput is a struct to use for creating Discounts
type DiscountCreationInput struct {
	Name          string     `json:"name,omitempty"`           // name
	DiscountType  string     `json:"discount_type,omitempty"`  // discount_type
	Amount        float64    `json:"amount,omitempty"`         // amount
	ExpiresOn     *Dairytime `json:"expires_on,omitempty"`     // expires_on
	RequiresCode  bool       `json:"requires_code,omitempty"`  // requires_code
	Code          string     `json:"code,omitempty"`           // code
	LimitedUse    bool       `json:"limited_use,omitempty"`    // limited_use
	NumberOfUses  uint64     `json:"number_of_uses,omitempty"` // number_of_uses
	LoginRequired bool       `json:"login_required,omitempty"` // login_required
	StartsOn      *Dairytime `json:"starts_on,omitempty"`      // starts_on
}

// DiscountUpdateInput is a struct to use for updating Discounts
type DiscountUpdateInput struct {
	Name          string     `json:"name,omitempty"`           // name
	DiscountType  string     `json:"discount_type,omitempty"`  // discount_type
	Amount        float64    `json:"amount,omitempty"`         // amount
	ExpiresOn     *Dairytime `json:"expires_on,omitempty"`     // expires_on
	RequiresCode  bool       `json:"requires_code,omitempty"`  // requires_code
	Code          string     `json:"code,omitempty"`           // code
	LimitedUse    bool       `json:"limited_use,omitempty"`    // limited_use
	NumberOfUses  uint64     `json:"number_of_uses,omitempty"` // number_of_uses
	LoginRequired bool       `json:"login_required,omitempty"` // login_required
	StartsOn      *Dairytime `json:"starts_on,omitempty"`      // starts_on
}

type DiscountListResponse struct {
	ListResponse
	Discounts []Discount `json:"discounts"`
}
