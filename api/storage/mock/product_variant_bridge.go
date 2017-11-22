package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) ProductVariantBridgeExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductVariantBridge(db storage.Querier, id uint64) (*models.ProductVariantBridge, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductVariantBridge), args.Error(1)
}

func (m *MockDB) CreateProductVariantBridge(db storage.Querier, nu *models.ProductVariantBridge) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductVariantBridge(db storage.Querier, updated *models.ProductVariantBridge) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridge(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridgeByProductID(db storage.Querier, productID uint64) (t time.Time, err error) {
	args := m.Called(db, productID)
	return args.Get(0).(time.Time), args.Error(1)
}
