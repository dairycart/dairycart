package imgmock

import (
	"image"

	"github.com/dairycart/dairycart/api/storage"

	"github.com/stretchr/testify/mock"
)

type MockImageStorer struct {
	mock.Mock
}

var _ storage.ImageStorer = (*MockImageStorer)(nil)

func (m *MockImageStorer) CreateThumbnails(in image.Image) storage.ProductImageSet {
	args := m.Called(in)
	return args.Get(0).(storage.ProductImageSet)
}

func (m *MockImageStorer) StoreImages(in storage.ProductImageSet, sku string, id uint) (*storage.ProductImageLocations, error) {
	args := m.Called(in, sku)
	return args.Get(0).(*storage.ProductImageLocations), args.Error(0)
}
