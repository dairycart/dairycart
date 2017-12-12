package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) ProductRootWithSKUPrefixExists(db storage.Querier, skuPrefix string) (bool, error) {
	args := m.Called(db, skuPrefix)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ProductRootExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductRoot(db storage.Querier, id uint64) (*models.ProductRoot, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductRoot), args.Error(1)
}

func (m *MockDB) GetProductRootList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductRoot, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductRoot), args.Error(1)
}

func (m *MockDB) GetProductRootCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductRoot(db storage.Querier, nu *models.ProductRoot) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductRoot(db storage.Querier, updated *models.ProductRoot) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductRoot(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
