package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
)

func buildNotFoundProductResponse(sku string) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The product you were looking for (sku '%s') does not exist"
		}
	`, sku)
}

func TestProductExists(t *testing.T) {
	existentSKU := "existent_sku"
	nonexistentSKU := "nonexistent_sku"

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", existentSKU):    generateHeadHandler(t, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU): generateHeadHandler(t, http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(*testing.T) {
		exists, err := c.ProductExists(existentSKU)
		assert.Nil(t, err)
		assert.True(t, exists)

	})

	t.Run("with nonexistent product", func(*testing.T) {
		exists, err := c.ProductExists(nonexistentSKU)
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}

func TestGetProduct(t *testing.T) {
	goodResponseSKU := "good"
	nonexistentSKU := "nonexistent"
	badResponseSKU := "bad"

	exampleResponse := loadExampleResponse(t, "product")

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", goodResponseSKU): generateGetHandler(t, exampleResponse, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", badResponseSKU):  generateGetHandler(t, exampleBadJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU):  generateGetHandler(t, buildNotFoundProductResponse(nonexistentSKU), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		expected := &models.Product{
			Name:               "Your Favorite Band's T-Shirt",
			Subtitle:           "A t-shirt you can wear",
			Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			OptionSummary:      "Size: Small, Color: Red",
			SKU:                "sku",
			UPC:                "",
			Manufacturer:       "Record Company",
			Brand:              "Your Favorite Band",
			Quantity:           666,
			QuantityPerPackage: 1,
			Taxable:            true,
			Price:              20,
			OnSale:             false,
			Cost:               10,
			ProductWeight:      1,
			ProductHeight:      5,
			ProductWidth:       5,
			ProductLength:      5,
			PackageWeight:      1,
			PackageHeight:      5,
			PackageWidth:       5,
			PackageLength:      5,
		}
		actual, err := c.GetProduct(goodResponseSKU)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected product doesn't match actual product")
	})

	t.Run("nonexistent product", func(*testing.T) {
		_, err := c.GetProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})

	t.Run("bad response from server", func(*testing.T) {
		_, err := c.GetProduct(badResponseSKU)
		assert.NotNil(t, err)
	})

	t.Run("with request error", func(*testing.T) {
		ts.Close()
		_, err := c.GetProduct(exampleSKU)
		assert.NotNil(t, err)
	})
}

func TestGetProducts(t *testing.T) {
	exampleGoodResponse := loadExampleResponse(t, "products")

	t.Run("normal usage", func(*testing.T) {
		expected := []models.Product{
			{
				ID:                 1,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Red",
				SKU:                "t-shirt-small-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 2,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Red",
				SKU:                "t-shirt-medium-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 3,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Large, Color: Red",
				SKU:                "t-shirt-large-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 4,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Blue",
				SKU:                "t-shirt-small-blue",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 5,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Blue",
				SKU:                "t-shirt-medium-blue",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
		}

		handlers := map[string]http.HandlerFunc{
			"/v1/products": generateGetHandler(t, exampleGoodResponse, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		actual, err := c.GetProducts(nil)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected product doesn't match actual product")
	})

	t.Run("with bad server response", func(*testing.T) {
		handlers := map[string]http.HandlerFunc{
			"/v1/products": generateGetHandler(t, exampleBadJSON, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		_, err := c.GetProducts(nil)
		assert.NotNil(t, err, "GetProducts should return an error when it receives nonsense")
	})
}

func TestCreateProduct(t *testing.T) {
	exampleProductCreationInput := models.ProductCreationInput{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		UPC:                "upc",
		Manufacturer:       "manufacturer",
		Brand:              "brand",
		Quantity:           666,
		Price:              20,
		SalePrice:          10,
		Cost:               1.23,
		ProductWeight:      9,
		ProductHeight:      9,
		ProductWidth:       9,
		ProductLength:      9,
		PackageWeight:      9,
		PackageHeight:      9,
		PackageWidth:       9,
		PackageLength:      9,
		QuantityPerPackage: 1,
	}

	t.Run("normal response", func(*testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			"/v1/product": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPost, "CreateProduct should only be making POST requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
						"price": 20,
						"sale_price": 10,
						"cost": 1.23,
						"product_weight": 9,
						"product_height": 9,
						"product_width": 9,
						"product_length": 9,
						"package_weight": 9,
						"package_height": 9,
						"package_width": 9,
						"package_length": 9,
						"quantity_per_package": 1
					}
				`
				actual := string(bodyBytes)
				assert.Equal(t, minifyJSON(t, expected), actual, "CreateProduct should attach the correct JSON to the request body")

				exampleResponse := loadExampleResponse(t, "created_product")
				fmt.Fprintf(res, exampleResponse)
			},
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := &models.Product{
			Name:               "name",
			Subtitle:           "subtitle",
			Description:        "description",
			OptionSummary:      "option_summary",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "manufacturer",
			Brand:              "brand",
			Quantity:           666,
			Price:              20,
			SalePrice:          10,
			Cost:               1.23,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
			QuantityPerPackage: 1,
		}
		actual, err := c.CreateProduct(exampleProductCreationInput)

		assert.Nil(t, err, "CreateProduct with valid input and response should never produce an error")
		assert.Equal(t, expected, actual, "expected and actual products should match")
		assert.True(t, normalEndpointCalled, "the normal endpoint should be called")
	})

	t.Run("with bad server response", func(*testing.T) {
		var badEndpointCalled bool
		handlers := map[string]http.HandlerFunc{
			"/v1/product": func(res http.ResponseWriter, req *http.Request) {
				badEndpointCalled = true
				fmt.Fprintf(res, exampleBadJSON)
			},
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.CreateProduct(exampleProductCreationInput)
		assert.NotNil(t, err, "CreateProduct should return an error when it fails to load a response")
		assert.True(t, badEndpointCalled, "the bad response endpoint should be called")
	})

	t.Run("with request error", func(*testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.CreateProduct(models.ProductCreationInput{})
		assert.NotNil(t, err, "CreateProduct should return an error when faililng to execute a request")
	})
}

// Note: this test is basically the same as TestCreateProduct, because those functions are incredibly similar, but with different
// purposes. I could probably sleep well at night with no tests for this, if only it wouldn't lower my precious coverage number.
func TestUpdateProduct(t *testing.T) {
	exampleProductUpdateInput := models.ProductUpdateInput{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		UPC:                "upc",
		Manufacturer:       "manufacturer",
		Brand:              "brand",
		Quantity:           666,
		Price:              20,
		SalePrice:          10,
		Cost:               1.23,
		ProductWeight:      9,
		ProductHeight:      9,
		ProductWidth:       9,
		ProductLength:      9,
		PackageWeight:      9,
		PackageHeight:      9,
		PackageWidth:       9,
		PackageLength:      9,
		QuantityPerPackage: 1,
	}

	t.Run("normal response", func(*testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			// UPGRADEME
			"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPatch, "UpdateProduct should only be making PATCH requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
						"price": 20,
						"sale_price": 10,
						"cost": 1.23,
						"product_weight": 9,
						"product_height": 9,
						"product_width": 9,
						"product_length": 9,
						"package_weight": 9,
						"package_height": 9,
						"package_width": 9,
						"package_length": 9,
						"quantity_per_package": 1
					}
				`
				actual := string(bodyBytes)
				assert.Equal(t, minifyJSON(t, expected), actual, "UpdateProduct should attach the correct JSON to the request body")

				exampleResponse := loadExampleResponse(t, "updated_product")
				fmt.Fprintf(res, exampleResponse)
			},
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := &models.Product{
			Name:               "name",
			Subtitle:           "subtitle",
			Description:        "description",
			OptionSummary:      "option_summary",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "manufacturer",
			Brand:              "brand",
			Quantity:           666,
			Price:              20,
			SalePrice:          10,
			Cost:               1.23,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
			QuantityPerPackage: 1,
		}
		actual, err := c.UpdateProduct(exampleSKU, exampleProductUpdateInput)
		assert.Nil(t, err, "UpdateProduct with valid input and response should never produce an error")
		assert.Equal(t, expected, actual, "expected and actual products should match")
		assert.True(t, normalEndpointCalled, "the normal endpoint should be called")
	})

	t.Run("bad response", func(*testing.T) {
		var badEndpointCalled bool
		handlers := map[string]http.HandlerFunc{
			"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
				badEndpointCalled = true
				fmt.Fprintf(res, exampleBadJSON)
			},
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.UpdateProduct(exampleSKU, exampleProductUpdateInput)
		assert.NotNil(t, err, "UpdateProduct should return an error when it fails to load a response")
		assert.True(t, badEndpointCalled, "the bad response endpoint should be called")
	})

	t.Run("with request error", func(*testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.UpdateProduct(exampleSKU, models.ProductUpdateInput{})
		assert.NotNil(t, err, "UpdateProduct should return an error when faililng to execute a request")
	})
}

func TestDeleteProduct(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "deleted_product")

	existentSKU := "existent_sku"
	nonexistentSKU := "nonexistent_sku"

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", existentSKU):    generateDeleteHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU): generateDeleteHandler(t, buildNotFoundProductResponse(nonexistentSKU), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(*testing.T) {
		err := c.DeleteProduct(existentSKU)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		err := c.DeleteProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})
}

func buildNotFoundProductRootResponse(id uint64) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The product_root you were looking for (identified by '%d') does not exist"
		}
	`, id)
}

func TestGetProductRoot(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "product_root")
	existentID := uint64(1)
	nonexistentID := uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_root/%d", existentID):    generateGetHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product_root/%d", nonexistentID): generateGetHandler(t, buildNotFoundProductRootResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		expected := &models.ProductRoot{
			ID:                 1,
			Name:               "Your Favorite Band's T-Shirt",
			Subtitle:           "A t-shirt you can wear",
			Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			SKUPrefix:          "t-shirt",
			Manufacturer:       "Record Company",
			Brand:              "Your Favorite Band",
			Taxable:            true,
			Cost:               20,
			ProductWeight:      1,
			ProductHeight:      5,
			ProductWidth:       5,
			ProductLength:      5,
			PackageWeight:      1,
			PackageHeight:      5,
			PackageWidth:       5,
			PackageLength:      5,
			QuantityPerPackage: 1,
			AvailableOn:        buildTestTime(t),
			CreatedOn:          buildTestTime(t),
			Options: []models.ProductOption{
				{
					ID:            1,
					ProductRootID: 1,
					Name:          "color",
					CreatedOn:     buildTestTime(t),
				},
				{
					ID:            2,
					ProductRootID: 1,
					Name:          "size",
					CreatedOn:     buildTestTime(t),
				},
			},
			Products: []models.Product{
				{
					ID:                 1,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-red",
					OptionSummary:      "Size: Small, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 2,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-red",
					OptionSummary:      "Size: Medium, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 3,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-red",
					OptionSummary:      "Size: Large, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 4,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-blue",
					OptionSummary:      "Size: Small, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 5,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-blue",
					OptionSummary:      "Size: Medium, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 6,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-blue",
					OptionSummary:      "Size: Large, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 7,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-green",
					OptionSummary:      "Size: Small, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 8,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-green",
					OptionSummary:      "Size: Medium, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
				{
					ID:                 9,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-green",
					OptionSummary:      "Size: Large, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        buildTestTime(t),
					CreatedOn:          buildTestTime(t),
				},
			},
		}

		actual, err := c.GetProductRoot(existentID)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected and actual product roots don't match.")
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		_, err := c.GetProductRoot(nonexistentID)
		assert.NotNil(t, err)
	})
}

func TestGetProductRoots(t *testing.T) {
	t.Run("normal usage", func(*testing.T) {
		exampleResponseJSON := loadExampleResponse(t, "product_roots")
		handlers := map[string]http.HandlerFunc{
			"/v1/product_roots": generateGetHandler(t, exampleResponseJSON, http.StatusOK),
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := []models.ProductRoot{
			{
				ID:                 5,
				Name:               "Animals As Leaders - The Joy Of Motion",
				Subtitle:           "A solid prog metal album",
				Description:        "Arbitrary description can go here because real product descriptions are technically copywritten.",
				Brand:              "Animals As Leaders",
				Manufacturer:       "Record Company",
				SKUPrefix:          "the-joy-of-motion",
				QuantityPerPackage: 1,
				ProductLength:      0.5,
				PackageHeight:      12,
				ProductWeight:      1,
				ProductWidth:       12,
				ProductHeight:      12,
				PackageLength:      0.5,
				PackageWeight:      1,
				PackageWidth:       12,
				Cost:               5,
				Taxable:            true,
				AvailableOn:        buildTestTime(t),
				CreatedOn:          buildTestTime(t),
			},
			{
				ID:                 6,
				Name:               "Mort Garson - Mother Earth's Plantasia",
				Subtitle:           "A solid synth album",
				Description:        "Arbitrary description can go here because real product descriptions are technically copywritten.",
				Brand:              "Mort Garson",
				Manufacturer:       "Record Company",
				SKUPrefix:          "mother-earths-plantasia",
				QuantityPerPackage: 1,
				ProductLength:      0.5,
				PackageHeight:      12,
				ProductWeight:      1,
				ProductWidth:       12,
				ProductHeight:      12,
				PackageLength:      0.5,
				PackageWeight:      1,
				PackageWidth:       12,
				Cost:               5,
				Taxable:            true,
				AvailableOn:        buildTestTime(t),
				CreatedOn:          buildTestTime(t),
			},
		}

		actual, err := c.GetProductRoots(nil)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("with error", func(*testing.T) {
		handlers := map[string]http.HandlerFunc{
			"/v1/product_roots": generateGetHandler(t, "{}", http.StatusInternalServerError),
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.GetProductRoots(nil)
		assert.NotNil(t, err)
	})
}

func TestDeleteProductRoot(t *testing.T) {
	existentID := uint64(1)
	nonexistentID := uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_root/%d", existentID):    generateDeleteHandler(t, "{}", http.StatusOK),
		fmt.Sprintf("/v1/product_root/%d", nonexistentID): generateDeleteHandler(t, buildNotFoundProductRootResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		err := c.DeleteProductRoot(existentID)
		assert.Nil(t, err)
	})

	t.Run("nonexistent product root", func(*testing.T) {
		err := c.DeleteProductRoot(nonexistentID)
		assert.NotNil(t, err)
	})
}

func buildNotFoundProductOptionsResponse(productID uint64) string {
	return `
		{
			"count": 2,
			"limit": 25,
			"page": 1,
			"data": null
		}
	`
}

func buildNotFoundProductOptionResponse(productID uint64) string {
	return fmt.Sprintf(`{"status":404,"message":"The product option you were looking for (id '%d') does not exist"}`, productID)
}

func TestGetProductOptions(t *testing.T) {
	existentID := uint64(1)
	nonexistentID := uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "product_options")

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%d/options", existentID):    generateGetHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%d/options", nonexistentID): generateGetHandler(t, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := []models.ProductOption{
			{
				ID:            1,
				Name:          "color",
				ProductRootID: 1,
				CreatedOn:     buildTestTime(t),
				Values: []models.ProductOptionValue{
					{
						ID:              1,
						ProductOptionID: 1,
						Value:           "red",
						CreatedOn:       buildTestTime(t),
					},
					{
						ID:              2,
						ProductOptionID: 1,
						Value:           "green",
						CreatedOn:       buildTestTime(t),
					},
					{
						ID:              3,
						ProductOptionID: 1,
						Value:           "blue",
						CreatedOn:       buildTestTime(t),
					},
				},
			},
			{
				ID:            2,
				Name:          "size",
				ProductRootID: 1,
				CreatedOn:     buildTestTime(t),
				Values: []models.ProductOptionValue{
					{
						ID:              4,
						ProductOptionID: 2,
						Value:           "small",
						CreatedOn:       buildTestTime(t),
					},
					{
						ID:              5,
						ProductOptionID: 2,
						Value:           "medium",
						CreatedOn:       buildTestTime(t),
					},
					{
						ID:              6,
						ProductOptionID: 2,
						Value:           "large",
						CreatedOn:       buildTestTime(t),
					},
				},
			},
		}

		actual, err := c.GetProductOptions(existentID, nil)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.GetProductOptions(nonexistentID, nil)
		assert.NotNil(t, err)
	})
}

func TestCreateProductOption(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "created_product_option")
	expectedBody := `
		{
			"name": "example_option",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`
	exampleInput := models.ProductOptionCreationInput{
		Name:   "example_option",
		Values: []string{"one", "two", "three"},
	}

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%d/options", existentID):    generatePostHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		fmt.Sprintf("/v1/product/%d/options", nonexistentID): generatePostHandler(t, expectedBody, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := &models.ProductOption{
			ID:            3,
			Name:          "example_option",
			ProductRootID: 1,
			CreatedOn:     buildTestTime(t),
			Values: []models.ProductOptionValue{
				{
					ID:              7,
					ProductOptionID: 3,
					Value:           "one",
					CreatedOn:       buildTestTime(t),
				},
				{
					ID:              8,
					ProductOptionID: 3,
					Value:           "two",
					CreatedOn:       buildTestTime(t),
				},
				{
					ID:              9,
					ProductOptionID: 3,
					Value:           "three",
					CreatedOn:       buildTestTime(t),
				},
			},
		}

		actual, err := c.CreateProductOption(existentID, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.CreateProductOption(nonexistentID, exampleInput)
		assert.NotNil(t, err)
	})
}

func TestUpdateProductOption(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "updated_product_option")
	// FIXME
	expectedBody := `
		{
			"name": "example_option_updated"
		}
	`
	exampleInput := models.ProductOptionUpdateInput{
		Name: "example_option_updated",
	}

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_options/%d", existentID):    generatePatchHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		fmt.Sprintf("/v1/product_options/%d", nonexistentID): generatePatchHandler(t, expectedBody, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := &models.ProductOption{
			ID:            3,
			Name:          "example_option_updated",
			ProductRootID: 1,
			CreatedOn:     buildTestTime(t),
			Values: []models.ProductOptionValue{
				{
					ID:              7,
					ProductOptionID: 3,
					Value:           "one",
					CreatedOn:       buildTestTime(t),
				},
				{
					ID:              8,
					ProductOptionID: 3,
					Value:           "two",
					CreatedOn:       buildTestTime(t),
				},
				{
					ID:              9,
					ProductOptionID: 3,
					Value:           "three",
					CreatedOn:       buildTestTime(t),
				},
			},
		}

		actual, err := c.UpdateProductOption(existentID, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.UpdateProductOption(nonexistentID, exampleInput)
		assert.NotNil(t, err)
	})
}

func TestDeleteProductOption(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "deleted_product_option")
	existentID, nonexistentID := uint64(1), uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_options/%d", existentID):    generateDeleteHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product_options/%d", nonexistentID): generateDeleteHandler(t, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(*testing.T) {
		err := c.DeleteProductOption(existentID)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		err := c.DeleteProductOption(nonexistentID)
		assert.NotNil(t, err)
	})
}

func TestCreateProductOptionValue(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "created_product_option_value")
	expectedBody := `
		{
			"value": "example_value"
		}
	`
	exampleInput := models.ProductOptionValueCreationInput{
		Value: "example_value",
	}

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_options/%d/value", existentID):    generatePostHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		fmt.Sprintf("/v1/product_options/%d/value", nonexistentID): generatePostHandler(t, expectedBody, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := &models.ProductOptionValue{
			ID:              8,
			ProductOptionID: 2,
			Value:           "large",
			CreatedOn:       buildTestTime(t),
		}

		actual, err := c.CreateProductOptionValue(existentID, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.CreateProductOptionValue(nonexistentID, exampleInput)
		assert.NotNil(t, err)
	})
}

func TestUpdateProductOptionValue(t *testing.T) {
	existentID, nonexistentID := uint64(1), uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "created_product_option_value")
	expectedBody := `
		{
			"value": "example_value"
		}
	`
	exampleInput := models.ProductOptionValueUpdateInput{
		Value: "example_value",
	}

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_option_values/%d", existentID):    generatePatchHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		fmt.Sprintf("/v1/product_option_values/%d", nonexistentID): generatePatchHandler(t, expectedBody, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal operation", func(*testing.T) {
		expected := &models.ProductOptionValue{
			ID:              8,
			ProductOptionID: 2,
			Value:           "large",
			CreatedOn:       buildTestTime(t),
		}

		actual, err := c.UpdateProductOptionValue(existentID, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.UpdateProductOptionValue(nonexistentID, exampleInput)
		assert.NotNil(t, err)
	})
}

func TestDeleteProductOptionValue(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "deleted_product_option_value")
	existentID, nonexistentID := uint64(1), uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_option_values/%d", existentID):    generateDeleteHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product_option_values/%d", nonexistentID): generateDeleteHandler(t, buildNotFoundProductOptionResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(*testing.T) {
		err := c.DeleteProductOptionValue(existentID)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		err := c.DeleteProductOptionValue(nonexistentID)
		assert.NotNil(t, err)
	})
}
