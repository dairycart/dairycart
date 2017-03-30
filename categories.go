package main

// Category represents a category that a product can belong to
type Category struct {
	ID int64

	Name string
}

// CategoryProducts is a table where links between categories and products are stored.
type CategoryProducts struct {
	ID         int64
	CategoryID int64
	ProductID  int64
}
