package api

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
	"github.com/dairycart/dairycart/storage/v1/images"

	"github.com/fatih/set"
	"github.com/pkg/errors"
	"io"
)

func loadImage(in io.Reader) (image.Image, error) {
	img, format, err := image.Decode(in)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding image")
	} else if format != "png" {
		return nil, errors.New("only pngs are accepted")
	}
	return img, err
}

func handleProductCreationImages(tx *sql.Tx, client database.Storer, imager images.ImageStorer, images []models.ProductImageCreationInput, sku string, rootID uint64) ([]models.ProductImage, *uint64, error) {
	var (
		primaryImageID *uint64
		imagesToCreate []models.ProductImageCreationInput
		returnImages   []models.ProductImage
	)

	createdImages := set.New()
	for _, img := range images {
		if createdImages.Has(img.Data) {
			continue
		}
		createdImages.Add(img.Data)
		imagesToCreate = append(imagesToCreate, img)
	}

	// FIXME: Make this whole process concurrent
	for i, imageInput := range imagesToCreate {
		var (
			img image.Image
			err error
		)
		imageType := strings.ToLower(imageInput.Type)

		switch imageType {
		case "base64":
			// note: base64 expects raw base64 data, not a data URI (`data:image/png;base64,blahblahblah`)
			reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageInput.Data))
			img, err = loadImage(reader)
			if err != nil {
				return nil, nil, errors.Wrap(err, fmt.Sprintf("image data at index %d is invalid", i))
			}
		case "url":
			// FIXME: this is suuuuuuuuper lame, support JPEGs at least
			if !strings.HasSuffix(imageInput.Data, "png") {
				return nil, nil, errors.New("only PNG images are supported")
			}
			response, err := http.Get(imageInput.Data)
			if err != nil {
				return nil, nil, errors.Wrap(err, fmt.Sprintf("error retrieving product image from url %s", imageInput.Data))
			} else {
				img, err = loadImage(response.Body)
				if err != nil {
					return nil, nil, errors.Wrap(err, fmt.Sprintf("image data at index %d is invalid", i))
				}
			}
			response.Body.Close()
		}

		thumbnails := imager.CreateThumbnails(img)
		locations, err := imager.StoreImages(thumbnails, sku, uint(i))
		if err != nil || locations == nil {
			return nil, nil, err
		}

		newImage := &models.ProductImage{
			ProductRootID: rootID,
			ThumbnailURL:  locations.Thumbnail,
			MainURL:       locations.Main,
			OriginalURL:   locations.Original,
		}

		if imageType == "url" {
			newImage.SourceURL = imageInput.Data
		}

		newImage.ID, newImage.CreatedOn, err = client.CreateProductImage(tx, newImage)
		if err != nil {
			return nil, nil, err
		}

		if imageInput.IsPrimary && primaryImageID == nil {
			primaryImageID = &newImage.ID
		}

		returnImages = append(returnImages, *newImage)
	}
	return returnImages, primaryImageID, nil
}
