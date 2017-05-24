package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	productProgenitorExistenceQuery = `SELECT * FROM product_progenitors WHERE id = $1 and archived_at is null;`
	productProgenitorQuery          = `SELECT * FROM product_progenitors WHERE id = $1 and archived_at is null;`
)

// ProductProgenitor is the parent product for every product
type ProductProgenitor struct {
	// Basic Info
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Price   float32 `json:"price"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`

	// // Housekeeping
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  pq.NullTime `json:"updated_at"`
	ArchivedAt pq.NullTime `json:"-"`
}

// generateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (g *ProductProgenitor) generateScanArgs() []interface{} {
	return []interface{}{
		&g.ID,
		&g.Name,
		&g.Description,
		&g.Taxable,
		&g.Price,
		&g.ProductWeight,
		&g.ProductHeight,
		&g.ProductWidth,
		&g.ProductLength,
		&g.PackageWeight,
		&g.PackageHeight,
		&g.PackageWidth,
		&g.PackageLength,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.ArchivedAt,
	}
}

// productProgenitorExistsInDB will return whether or not a product/attribute/etc with a given identifier exists in the database
func productProgenitorExistsInDB(db *sql.DB, id int64) (bool, error) {
	var exists string

	err := db.QueryRow(productProgenitorExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, errors.Wrap(err, "Error querying for product progenitor")
	}

	return exists == "true", err
}

func respondThatProductProgenitorDoesNotExist(req *http.Request, res http.ResponseWriter, id string) {
	log.Printf(`informing user that the product they were looking for (id %s) does not exist`, id)
	http.NotFound(res, req)
}

// retrieveProductProgenitorFromDB retrieves a product progenitor with a given ID from the database
func retrieveProductProgenitorFromDB(db *sql.DB, id int64) (ProductProgenitor, error) {
	var progenitor ProductProgenitor
	scanArgs := progenitor.generateScanArgs()

	err := db.QueryRow(productProgenitorQuery, id).Scan(scanArgs...)

	return progenitor, err
}
