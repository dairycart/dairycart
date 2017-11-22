package models

import (
	"time"
)

// Product represents a Dairycart product
type Product struct {
	ID                 uint64    `json:"id"`                   // id
	ProductRootID      uint64    `json:"product_root_id"`      // product_root_id
	Name               string    `json:"name"`                 // name
	Subtitle           string    `json:"subtitle"`             // subtitle
	Description        string    `json:"description"`          // description
	OptionSummary      string    `json:"option_summary"`       // option_summary
	SKU                string    `json:"sku"`                  // sku
	UPC                string    `json:"upc"`                  // upc
	Manufacturer       string    `json:"manufacturer"`         // manufacturer
	Brand              string    `json:"brand"`                // brand
	Quantity           uint32    `json:"quantity"`             // quantity
	Taxable            bool      `json:"taxable"`              // taxable
	Price              float64   `json:"price"`                // price
	OnSale             bool      `json:"on_sale"`              // on_sale
	SalePrice          float64   `json:"sale_price"`           // sale_price
	Cost               float64   `json:"cost"`                 // cost
	ProductWeight      float64   `json:"product_weight"`       // product_weight
	ProductHeight      float64   `json:"product_height"`       // product_height
	ProductWidth       float64   `json:"product_width"`        // product_width
	ProductLength      float64   `json:"product_length"`       // product_length
	PackageWeight      float64   `json:"package_weight"`       // package_weight
	PackageHeight      float64   `json:"package_height"`       // package_height
	PackageWidth       float64   `json:"package_width"`        // package_width
	PackageLength      float64   `json:"package_length"`       // package_length
	QuantityPerPackage uint32    `json:"quantity_per_package"` // quantity_per_package
	AvailableOn        time.Time `json:"available_on"`         // available_on
	CreatedOn          time.Time `json:"created_on"`           // created_on
	UpdatedOn          NullTime  `json:"updated_on"`           // updated_on
	ArchivedOn         NullTime  `json:"archived_on"`          // archived_on

	// useful for responses
	ApplicableOptionValues []ProductOptionValue `json:"applicable_options,omitempty"`
}
