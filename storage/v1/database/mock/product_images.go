package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) GetProductImagesByProductID(db database.Querier, productID uint64) ([]models.ProductImage, error) {
	args := m.Called(db, productID)
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockDB) SetPrimaryProductImageForProduct(db database.Querier, productID, imageID uint64) (time.Time, error) {
	args := m.Called(db, productID, imageID)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) ProductImageExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetProductImage(db database.Querier, id uint64) (*models.ProductImage, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.ProductImage), args.Error(1)
}

func (m *MockDB) GetProductImageList(db database.Querier, qf *models.QueryFilter) ([]models.ProductImage, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockDB) GetProductImageCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateProductImage(db database.Querier, nu *models.ProductImage) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductImage(db database.Querier, updated *models.ProductImage) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductImage(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
