package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) ProductVariantBridgeExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductVariantBridge(db database.Querier, id uint64) (*models.ProductVariantBridge, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductVariantBridge), args.Error(1)
}

func (m *MockDB) GetProductVariantBridgeList(db database.Querier, qf *models.QueryFilter) ([]models.ProductVariantBridge, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductVariantBridge), args.Error(1)
}

func (m *MockDB) GetProductVariantBridgeCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductVariantBridge(db database.Querier, nu *models.ProductVariantBridge) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) CreateMultipleProductVariantBridgesForProductID(db database.Querier, productID uint64, optionValueIDs []uint64) error {
	args := m.Called(db, productID, optionValueIDs)
	return args.Error(0)
}

func (m *MockDB) UpdateProductVariantBridge(db database.Querier, updated *models.ProductVariantBridge) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridge(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
func (m *MockDB) ArchiveProductVariantBridgesWithProductRootID(db database.Querier, id uint64) (t time.Time, err error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductVariantBridgeByProductID(db database.Querier, productID uint64) (t time.Time, err error) {
	args := m.Called(db, productID)
	return args.Get(0).(time.Time), args.Error(1)
}
