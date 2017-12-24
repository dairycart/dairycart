package local

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

type LocalImageStorer struct{}

func (lis *LocalImageStorer) CreateThumbnails(in image.Image) []image.Image {
	out := []image.Image{}
	sizes := []uint{100, 500}
	for _, size := range sizes {
		out = append(out, resize.Thumbnail(size, size, in, resize.NearestNeighbor))
		// newFilename := fmt.Sprintf("%d_%d_x_%d.png", timestamp, size, size)
	}
	return []image.Image{}
}

func (lis *LocalImageStorer) StoreImage(in image.Image, filename string) error {
	path := fmt.Sprintf("images/%s/", sku)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "error creating necessary folders")
		}
	}

	path = fmt.Sprintf("images/%s/%s", sku, filename)
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "error creating local file")
	}

	err = png.Encode(f, in)
	if err != nil {
		return errors.Wrap(err, "error encoding png")
	}

	return nil
}
