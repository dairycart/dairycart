package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) ProductRootWithSKUPrefixExists(db database.Querier, skuPrefix string) (bool, error) {
	args := m.Called(db, skuPrefix)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ProductRootExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductRoot(db database.Querier, id uint64) (*models.ProductRoot, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductRoot), args.Error(1)
}

func (m *MockDB) GetProductRootList(db database.Querier, qf *models.QueryFilter) ([]models.ProductRoot, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductRoot), args.Error(1)
}

func (m *MockDB) GetProductRootCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductRoot(db database.Querier, nu *models.ProductRoot) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductRoot(db database.Querier, updated *models.ProductRoot) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductRoot(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
