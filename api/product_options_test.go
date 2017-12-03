package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateCartesianProductForOptions(t *testing.T) {
	t.Parallel()

	small := models.ProductOptionValue{ID: 1, Value: "small"}
	medium := models.ProductOptionValue{ID: 2, Value: "medium"}
	large := models.ProductOptionValue{ID: 3, Value: "large"}
	red := models.ProductOptionValue{ID: 4, Value: "red"}
	green := models.ProductOptionValue{ID: 5, Value: "green"}
	blue := models.ProductOptionValue{ID: 6, Value: "blue"}
	xtraLarge := models.ProductOptionValue{ID: 7, Value: "xtra-large"}
	polyester := models.ProductOptionValue{ID: 8, Value: "polyester"}
	cotton := models.ProductOptionValue{ID: 9, Value: "cotton"}

	tt := []struct {
		in       []models.ProductOption
		expected []simpleProductOption
		len      int
	}{
		{
			in: []models.ProductOption{
				{Name: "Size", Values: []models.ProductOptionValue{small, medium, large}},
				{Name: "Color", Values: []models.ProductOptionValue{red, green, blue}},
			},
			expected: []simpleProductOption{
				{IDs: []uint64{small.ID, red.ID}, OptionSummary: "Size: small, Color: red", SKUPostfix: "small_red", OriginalValues: []models.ProductOptionValue{small, red}},
				{IDs: []uint64{small.ID, green.ID}, OptionSummary: "Size: small, Color: green", SKUPostfix: "small_green", OriginalValues: []models.ProductOptionValue{small, green}},
				{IDs: []uint64{small.ID, blue.ID}, OptionSummary: "Size: small, Color: blue", SKUPostfix: "small_blue", OriginalValues: []models.ProductOptionValue{small, blue}},
				{IDs: []uint64{medium.ID, red.ID}, OptionSummary: "Size: medium, Color: red", SKUPostfix: "medium_red", OriginalValues: []models.ProductOptionValue{medium, red}},
				{IDs: []uint64{medium.ID, green.ID}, OptionSummary: "Size: medium, Color: green", SKUPostfix: "medium_green", OriginalValues: []models.ProductOptionValue{medium, green}},
				{IDs: []uint64{medium.ID, blue.ID}, OptionSummary: "Size: medium, Color: blue", SKUPostfix: "medium_blue", OriginalValues: []models.ProductOptionValue{medium, blue}},
				{IDs: []uint64{large.ID, red.ID}, OptionSummary: "Size: large, Color: red", SKUPostfix: "large_red", OriginalValues: []models.ProductOptionValue{large, red}},
				{IDs: []uint64{large.ID, green.ID}, OptionSummary: "Size: large, Color: green", SKUPostfix: "large_green", OriginalValues: []models.ProductOptionValue{large, green}},
				{IDs: []uint64{large.ID, blue.ID}, OptionSummary: "Size: large, Color: blue", SKUPostfix: "large_blue", OriginalValues: []models.ProductOptionValue{large, blue}},
			},
			len: 9,
		},
		{
			// test that name: value pairs can be completely different sizes
			in: []models.ProductOption{
				{Name: "Size", Values: []models.ProductOptionValue{small, medium, large, xtraLarge}},
				{Name: "Color", Values: []models.ProductOptionValue{red, green, blue}},
				{Name: "Fabric", Values: []models.ProductOptionValue{polyester, cotton}},
			},
			expected: []simpleProductOption{
				{
					IDs:            []uint64{small.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: red, Fabric: polyester",
					SKUPostfix:     "small_red_polyester",
					OriginalValues: []models.ProductOptionValue{small, red, polyester},
				},
				{
					IDs:            []uint64{small.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: red, Fabric: cotton",
					SKUPostfix:     "small_red_cotton",
					OriginalValues: []models.ProductOptionValue{small, red, cotton},
				},
				{
					IDs:            []uint64{small.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: green, Fabric: polyester",
					SKUPostfix:     "small_green_polyester",
					OriginalValues: []models.ProductOptionValue{small, green, polyester},
				},
				{
					IDs:            []uint64{small.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: green, Fabric: cotton",
					SKUPostfix:     "small_green_cotton",
					OriginalValues: []models.ProductOptionValue{small, green, cotton},
				},
				{
					IDs:            []uint64{small.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: small, Color: blue, Fabric: polyester",
					SKUPostfix:     "small_blue_polyester",
					OriginalValues: []models.ProductOptionValue{small, blue, polyester},
				},
				{
					IDs:            []uint64{small.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: small, Color: blue, Fabric: cotton",
					SKUPostfix:     "small_blue_cotton",
					OriginalValues: []models.ProductOptionValue{small, blue, cotton},
				},
				{
					IDs:            []uint64{medium.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: red, Fabric: polyester",
					SKUPostfix:     "medium_red_polyester",
					OriginalValues: []models.ProductOptionValue{medium, red, polyester},
				},
				{
					IDs:            []uint64{medium.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: red, Fabric: cotton",
					SKUPostfix:     "medium_red_cotton",
					OriginalValues: []models.ProductOptionValue{medium, red, cotton},
				},
				{
					IDs:            []uint64{medium.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: green, Fabric: polyester",
					SKUPostfix:     "medium_green_polyester",
					OriginalValues: []models.ProductOptionValue{medium, green, polyester},
				},
				{
					IDs:            []uint64{medium.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: green, Fabric: cotton",
					SKUPostfix:     "medium_green_cotton",
					OriginalValues: []models.ProductOptionValue{medium, green, cotton},
				},
				{
					IDs:            []uint64{medium.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: medium, Color: blue, Fabric: polyester",
					SKUPostfix:     "medium_blue_polyester",
					OriginalValues: []models.ProductOptionValue{medium, blue, polyester},
				},
				{
					IDs:            []uint64{medium.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: medium, Color: blue, Fabric: cotton",
					SKUPostfix:     "medium_blue_cotton",
					OriginalValues: []models.ProductOptionValue{medium, blue, cotton},
				},
				{
					IDs:            []uint64{large.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: red, Fabric: polyester",
					SKUPostfix:     "large_red_polyester",
					OriginalValues: []models.ProductOptionValue{large, red, polyester},
				},
				{
					IDs:            []uint64{large.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: red, Fabric: cotton",
					SKUPostfix:     "large_red_cotton",
					OriginalValues: []models.ProductOptionValue{large, red, cotton},
				},
				{
					IDs:            []uint64{large.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: green, Fabric: polyester",
					SKUPostfix:     "large_green_polyester",
					OriginalValues: []models.ProductOptionValue{large, green, polyester},
				},
				{
					IDs:            []uint64{large.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: green, Fabric: cotton",
					SKUPostfix:     "large_green_cotton",
					OriginalValues: []models.ProductOptionValue{large, green, cotton},
				},
				{
					IDs:            []uint64{large.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: large, Color: blue, Fabric: polyester",
					SKUPostfix:     "large_blue_polyester",
					OriginalValues: []models.ProductOptionValue{large, blue, polyester},
				},
				{
					IDs:            []uint64{large.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: large, Color: blue, Fabric: cotton",
					SKUPostfix:     "large_blue_cotton",
					OriginalValues: []models.ProductOptionValue{large, blue, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, red.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: red, Fabric: polyester",
					SKUPostfix:     "xtra-large_red_polyester",
					OriginalValues: []models.ProductOptionValue{xtraLarge, red, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, red.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: red, Fabric: cotton",
					SKUPostfix:     "xtra-large_red_cotton",
					OriginalValues: []models.ProductOptionValue{xtraLarge, red, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, green.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: green, Fabric: polyester",
					SKUPostfix:     "xtra-large_green_polyester",
					OriginalValues: []models.ProductOptionValue{xtraLarge, green, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, green.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: green, Fabric: cotton",
					SKUPostfix:     "xtra-large_green_cotton",
					OriginalValues: []models.ProductOptionValue{xtraLarge, green, cotton},
				},
				{
					IDs:            []uint64{xtraLarge.ID, blue.ID, polyester.ID},
					OptionSummary:  "Size: xtra-large, Color: blue, Fabric: polyester",
					SKUPostfix:     "xtra-large_blue_polyester",
					OriginalValues: []models.ProductOptionValue{xtraLarge, blue, polyester},
				},
				{
					IDs:            []uint64{xtraLarge.ID, blue.ID, cotton.ID},
					OptionSummary:  "Size: xtra-large, Color: blue, Fabric: cotton",
					SKUPostfix:     "xtra-large_blue_cotton",
					OriginalValues: []models.ProductOptionValue{xtraLarge, blue, cotton},
				},
			},
			len: 24,
		},
	}

	for _, tc := range tt {
		actual := generateCartesianProductForOptions(tc.in)
		assert.Equal(t, tc.len, len(actual), "there should be %d simpleProductOptions, but we generated %d", tc.len, len(actual))
		assert.Equal(t, tc.expected, actual, "expected output should match actual output")
	}
}

func TestCreateProductOptionAndValuesInDBFromInput(t *testing.T) {
	exampleID := uint64(1)

	t.Run("optimal conditions", func(*testing.T) {
		exampleProductOptionCreationInput := &ProductOptionCreationInput{
			Name:   "name",
			Values: []string{"one", "two", "three"},
		}

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), nil)
		tx, err := testUtil.PlainDB.Begin()
		assert.NoError(t, err)

		actual, err := createProductOptionAndValuesInDBFromInput(tx, exampleProductOptionCreationInput, exampleID, testUtil.MockDB)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("with error creating product option value", func(*testing.T) {
		exampleProductOptionCreationInput := &ProductOptionCreationInput{
			Name:   "name",
			Values: []string{"one", "two", "three"},
		}

		testUtil := setupTestVariablesWithMock(t)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("CreateProductOption", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleID, generateExampleTimeForTests(), generateArbitraryError())
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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%d/options", exampleProductRootID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error getting product option count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOptionCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		CreatedOn:     generateExampleTimeForTests(),
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
			Return(exampleProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(false, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/product/%d/options", exampleProductOption.ProductRootID), strings.NewReader(exampleProductOptionCreationBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error getting product root", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("ProductRootExists", mock.Anything, exampleProductOption.ProductRootID).
			Return(true, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(exampleProductOption.ID, generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(exampleProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("CreateProductOptionValue", mock.Anything, mock.Anything).
			Return(exampleProductOption.ID, generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.MockDB.On("UpdateProductOption", mock.Anything, mock.Anything).
			Return(generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("GetProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return([]models.ProductOptionValue{}, nil)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving product option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(generateExampleTimeForTests(), generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("GetProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return([]models.ProductOptionValue{}, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), strings.NewReader(exampleProductOptionUpdateBody))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestProductOptionDeletionHandler(t *testing.T) {
	exampleProductOption := &models.ProductOption{
		ID:            123,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "something",
		ProductRootID: 2,
	}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, nil)
		testUtil.Mock.ExpectBegin()
		testUtil.MockDB.On("ArchiveProductOptionValuesForOption", mock.Anything, exampleProductOption.ID).
			Return(generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, sql.ErrNoRows)
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving option", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetProductOption", mock.Anything, exampleProductOption.ID).
			Return(exampleProductOption, generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(generateExampleTimeForTests(), generateArbitraryError())
		testUtil.Mock.ExpectRollback()
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

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
			Return(generateExampleTimeForTests(), nil)
		testUtil.MockDB.On("DeleteProductOption", mock.Anything, exampleProductOption.ID).
			Return(generateExampleTimeForTests(), nil)
		testUtil.Mock.ExpectCommit().WillReturnError(generateArbitraryError())
		SetupAPIRoutes(testUtil.Router, testUtil.PlainDB, testUtil.Store, testUtil.MockDB)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/product_options/%d", exampleProductOption.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
