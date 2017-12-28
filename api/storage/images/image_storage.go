package dairyphoto

import (
	//  "github.com/dairycart/dairymodels/v1"
	"image"
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
	CreateThumbnails(img image.Image) ProductImageSet
	StoreImages(imgset ProductImageSet, sku string, id uint) (*ProductImageLocations, error)
}
