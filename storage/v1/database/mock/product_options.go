package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) GetProductOptionsByProductRootID(db database.Querier, productRootID uint64) ([]models.ProductOption, error) {
	args := m.Called(db, productRootID)
	return args.Get(0).([]models.ProductOption), args.Error(1)
}
func (m *MockDB) ProductOptionWithNameExistsForProductRoot(db database.Querier, name string, productRootID uint64) (bool, error) {
	args := m.Called(db, name, productRootID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ProductOptionExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductOption(db database.Querier, id uint64) (*models.ProductOption, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductOption), args.Error(1)
}

func (m *MockDB) GetProductOptionList(db database.Querier, qf *models.QueryFilter) ([]models.ProductOption, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductOption), args.Error(1)
}

func (m *MockDB) GetProductOptionCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductOption(db database.Querier, nu *models.ProductOption) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductOption(db database.Querier, updated *models.ProductOption) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductOption(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
func (m *MockDB) ArchiveProductOptionsWithProductRootID(db database.Querier, id uint64) (t time.Time, err error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
