package images

import (
	"image"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

type ProductImageSet struct {
	Thumbnail image.Image
	Main      image.Image
	Original  image.Image
}

type ProductImageLocations struct {
	Thumbnail string
	Main      string
	Original  string
}

type ImageStorer interface {
	Init(config *viper.Viper, router chi.Router) error
	CreateThumbnails(img image.Image) ProductImageSet
	StoreImages(imgset ProductImageSet, sku string, id uint) (*ProductImageLocations, error)
}
