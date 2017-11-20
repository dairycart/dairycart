package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetProductRoot(id uint64) (*models.ProductRoot, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ProductRoot), args.Error(1)
}

func (m *MockDB) CreateProductRoot(nu *models.ProductRoot) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductRoot(updated *models.ProductRoot) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductRoot(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0).(time.Time), args.Error(1)
}
