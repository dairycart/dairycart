package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

func (m *MockDB) GetProductBySKU(sku string) (models.Product, error) {
	args := m.Called(sku)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) GetProduct(id uint64) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) CreateProduct(nu models.Product) (uint64, time.Time, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0), args.Get(1), args.Get(2), args.Error(3)
}

func (m *MockDB) UpdateProduct(updated models.Product) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) DeleteProduct(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
