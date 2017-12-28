package local

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/dairycart/dairycart/api/storage/images"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

const LocalProductImagesDirectory = "product_images/"

type LocalImageStorer struct{}

var _ dairyphoto.ImageStorer = (*LocalImageStorer)(nil)

func (lis *LocalImageStorer) CreateThumbnails(in image.Image) dairyphoto.ProductImageSet {
	return dairyphoto.ProductImageSet{
		Thumbnail: resize.Thumbnail(100, 100, in, resize.NearestNeighbor),
		Main:      resize.Thumbnail(500, 500, in, resize.NearestNeighbor),
		Original:  in,
	}
}

func saveImage(in image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "error creating local file")
	}

	return png.Encode(f, in)
}

func (lis *LocalImageStorer) StoreImages(in dairyphoto.ProductImageSet, sku string, id uint) (*dairyphoto.ProductImageLocations, error) {
	var err error
	if _, err = os.Stat(LocalProductImagesDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(LocalProductImagesDirectory, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "error creating necessary folders")
		}
	}
	out := &dairyphoto.ProductImageLocations{}

	thumbnailPath := fmt.Sprintf("images/%s/%d/thumbnail.png", sku, id)
	err = saveImage(in.Thumbnail, thumbnailPath)
	if err != nil {
		return nil, err
	}
	out.Thumbnail = thumbnailPath

	mainPath := fmt.Sprintf("images/%s/%d/main.png", sku, id)
	err = saveImage(in.Main, mainPath)
	if err != nil {
		return out, err
	}
	out.Main = mainPath

	originalPath := fmt.Sprintf("images/%s/%d/original.png", sku, id)
	err = saveImage(in.Original, originalPath)
	if err != nil {
		return out, err
	}
	out.Original = originalPath

	return out, nil
}
