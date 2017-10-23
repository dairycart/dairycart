package storage

import (
	"github.com/dairycart/dairycart/api/storage/models"
)

type Storage interface {
	GetProductBySKU(sku string) (models.Product, error)
	GetProductByID(id uint64) (models.Product, error)
}
