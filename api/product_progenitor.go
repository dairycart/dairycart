package api

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const (
	productProgenitorExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_progenitors WHERE id = $1 and archived_at is null);`
	productProgenitorQuery          = `SELECT * FROM product_progenitors WHERE id = $1 and archived_at is null;`
	productProgenitorCreationQuery  = `
		INSERT INTO product_progenitors (
			"name",
			"description",
			"taxable",
			"price",
			"product_weight",
			"product_height",
			"product_width",
			"product_length",
			"package_weight",
			"package_height",
			"package_width",
			"package_length",
			"created_at"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW()) RETURNING id;
	`
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
	UpdatedAt  pq.NullTime `json:"-"`
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

func newProductProgenitorFromProductCreationInput(in *ProductCreationInput) *ProductProgenitor {
	return &ProductProgenitor{
		Name:          in.Name,
		Description:   in.Description,
		Taxable:       in.Taxable,
		Price:         in.Price,
		ProductWeight: in.ProductWeight,
		ProductHeight: in.ProductHeight,
		ProductWidth:  in.ProductWidth,
		ProductLength: in.ProductLength,
		PackageWeight: in.PackageWeight,
		PackageHeight: in.PackageHeight,
		PackageWidth:  in.PackageWidth,
		PackageLength: in.PackageLength,
	}
}

func createProductProgenitorInDB(db *sql.DB, g *ProductProgenitor) (*ProductProgenitor, error) {
	var newProgenitorID int64
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	err := db.QueryRow(productProgenitorCreationQuery,
		g.Name,
		g.Description,
		g.Taxable,
		g.Price,
		g.ProductWeight,
		g.ProductHeight,
		g.ProductWidth,
		g.ProductLength,
		g.PackageWeight,
		g.PackageHeight,
		g.PackageWidth,
		g.PackageLength,
	).Scan(&newProgenitorID)

	g.ID = newProgenitorID
	return g, err
}

// retrieveProductProgenitorFromDB retrieves a product progenitor with a given ID from the database
func retrieveProductProgenitorFromDB(db *sql.DB, id int64) (*ProductProgenitor, error) {
	progenitor := &ProductProgenitor{}
	scanArgs := progenitor.generateScanArgs()

	err := db.QueryRow(productProgenitorQuery, id).Scan(scanArgs...)

	return progenitor, err
}
