package main

// Variation describes children of products with different attributes from the parent
type Variation struct {
	ID int64

	SKU             string
	Type            string
	Value           string
	HasSpecialPrice bool
	Price           float32
}

// Category represents a category that a product can belong to
type Category struct {
	ID int64

	Name string
}

// Product describes...well, a product
type Product struct {
	ID int64

	SKU         string
	Title       string
	Description string
	Categories  []Category

	OnSale    bool
	Price     float32
	SalePrice float32
	Taxable   bool

	Weight float32
	Height float32
	Width  float32
	Length float32

	Quantity   int32
	UPC        string
	Variations []Variation
	Images     []string
}
