package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/fatih/set"
	"github.com/pkg/errors"
	"io/ioutil"
)

func handleProductCreationImages(tx *sql.Tx, client storage.Storer, imager storage.ImageStorer, images []models.ProductImageCreationInput, sku string, rootID uint64) ([]models.ProductImage, *uint64, error) {
	var imagesToCreate []models.ProductImageCreationInput
	createdImages := set.New()
	for _, img := range images {
		if createdImages.Has(img.Data) {
			continue
		}
		createdImages.Add(img.Data)
		imagesToCreate = append(imagesToCreate, img)
	}

	// FIXME: Make this whole process concurrent
	var primaryImageID *uint64
	returnImages := []models.ProductImage{}
	for i, imageInput := range imagesToCreate {
		var (
			format string
			img image.Image
			err error
		)
		imageType := strings.ToLower(imageInput.Type)

		switch imageType {
		case "base64":
			// note: base64 expects raw base64 data, not a data URI (`data:image/png;base64,blahblahblah`)
			reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageInput.Data))
			img, format, err = image.Decode(reader)
			if format != "png" {
				return nil, nil, errors.New("only pngs are accepted")
			} else if err != nil {
				return nil, nil, fmt.Errorf("Image data at index %d is invalid", i)
			}
		case "url":
			// FIXME: this is almost definitely the wrong way to do this,
			// we should support conversion from known data types (mainly JPEGs) to PNGs
			if !strings.HasSuffix(imageInput.Data, "png") {
				return nil, nil, errors.New("only PNG images are supported")
			}
			response, err := http.Get(imageInput.Data)
			if err != nil {
				return nil, nil, errors.Wrap(err, fmt.Sprintf("error retrieving product image from url %s", imageInput.Data))
			} else {
				defer response.Body.Close()

				

				img, _, err = image.Decode(response.Body)
				if err != nil {
					return nil, nil, fmt.Errorf("Image data at index %d is invalid: %v", i, err)
				}
			}
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