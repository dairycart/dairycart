package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) ProductOptionValueForOptionIDExists(db database.Querier, optionID uint64, value string) (bool, error) {
	args := m.Called(db, optionID, value)
	return args.Bool(0), args.Error(1)
}
func (m *MockDB) ArchiveProductOptionValuesForOption(db database.Querier, optionID uint64) (time.Time, error) {
	args := m.Called(db, optionID)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) GetProductOptionValuesForOption(db database.Querier, optionID uint64) ([]models.ProductOptionValue, error) {
	args := m.Called(db, optionID)
	return args.Get(0).([]models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) ProductOptionValueExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductOptionValue(db database.Querier, id uint64) (*models.ProductOptionValue, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) GetProductOptionValueList(db database.Querier, qf *models.QueryFilter) ([]models.ProductOptionValue, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) GetProductOptionValueCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductOptionValue(db database.Querier, nu *models.ProductOptionValue) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductOptionValue(db database.Querier, updated *models.ProductOptionValue) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductOptionValue(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
func (m *MockDB) ArchiveProductOptionValuesWithProductRootID(db database.Querier, id uint64) (t time.Time, err error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
