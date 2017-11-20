package dairymock

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetProductBySKU(sku string) (*models.Product, error) {
	args := m.Called(sku)
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockDB) ProductWithSKUExists(sku string) (bool, error) {
	args := m.Called(sku)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ProductExists(id uint64) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProduct(id uint64) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockDB) CreateProduct(nu *models.Product) (uint64, time.Time, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Get(2).(time.Time), args.Error(3)
}

func (m *MockDB) UpdateProduct(updated *models.Product) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProduct(id uint64, tx *sql.Tx) (time.Time, error) {
	args := m.Called(id, tx)
	return args.Get(0).(time.Time), args.Error(1)
}
