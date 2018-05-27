package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) GetDiscountByCode(db database.Querier, code string) (*models.Discount, error) {
	args := m.Called(db, code)
	return args.Get(0).(*models.Discount), args.Error(1)
}

func (m *MockDB) DiscountExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetDiscount(db database.Querier, id uint64) (*models.Discount, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.Discount), args.Error(1)
}

func (m *MockDB) GetDiscountList(db database.Querier, qf *models.QueryFilter) ([]models.Discount, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.Discount), args.Error(1)
}

func (m *MockDB) GetDiscountCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateDiscount(db database.Querier, nu *models.Discount) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateDiscount(db database.Querier, updated *models.Discount) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteDiscount(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
