package imgmock

import (
	"image"

	"github.com/stretchr/testify/mock"
)

type MockImageStorer struct {
	mock.Mock
}

func (m *MockImageStorer) CreateThumbnails(in image.Image) []image.Image {
	args := m.Called(in)
	return args.Get(0).([]image.Image)
}

func (m *MockImageStorer) StoreImage(in image.Image, filename string) error {
	args := m.Called(in, filename)
	return args.Error(0)
}
