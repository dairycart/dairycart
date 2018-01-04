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

func buildExampleImageResponse(t *testing.T) string {
	t.Helper()
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(smallGreenPNG))
	img, _, err := image.Decode(reader)
	require.Nil(t, err)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	require.Nil(t, err)

	return buf.String()
}

func TestHandleProductCreationImages(t *testing.T) {
	t.Parallel()

	exampleSKU := "example"
	exampleRootID := uint64(1)
	exampleThumbnailLocation := "https://dairycart.com/product_images/sku/0/thumbnail.png"
	exampleMainLocation := "https://dairycart.com/product_images/sku/0/main.png"
	exampleOriginalLocation := "https://dairycart.com/product_images/sku/0/original.png"

	t.Run("optimal conditions", func(_t *testing.T) {
		_t.Parallel()
		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()

		handlers := map[string]http.HandlerFunc{
			"cool.png": generateHandler(t, "", buildExampleImageResponse(t), http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleImageInputs := []models.ProductImageCreationInput{
			{
				Type: "base64",
				Data: smallGreenPNG,
			},
			{
				Type: "url",
				Data: fmt.Sprintf("%s/cool.png", ts.URL),
			},
		}

		expectedPrimaryImageID := uint64(123)
		expectedImages := []models.ProductImage{
			{
				ID:            exampleRootID,
				ProductRootID: exampleRootID,
				ThumbnailURL:  exampleThumbnailLocation,
				MainURL:       exampleMainLocation,
				OriginalURL:   exampleOriginalLocation,
			},
		}

		exampleProductImageLocations := &storage.ProductImageLocations{
			Thumbnail: exampleThumbnailLocation,
			Main:      exampleMainLocation,
			Original:  exampleOriginalLocation,
		}
		arbitaryImageSet := storage.ProductImageSet{}
		testUtil.MockImageStorage.On("CreateThumbnails", mock.Anything).
			Return(arbitaryImageSet)
		testUtil.MockImageStorage.On("StoreImages", mock.Anything, exampleSKU, mock.AnythingOfType("uint")).
			Return(exampleProductImageLocations, nil)

		testUtil.MockDB.On("CreateProductImage", mock.AnythingOfType("*sql.Tx"), mock.Anything).
			Return(exampleRootID, buildTestTime(), nil)

		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actualImages, actualPrimaryImageID, err := handleProductCreationImages(tx, testUtil.MockDB, testUtil.MockImageStorage, exampleImageInputs, exampleSKU, exampleRootID)

		assert.NoError(t, err)
		assert.Equal(t, expectedImages, actualImages, "expected and actual images should match")
		assert.Equal(t, expectedPrimaryImageID, actualPrimaryImageID, "expected and actual primary image IDs should match")
	})
}
