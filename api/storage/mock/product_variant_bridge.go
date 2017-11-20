package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

func (m *MockDB) GetProductVariantBridge(id uint64) (models.ProductVariantBridge, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) CreateProductVariantBridge(nu models.ProductVariantBridge) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0), args.Get(1), args.Error(2)
}

func (m *MockDB) UpdateProductVariantBridge(updated models.ProductVariantBridge) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridge(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
