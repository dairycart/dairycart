package dairymock

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetDiscountByCode(code string) (*models.Discount, error) {
	args := m.Called(code)
	return args.Get(0).(*models.Discount), args.Error(1)
}

func (m *MockDB) DiscountExists(id uint64) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetDiscount(id uint64) (*models.Discount, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Discount), args.Error(1)
}

func (m *MockDB) CreateDiscount(nu *models.Discount) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateDiscount(updated *models.Discount) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteDiscount(id uint64, tx *sql.Tx) (time.Time, error) {
	args := m.Called(id, tx)
	return args.Get(0).(time.Time), args.Error(1)
}
