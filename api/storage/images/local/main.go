package localimagestorage

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/dairycart/dairycart/api/storage/images"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

const LocalProductImagesDirectory = "product_images"

type LocalImageStorer struct {
	BaseURL string
}

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
	baseURL := lis.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:4321"
	}

	photoDir := fmt.Sprintf("%s/%s/%d", LocalProductImagesDirectory, sku, id)

	var err error
	if _, err = os.Stat(photoDir); os.IsNotExist(err) {
		err = os.MkdirAll(photoDir, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "error creating necessary folders")
		}
	}
	out := &dairyphoto.ProductImageLocations{}

	thumbnailPath := fmt.Sprintf("%s/thumbnail.png", photoDir)
	err = saveImage(in.Thumbnail, thumbnailPath)
	if err != nil {
		return nil, err
	}
	out.Thumbnail = fmt.Sprintf("%s/%s", baseURL, thumbnailPath)

	mainPath := fmt.Sprintf("%s/main.png", photoDir)
	err = saveImage(in.Main, mainPath)
	if err != nil {
		return out, err
	}
	out.Main = fmt.Sprintf("%s/%s", baseURL, mainPath)

	originalPath := fmt.Sprintf("%s/original.png", photoDir)
	err = saveImage(in.Original, originalPath)
	if err != nil {
		return out, err
	}
	out.Original = fmt.Sprintf("%s/%s", baseURL, originalPath)

	return out, nil
}
