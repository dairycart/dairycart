package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuildProductsFromOptions(t *testing.T) {
	t.Parallel()

	small := models.ProductOptionValue{ID: 1, Value: "small"}
	medium := models.ProductOptionValue{ID: 2, Value: "medium"}
	large := models.ProductOptionValue{ID: 3, Value: "large"}
	red := models.ProductOptionValue{ID: 4, Value: "red"}
	green := models.ProductOptionValue{ID: 5, Value: "green"}
	blue := models.ProductOptionValue{ID: 6, Value: "blue"}
	// xtraLarge := models.ProductOptionValue{ID: 7, Value: "xtra-large"}
	// polyester := models.ProductOptionValue{ID: 8, Value: "polyester"}
	// cotton := models.ProductOptionValue{ID: 9, Value: "cotton"}

	tt := []struct {
		input     *models.ProductCreationInput
		inOptions []models.ProductOption
		expected  []*models.Product
	}{
		{
			input: &models.ProductCreationInput{
				SKU: "t-shirt",
				Options: []models.ProductOptionCreationInput{
					{Name: "Size", Values: []string{"small", "medium", "large"}},
					{Name: "Color", Values: []string{"red", "green", "blue"}},
				},
			},
			inOptions: []models.ProductOption{
				{
					Name: "Size",
					Values: []models.ProductOptionValue{
						small,
						medium,
						large,
					},
				},
				{
					Name: "Color",
					Values: []models.ProductOptionValue{
						red,
						green,
						blue,
					},
				},
			},
			expected: []*models.Product{
				{
					OptionSummary: "Size: small, Color: red",
					SKU:           "t-shirt_small_red",
					ApplicableOptionValues: []models.ProductOptionValue{
						small,
						red,
					},
				},
				{
					OptionSummary: "Size: small, Color: green",
					SKU:           "t-shirt_small_green",
					ApplicableOptionValues: []models.ProductOptionValue{
						small,
						green,
					},
				},
				{
					OptionSummary: "Size: small, Color: blue",
					SKU:           "t-shirt_small_blue",
					ApplicableOptionValues: []models.ProductOptionValue{
						small,
						blue,
					},
				},
				{
					OptionSummary: "Size: medium, Color: red",
					SKU:           "t-shirt_medium_red",
					ApplicableOptionValues: []models.ProductOptionValue{
						medium,
						red,
					},
				},
				{
					OptionSummary: "Size: medium, Color: green",
					SKU:           "t-shirt_medium_green",
					ApplicableOptionValues: []models.ProductOptionValue{
						medium,
						green,
					},
				},
				{
					OptionSummary: "Size: medium, Color: blue",
					SKU:           "t-shirt_medium_blue",
					ApplicableOptionValues: []models.ProductOptionValue{
						medium,
						blue,
					},
				},
				{
					OptionSummary: "Size: large, Color: red",
					SKU:           "t-shirt_large_red",
					ApplicableOptionValues: []models.ProductOptionValue{
						large,
						red,
					},
				},
				{
					OptionSummary: "Size: large, Color: green",
					SKU:           "t-shirt_large_green",
					ApplicableOptionValues: []models.ProductOptionValue{
						large,
						green,
					},
				},
				{
					OptionSummary: "Size: large, Color: blue",
					SKU:           "t-shirt_large_blue",
					ApplicableOptionValues: []models.ProductOptionValue{
						large,
						blue,
					},
				},
			},
		},
	}

	for _, tc := range tt {
		actual := buildProductsFromOptions(tc.input, tc.inOptions)
		assert.Equal(t, tc.expected, actual, "expected output should match actual output")
	}
}

func TestCreateProductOptionAndValuesInDBFromInput(t *testing.T) {
	exampleID := uint64(1)

	t.Run("optimal conditions", func(*testing.T) {
		exampleProductOptionCreationInput := models.ProductOptionCreationInput{
			Name:   "name",
			Values: []string{"one", "two", "three"},
		}

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), nil)
		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actual, err := createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionCreationInput, exampleID, testUtil.MockDB)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("with error creating product option value", func(*testing.T) {
		exampleProductOptionCreationInput := models.ProductOptionCreationInput{
			Name:   "name",
			Values: []string{"one", "two", "three"},
		}

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleID, buildTestTime(), generateArbitraryError())
		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		_, err = createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionCreationInput, exampleID, testUtil.MockDB)
		assert.NotNil(t, err)
	})
}

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestProductOptionListHandler(t *testing.T) {
	exampleProductRootID := uint64(1)

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetProductOptionsByProductRootID", mock.Anything, exampleProductRootID).
			Return([]models.ProductOption{}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%d/options", exampleProductRootID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error getting product option count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%d/options", exampleProductRootID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error getting product options", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetProductOptionsByProductRootID", mock.Anything, exampleProductRootID).
			Return([]models.ProductOption{}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%d/options", exampleProductRootID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

}

func TestProductOptionCreationHandler(t *testing.T) {
	exampleProductOptionCreationBody := `
		{
			"name": "something",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(false, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error getting product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with already existent product option name", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(true, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error checking for duplicate product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(false, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error setting up databse transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(false, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error creating product option and values", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error committing transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, nil)
		testUtil.MockDB.On("ProductOptionWithNameExistsForProductRoot", mock.Anything, exampleProductOption.Name, exampleProductOption.ProductRootID).
			Return(false, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, buildTestTime(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, buildTestTime(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductOptionUpdateHandler(t *testing.T) {
	exampleProductOptionUpdateBody := `
		{
			"name": "something else"
		}
	`
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.MockDB.On("UpdateProductOption", mock.Anything, mock.Anything).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("GetProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return([]models.ProductOptionValue{}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.MockDB.On("UpdateProductOption", mock.Anything, mock.Anything).
			Return(buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving product option values", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.MockDB.On("UpdateProductOption", mock.Anything, mock.Anything).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("GetProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return([]models.ProductOptionValue{}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductOptionDeletionHandler(t *testing.T) {
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     buildTestTime(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), nil)
		testUtil.Mock.ExpectCommit()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error beginning transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving product option values", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error committing transaction", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(buildTestTime(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRoutes(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
