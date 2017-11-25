package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) ProductOptionValueForOptionIDExists(db storage.Querier, optionID uint64, value string) (bool, error) {
	args := m.Called(db, optionID, value)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ProductOptionValueExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductOptionValue(db storage.Querier, id uint64) (*models.ProductOptionValue, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) GetProductOptionValueList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductOptionValue, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) GetProductOptionValueCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductOptionValue(db storage.Querier, nu *models.ProductOptionValue) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductOptionValue(db storage.Querier, updated *models.ProductOptionValue) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductOptionValue(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
func (m *MockDB) ArchiveProductOptionValuesWithProductRootID(db storage.Querier, id uint64) (t time.Time, err error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
