package models

import (
	"time"
)

// Discount represents a Dairycart discount
type Discount struct {
	ID            uint64    `json:"id"`             // id
	Name          string    `json:"name"`           // name
	DiscountType  string    `json:"discount_type"`  // discount_type
	Amount        float64   `json:"amount"`         // amount
	StartsOn      time.Time `json:"starts_on"`      // starts_on
	ExpiresOn     NullTime  `json:"expires_on"`     // expires_on
	RequiresCode  bool      `json:"requires_code"`  // requires_code
	Code          string    `json:"code"`           // code
	LimitedUse    bool      `json:"limited_use"`    // limited_use
	NumberOfUses  uint64    `json:"number_of_uses"` // number_of_uses
	LoginRequired bool      `json:"login_required"` // login_required
	CreatedOn     time.Time `json:"created_on"`     // created_on
	UpdatedOn     NullTime  `json:"updated_on"`     // updated_on
	ArchivedOn    NullTime  `json:"archived_on"`    // archived_on
}
