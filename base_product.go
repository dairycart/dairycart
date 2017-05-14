package main

import "time"

// BaseProduct is the parent product for every product
type BaseProduct struct {
	// Basic Info
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Pricing Fields
	Taxable               bool    `json:"taxable"`
	CustomerCanSetPricing bool    `json:"customer_can_set_pricing"`
	BasePrice             float32 `json:"base_price"`

	// Product Dimensions
	BaseProductWeight float32 `json:"base_product_weight"`
	BaseProductHeight float32 `json:"base_product_height"`
	BaseProductWidth  float32 `json:"base_product_width"`
	BaseProductLength float32 `json:"base_product_length"`

	// Package dimensions
	BasePackageWeight float32 `json:"base_package_weight"`
	BasePackageHeight float32 `json:"base_package_height"`
	BasePackageWidth  float32 `json:"base_package_width"`
	BasePackageLength float32 `json:"base_package_length"`

	// Housekeeping
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created"`
	ArchivedAt time.Time `json:"-"`
}

// NewBaseProductFromProduct takes a Product object and create a BaseProduct from it
func NewBaseProductFromProduct(p *Product) *BaseProduct {
	bp := &BaseProduct{
		Name:                  p.Name,
		Description:           p.Description,
		Taxable:               p.Taxable,
		CustomerCanSetPricing: p.CustomerCanSetPricing,
		BasePrice:             p.Price,
		BaseProductWeight:     p.ProductWeight,
		BaseProductHeight:     p.ProductHeight,
		BaseProductWidth:      p.ProductWidth,
		BaseProductLength:     p.ProductLength,
		BasePackageWeight:     p.PackageWeight,
		BasePackageHeight:     p.PackageHeight,
		BasePackageWidth:      p.PackageWidth,
		BasePackageLength:     p.PackageLength,
	}

	return bp
}

// RetrieveBaseProductFromDB retrieves a product with a given SKU from the database
func RetrieveBaseProductFromDB(id int64) (*BaseProduct, error) {
	bp := &BaseProduct{}
	product := db.Model(bp).
		Where("id = ?", id).
		Where("base_product.archived_at is null")

	err := product.Select()
	return bp, err
}
