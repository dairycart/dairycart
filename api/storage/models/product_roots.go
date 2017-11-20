package models

import (
	"time"
)

// ProductRoot represents a Dairycart productroot
type ProductRoot struct {
	ID                 uint64    `json:"id"`                   // id
	Name               string    `json:"name"`                 // name
	Subtitle           string    `json:"subtitle"`             // subtitle
	Description        string    `json:"description"`          // description
	SkuPrefix          string    `json:"sku_prefix"`           // sku_prefix
	Manufacturer       string    `json:"manufacturer"`         // manufacturer
	Brand              string    `json:"brand"`                // brand
	Taxable            bool      `json:"taxable"`              // taxable
	Cost               float64   `json:"cost"`                 // cost
	ProductWeight      float64   `json:"product_weight"`       // product_weight
	ProductHeight      float64   `json:"product_height"`       // product_height
	ProductWidth       float64   `json:"product_width"`        // product_width
	ProductLength      float64   `json:"product_length"`       // product_length
	PackageWeight      float64   `json:"package_weight"`       // package_weight
	PackageHeight      float64   `json:"package_height"`       // package_height
	PackageWidth       float64   `json:"package_width"`        // package_width
	PackageLength      float64   `json:"package_length"`       // package_length
	QuantityPerPackage int       `json:"quantity_per_package"` // quantity_per_package
	AvailableOn        time.Time `json:"available_on"`         // available_on
	CreatedOn          time.Time `json:"created_on"`           // created_on
	UpdatedOn          NullTime  `json:"updated_on"`           // updated_on
	ArchivedOn         NullTime  `json:"archived_on"`          // archived_on
}
