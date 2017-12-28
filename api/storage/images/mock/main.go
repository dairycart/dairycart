package imgmock

import (
	"image"

	"github.com/dairycart/dairycart/api/storage/images"

	"github.com/stretchr/testify/mock"
)

type MockImageStorer struct {
	mock.Mock
}

var _ dairyphoto.ImageStorer = (*MockImageStorer)(nil)

func (m *MockImageStorer) CreateThumbnails(in image.Image) dairyphoto.ProductImageSet {
	args := m.Called(in)
	return args.Get(0).(dairyphoto.ProductImageSet)
}

func (m *MockImageStorer) StoreImages(in dairyphoto.ProductImageSet, sku string) (*dairyphoto.ProductImageLocations, error) {
	args := m.Called(in, sku)
	return args.Get(0).(*dairyphoto.ProductImageLocations), args.Error(0)
}
