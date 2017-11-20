package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetProductVariantBridge(id uint64) (*models.ProductVariantBridge, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ProductVariantBridge), args.Error(1)
}

func (m *MockDB) CreateProductVariantBridge(nu *models.ProductVariantBridge) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductVariantBridge(updated *models.ProductVariantBridge) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridge(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0).(time.Time), args.Error(1)
}
