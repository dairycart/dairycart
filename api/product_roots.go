package main

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	productRootSkuExistenceQuery = `SELECT EXISTS(SELECT 1 FROM product_roots WHERE sku_prefix = $1 AND archived_on IS NULL)`
	productRootExistenceQuery    = `SELECT EXISTS(SELECT 1 FROM product_roots WHERE id = $1 AND archived_on IS NULL)`
	productRootRetrievalQuery    = `SELECT * FROM product_roots WHERE id = $1`
)

// ProductRoot represents the object that products inherit from
type ProductRoot struct {
	DBRow

	// Basic Info
	Name               string     `json:"name"`
	Subtitle           NullString `json:"subtitle"`
	Description        string     `json:"description"`
	SKUPrefix          string     `json:"sku_prefix"`
	Manufacturer       NullString `json:"manufacturer"`
	Brand              NullString `json:"brand"`
	AvailableOn        time.Time  `json:"available_on"`
	QuantityPerPackage uint32     `json:"quantity_per_package"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Cost    float32 `json:"cost"`

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

	Options  []ProductOption `json:"options"`
	Products []Product       `json:"products"`
}

func createProductRootFromProduct(p *Product) *ProductRoot {
	r := &ProductRoot{
		Name:               p.Name,
		Subtitle:           p.Subtitle,
		Description:        p.Description,
		SKUPrefix:          p.SKU,
		Manufacturer:       p.Manufacturer,
		Brand:              p.Brand,
		QuantityPerPackage: p.QuantityPerPackage,
		Taxable:            p.Taxable,
		Cost:               p.Cost,
		ProductWeight:      p.ProductWeight,
		ProductHeight:      p.ProductHeight,
		ProductWidth:       p.ProductWidth,
		ProductLength:      p.ProductLength,
		PackageWeight:      p.PackageWeight,
		PackageHeight:      p.PackageHeight,
		PackageWidth:       p.PackageWidth,
		PackageLength:      p.PackageLength,
		AvailableOn:        p.AvailableOn,
	}
	return r
}

func createProductRootInDB(tx *sql.Tx, r *ProductRoot) (uint64, time.Time, error) {
	var newRootID uint64
	var createdOn time.Time
	// using QueryRow instead of Exec because we want it to return the newly created row's ID
	// Exec normally returns a sql.Result, which has a LastInsertedID() method, but when I tested
	// this locally, it never worked. ¯\_(ツ)_/¯
	query, queryArgs := buildProductRootCreationQuery(r)
	err := tx.QueryRow(query, queryArgs...).Scan(&newRootID, &createdOn)

	return newRootID, createdOn, err
}

// retrieveProductRootFromDB retrieves a product root with a given ID from the database
func retrieveProductRootFromDB(db *sqlx.DB, id uint64) (*ProductRoot, error) {
	root := &ProductRoot{}
	err := db.QueryRowx(productRootRetrievalQuery, id).StructScan(root)
	return root, err
}
