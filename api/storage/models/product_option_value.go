package models

// ProductOptionValue represents a product's option values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size, and six ProductOptionValues,
// One for each color and one for each size.
type ProductOptionValue struct {
	DBRow
	ProductOptionID uint64 `json:"product_option_id"`
	Value           string `json:"value"`
}
