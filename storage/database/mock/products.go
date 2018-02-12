package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) GetProductBySKU(db database.Querier, sku string) (*models.Product, error) {
	args := m.Called(db, sku)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockDB) ProductWithSKUExists(db database.Querier, sku string) (bool, error) {
	args := m.Called(db, sku)
	return args.Bool(0), args.Error(1)
}
func (m *MockDB) GetProductsByProductRootID(db database.Querier, productRootID uint64) ([]models.Product, error) {
	args := m.Called(db, productRootID)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockDB) ProductExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProduct(db database.Querier, id uint64) (*models.Product, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockDB) GetProductList(db database.Querier, qf *models.QueryFilter) ([]models.Product, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockDB) GetProductCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProduct(db database.Querier, nu *models.Product) (uint64, time.Time, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Get(2).(time.Time), args.Error(3)
}

func (m *MockDB) UpdateProduct(db database.Querier, updated *models.Product) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProduct(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
func (m *MockDB) ArchiveProductsWithProductRootID(db database.Querier, id uint64) (t time.Time, err error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
