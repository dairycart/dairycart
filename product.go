package main

// Variant describes children of products with different attributes from the parent
type Variant struct {
	ID int64

	SKU   string
	Type  string
	Value string

	HasSpecialPrice bool
	Price           float32
}

// Product describes...well, a product
type Product struct {
	ID                    int64   `json:"id"`
	SKU                   string  `json:"sku"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	UPC                   string  `json:"upc"`
	OnSale                bool    `json:"on_sale"`
	Taxable               bool    `json:"taxable"`
	CustomerCanSetPricing bool    `json:"customer_can_set_pricing"`
	Price                 float32 `json:"price"`
	Weight                float32 `json:"weight"`
	Height                float32 `json:"height"`
	Width                 float32 `json:"width"`
	Length                float32 `json:"length"`
	Quantity              int32   `json:"quantity"`

	// SalePrice float32
}
