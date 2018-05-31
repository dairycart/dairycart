package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) ProductImageBridgeExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductImageBridge(db database.Querier, id uint64) (*models.ProductImageBridge, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductImageBridge), args.Error(1)
}

func (m *MockDB) GetProductImageBridgeList(db database.Querier, qf *models.QueryFilter) ([]models.ProductImageBridge, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductImageBridge), args.Error(1)
}

func (m *MockDB) GetProductImageBridgeCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductImageBridge(db database.Querier, nu *models.ProductImageBridge) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductImageBridge(db database.Querier, updated *models.ProductImageBridge) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductImageBridge(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
