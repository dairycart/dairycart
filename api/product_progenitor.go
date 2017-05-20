package api

import (
	"database/sql"
	"time"
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
	CreatedAt  time.Time `json:"created"`
	ArchivedAt NullTime  `json:"-"`
}

// GenerateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (g *ProductProgenitor) GenerateScanArgs() []interface{} {
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
		&g.ArchivedAt,
	}
}

// retrieveProductProgenitorFromDB retrieves a product with a given SKU from the database
func retrieveProductProgenitorFromDB(db *sql.DB, id int64) (ProductProgenitor, error) {
	var progenitor ProductProgenitor
	scanArgs := progenitor.GenerateScanArgs()

	err := db.QueryRow("SELECT * FROM product_progenitors WHERE id = $1;", id).Scan(scanArgs...)

	return progenitor, err
}
