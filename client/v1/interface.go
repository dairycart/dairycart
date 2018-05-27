package dairyclient

import (
	"github.com/dairycart/dairymodels/v1"
)

type DairyclientV1 interface {
	ProductExists(sku string) (bool, error)
	GetProduct(sku string) (*models.Product, error)
	GetProducts(queryFilter map[string]string) ([]models.Product, error)
	CreateProduct(np models.ProductCreationInput) (*models.Product, error)
	UpdateProduct(sku string, up models.ProductCreationInput) (*models.Product, error)
	DeleteProduct(sku string) error
	GetProductRoot(rootID uint64) (*models.ProductRoot, error)
	GetProductRoots(queryFilter map[string]string) ([]models.ProductRoot, error)
	DeleteProductRoot(rootID uint64) error
	GetProductOptions(productID uint64, queryFilter map[string]string) ([]models.ProductOption, error)
	CreateProductOptionForProduct(productID uint64, no models.ProductOption) (*models.ProductOption, error)
	UpdateProductOption(optionID uint64, uo models.ProductOption) (*models.ProductOption, error)
	DeleteProductOption(optionID uint64) error
	CreateProductOptionValueForOption(optionID uint64, nv models.ProductOptionValue) (*models.ProductOptionValue, error)
	UpdateProductOptionValueForOption(valueID uint64, uv models.ProductOptionValue) (*models.ProductOptionValue, error)
	DeleteProductOptionValueForOption(optionID uint64) error
}
