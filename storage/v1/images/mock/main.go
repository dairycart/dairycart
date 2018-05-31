package imgmock

import (
	"image"

	"github.com/dairycart/dairycart/storage/v1/images"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

type MockImageStorer struct {
	mock.Mock
}

var _ images.ImageStorer = (*MockImageStorer)(nil)

func (m *MockImageStorer) Init(config *viper.Viper, router chi.Router) error {
	args := m.Called(config, router)

	return args.Error(0)
}

func (m *MockImageStorer) CreateThumbnails(in image.Image) images.ProductImageSet {
	args := m.Called(in)
	return args.Get(0).(images.ProductImageSet)
}

func (m *MockImageStorer) StoreImages(in images.ProductImageSet, sku string, id uint) (*images.ProductImageLocations, error) {
	args := m.Called(in, sku, id)
	return args.Get(0).(*images.ProductImageLocations), args.Error(1)
}
