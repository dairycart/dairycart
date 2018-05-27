package models

import (
	"time"
)

// ProductRoot represents a Dairycart product root
type ProductRoot struct {
	ID                 uint64     `json:"id"`                   // id
	Name               string     `json:"name"`                 // name
	PrimaryImageID     *uint64    `json:"primary_image_id"`     // primary_image_id
	Subtitle           string     `json:"subtitle"`             // subtitle
	Description        string     `json:"description"`          // description
	SKUPrefix          string     `json:"sku_prefix"`           // sku_prefix
	Manufacturer       string     `json:"manufacturer"`         // manufacturer
	Brand              string     `json:"brand"`                // brand
	Taxable            bool       `json:"taxable"`              // taxable
	Cost               float64    `json:"cost"`                 // cost
	ProductWeight      float64    `json:"product_weight"`       // product_weight
	ProductHeight      float64    `json:"product_height"`       // product_height
	ProductWidth       float64    `json:"product_width"`        // product_width
	ProductLength      float64    `json:"product_length"`       // product_length
	PackageWeight      float64    `json:"package_weight"`       // package_weight
	PackageHeight      float64    `json:"package_height"`       // package_height
	PackageWidth       float64    `json:"package_width"`        // package_width
	PackageLength      float64    `json:"package_length"`       // package_length
	QuantityPerPackage uint32     `json:"quantity_per_package"` // quantity_per_package
	AvailableOn        time.Time  `json:"available_on"`         // available_on
	CreatedOn          time.Time  `json:"created_on"`           // created_on
	UpdatedOn          *Dairytime `json:"updated_on"`           // updated_on
	ArchivedOn         *Dairytime `json:"archived_on"`          // archived_on

	// useful for responses
	Options  []ProductOption `json:"options"`
	Images   []ProductImage  `json:"images"`
	Products []Product       `json:"products"`
}

// ProductRootCreationInput is a struct to use for creating ProductRoots
type ProductRootCreationInput struct {
	Name               string     `json:"name,omitempty"`                 // name
	Subtitle           string     `json:"subtitle,omitempty"`             // subtitle
	Description        string     `json:"description,omitempty"`          // description
	SKUPrefix          string     `json:"sku_prefix,omitempty"`           // sku_prefix
	Manufacturer       string     `json:"manufacturer,omitempty"`         // manufacturer
	Brand              string     `json:"brand,omitempty"`                // brand
	Taxable            bool       `json:"taxable,omitempty"`              // taxable
	Cost               float64    `json:"cost,omitempty"`                 // cost
	ProductWeight      float64    `json:"product_weight,omitempty"`       // product_weight
	ProductHeight      float64    `json:"product_height,omitempty"`       // product_height
	ProductWidth       float64    `json:"product_width,omitempty"`        // product_width
	ProductLength      float64    `json:"product_length,omitempty"`       // product_length
	PackageWeight      float64    `json:"package_weight,omitempty"`       // package_weight
	PackageHeight      float64    `json:"package_height,omitempty"`       // package_height
	PackageWidth       float64    `json:"package_width,omitempty"`        // package_width
	PackageLength      float64    `json:"package_length,omitempty"`       // package_length
	QuantityPerPackage uint32     `json:"quantity_per_package,omitempty"` // quantity_per_package
	AvailableOn        *Dairytime `json:"available_on,omitempty"`         // available_on
}

// ProductRootUpdateInput is a struct to use for updating ProductRoots
type ProductRootUpdateInput struct {
	Name               string     `json:"name,omitempty"`                 // name
	PrimaryImageID     *uint64    `json:"primary_image_id,omitempty"`     // primary_image_id
	Subtitle           string     `json:"subtitle,omitempty"`             // subtitle
	Description        string     `json:"description,omitempty"`          // description
	SKUPrefix          string     `json:"sku_prefix,omitempty"`           // sku_prefix
	Manufacturer       string     `json:"manufacturer,omitempty"`         // manufacturer
	Brand              string     `json:"brand,omitempty"`                // brand
	Taxable            bool       `json:"taxable,omitempty"`              // taxable
	Cost               float64    `json:"cost,omitempty"`                 // cost
	ProductWeight      float64    `json:"product_weight,omitempty"`       // product_weight
	ProductHeight      float64    `json:"product_height,omitempty"`       // product_height
	ProductWidth       float64    `json:"product_width,omitempty"`        // product_width
	ProductLength      float64    `json:"product_length,omitempty"`       // product_length
	PackageWeight      float64    `json:"package_weight,omitempty"`       // package_weight
	PackageHeight      float64    `json:"package_height,omitempty"`       // package_height
	PackageWidth       float64    `json:"package_width,omitempty"`        // package_width
	PackageLength      float64    `json:"package_length,omitempty"`       // package_length
	QuantityPerPackage uint32     `json:"quantity_per_package,omitempty"` // quantity_per_package
	AvailableOn        *Dairytime `json:"available_on,omitempty"`         // available_on
}

type ProductRootListResponse struct {
	ListResponse
	ProductRoots []ProductRoot `json:"product_roots"`
}
