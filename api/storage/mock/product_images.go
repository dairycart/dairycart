package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) GetProductImagesByProductID(db storage.Querier, productID uint64) ([]models.ProductImage, error) {
	args := m.Called(db, productID)
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockDB) ProductImageExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductImage(db storage.Querier, id uint64) (*models.ProductImage, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductImage), args.Error(1)
}

func (m *MockDB) GetProductImageList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductImage, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockDB) GetProductImageCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductImage(db storage.Querier, nu *models.ProductImage) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductImage(db storage.Querier, updated *models.ProductImage) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductImage(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
