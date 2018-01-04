package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	smallGreenPNG = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKAQMAAAC3/F3+AAAABlBMVEUA/wAA/wD8J4MxAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAAC0lEQVQImWNgwAcAAB4AAe72cCEAAAAASUVORK5CYII="
)

func TestHandleProductCreationImages(t *testing.T) {
	t.Parallel()

	exampleSKU := "example"
	exampleID := uint64(1)
	exampleThumbnailLocation := "https://dairycart.com/product_images/sku/0/thumbnail.png"
	exampleMainLocation := "https://dairycart.com/product_images/sku/0/main.png"
	exampleOriginalLocation := "https://dairycart.com/product_images/sku/0/original.png"

	t.Run("optimal conditions", func(_t *testing.T) {
		_t.Parallel()
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		handlers := map[string]http.HandlerFunc{
			"/cool.png": func(res http.ResponseWriter, req *http.Request) {
				reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(smallGreenPNG))
				img, _, err := image.Decode(reader)
				require.Nil(t, err)

				buffer := new(bytes.Buffer)
				err = png.Encode(buffer, img)
				require.Nil(t, err)
				res.Write(buffer.Bytes())
			},
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "base64",
				Data: smallGreenPNG,
			},
			{
				Type:      "url",
				IsPrimary: true,
				Data:      fmt.Sprintf("%s/cool.png", ts.URL),
			},
		}

		expectedPrimaryImageID := &exampleID
		expectedImages := []models.ProductImage{
			{
				ID:            exampleID,
				ProductRootID: exampleID,
				ThumbnailURL:  exampleThumbnailLocation,
				MainURL:       exampleMainLocation,
				OriginalURL:   exampleOriginalLocation,
				CreatedOn:     buildTestTime(),
			},
			{
				ID:            exampleID,
				ProductRootID: exampleID,
				ThumbnailURL:  exampleThumbnailLocation,
				MainURL:       exampleMainLocation,
				OriginalURL:   exampleOriginalLocation,
				SourceURL:     exampleImageInputs[1].Data,
				CreatedOn:     buildTestTime(),
			},
		}

		exampleProductImageLocations := &storage.ProductImageLocations{
			Thumbnail: exampleThumbnailLocation,
			Main:      exampleMainLocation,
			Original:  exampleOriginalLocation,
		}
		arbitraryImageSet := storage.ProductImageSet{}
		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(arbitraryImageSet)
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, exampleSKU, mock.AnythingOfType("uint")).
			Return(exampleProductImageLocations, nil)

		testUtil.MockDB.On("CreateProductImage", mock.AnythingOfType("*sql.Tx"), mock.Anything).
			Return(uint64(1), buildTestTime(), nil)

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actualImages, actualPrimaryImageID, err := handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expectedImages, actualImages, "expected and actual images should match")
		assert.Equal(t, expectedPrimaryImageID, actualPrimaryImageID, "expected and actual primary image IDs should match")
	})

	t.Run("with duplicates in creation input", func(_t *testing.T) {
		_t.Parallel()
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "base64",
				Data: smallGreenPNG,
			},
			{
				Type: "base64",
				Data: smallGreenPNG,
			},
		}

		expectedImages := []models.ProductImage{
			{
				ID:            exampleID,
				ProductRootID: exampleID,
				ThumbnailURL:  exampleThumbnailLocation,
				MainURL:       exampleMainLocation,
				OriginalURL:   exampleOriginalLocation,
				CreatedOn:     buildTestTime(),
			},
		}

		exampleProductImageLocations := &storage.ProductImageLocations{
			Thumbnail: exampleThumbnailLocation,
			Main:      exampleMainLocation,
			Original:  exampleOriginalLocation,
		}
		arbitraryImageSet := storage.ProductImageSet{}
		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(arbitraryImageSet).
			Once()
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, exampleSKU, mock.AnythingOfType("uint")).
			Return(exampleProductImageLocations, nil).
			Once()

		testUtil.MockDB.On("CreateProductImage", mock.AnythingOfType("*sql.Tx"), mock.Anything).
			Return(uint64(1), buildTestTime(), nil).
			Once()

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actualImages, _, err := handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expectedImages, actualImages, "expected and actual images should match")
	})

	// FIXME: this isn't working the way it should be
	t.Run("with non png value in base64", func(_t *testing.T) {
		_t.Parallel()

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "base64",
				Data: "/9j/4AAQSkZJRgABAQEAYABgAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4wICh1c2luZyBJSkcgSlBFRyB2ODApLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAIBgYHBgUIBwcHCQkICgwUDQwLCwwZEhMPFB0aHx4dGhwcICQuJyAiLCMcHCg3KSwwMTQ0NB8nOT04MjwuMzQy/9sAQwEJCQkMCwwYDQ0YMiEcITIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIy/8AAEQgACgAKAwEiAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A1qKKK/ND8gP/2Q==",
			},
		}

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		_, _, err = handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleID)
		assert.Error(t, err)
	})

	t.Run("with erroneous base64", func(_t *testing.T) {
		_t.Parallel()

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "base64",
				Data: "lol there's no way this is a valid base64 image thing",
			},
		}

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		_, _, err = handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleID)
		assert.Error(t, err)
	})

	t.Run("rejects urls with invalid extensions", func(_t *testing.T) {
		_t.Parallel()

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "url",
				Data: "http://somesite.com/lolrememberbitmaps.bmp",
			},
		}

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		_, _, err = handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleID)
		assert.Error(t, err)
	})
}
