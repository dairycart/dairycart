package models

import (
	"time"
)

// Product describes something a user can buy
type Product struct {
	DBRow
	// Basic Info
	ProductRootID      uint64 `json:"product_root_id"`
	Name               string `json:"name"`
	Subtitle           string `json:"subtitle"`
	Description        string `json:"description"`
	OptionSummary      string `json:"option_summary"`
	SKU                string `json:"sku"`
	UPC                string `json:"upc"`
	Manufacturer       string `json:"manufacturer"`
	Brand              string `json:"brand"`
	Quantity           uint32 `json:"quantity"`
	QuantityPerPackage uint32 `json:"quantity_per_package"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

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

	ApplicableOptionValues []ProductOptionValue `json:"applicable_options,omitempty"`

	AvailableOn time.Time `json:"available_on"`
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	// Core Product stuff
	Name         string `json:"name"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
	SKU          string `json:"sku"`
	UPC          string `json:"upc"`
	Manufacturer string `json:"manufacturer"`
	Brand        string `json:"brand"`
	Quantity     uint32 `json:"quantity"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight      float32 `json:"package_weight"`
	PackageHeight      float32 `json:"package_height"`
	PackageWidth       float32 `json:"package_width"`
	PackageLength      float32 `json:"package_length"`
	QuantityPerPackage uint32  `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`

	// Other things
	Options []*ProductOptionCreationInput `json:"options"`
}
