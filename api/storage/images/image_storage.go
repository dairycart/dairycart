package image_storage

import (
	//  "github.com/dairycart/dairymodels/v1"
	"image"
)

type ImageStorer interface {
	CreateThumbnails(image.Image) []image.Image
	StoreImage(image.Image, string) error
}
