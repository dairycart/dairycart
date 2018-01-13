package imgmock

import (
	"image"

	"github.com/dairycart/dairycart/storage/images"

	"github.com/stretchr/testify/mock"
)

type MockImageStorer struct {
	mock.Mock
}

var _ images.ImageStorer = (*MockImageStorer)(nil)

func (m *MockImageStorer) CreateThumbnails(in image.Image) images.ProductImageSet {
	args := m.Called(in)
	return args.Get(0).(images.ProductImageSet)
}

func (m *MockImageStorer) StoreImages(in images.ProductImageSet, sku string, id uint) (*images.ProductImageLocations, error) {
	args := m.Called(in, sku, id)
	return args.Get(0).(*images.ProductImageLocations), args.Error(1)
}
