package api

import (
	"time"

	"github.com/lib/pq"
)

// Discount represents pricing changes that apply temporarily to products
type Discount struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	ProductID int64     `json:"product_id"`
	Amount    float32   `json:"amount"`
	StartsOn  time.Time `json:"starts_on"`
	ExpiresOn time.Time `json:"expires_on"`

	// Housekeeping
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  pq.NullTime `json:"-"`
	ArchivedAt pq.NullTime `json:"-"`
}

// generateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (d *Discount) generateScanArgs() []interface{} {
	return []interface{}{
		&d.ID,
		&d.ProductID,
		&d.Amount,
		&d.StartsOn,
		&d.ExpiresOn,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.ArchivedAt,
	}
}

func (d *Discount) discountTypeIsValid() bool {
	return d.Type == "percentage" || d.Type == "flat_amount"
}
