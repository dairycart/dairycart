package api

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	productAttributeExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_attributes WHERE id = $1 and archived_at is null);`
	productAttributeCreationQuery  = `INSERT INTO product_attributes ("name", "product_progenitor_id") VALUES ($1, $2);`
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
func productAttributeExists(db *sql.DB, id int64) (bool, error) {
	var exists string

	err := db.QueryRow(productAttributeExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, errors.Wrap(err, "Error querying for product")
	}

	return exists == "true", err
}
