package dairytest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/dairycart/dairycart/models/v1"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func interfaceArgIsNotPointerOrNil(i interface{}) error {
	if i == nil {
		return errors.New("unmarshalBody cannot accept nil values")
	}
	isNotPtr := reflect.TypeOf(i).Kind() != reflect.Ptr
	if isNotPtr {
		return errors.New("unmarshalBody can only accept pointers")
	}
	return nil
}

func unmarshalBody(t *testing.T, res *http.Response, dest interface{}) {
	t.Helper()
	// These paths should only ever be reached in tests, an should never be encountered by an end user.
	require.Nil(t, interfaceArgIsNotPointerOrNil(dest), "unmarshalBody can only accept pointers")

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)

	require.Nil(t, json.Unmarshal(bodyBytes, &dest))
}

func logBodyAndResetResponse(t *testing.T, resp *http.Response) {
	t.Helper()
	respStr := turnResponseBodyIntoString(t, resp)
	log.Printf(`

		'%s'

	`, respStr)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(respStr)))
}

func convertCreationInputToProduct(in models.ProductCreationInput) models.Product {
	np := models.Product{
		Name:               in.Name,
		Subtitle:           in.Subtitle,
		Description:        in.Description,
		SKU:                in.SKU,
		UPC:                in.UPC,
		Manufacturer:       in.Manufacturer,
		Brand:              in.Brand,
		Quantity:           in.Quantity,
		QuantityPerPackage: in.QuantityPerPackage,
		Taxable:            in.Taxable,
		Price:              in.Price,
		OnSale:             in.OnSale,
		SalePrice:          in.SalePrice,
		Cost:               in.Cost,
		ProductWeight:      in.ProductWeight,
		ProductHeight:      in.ProductHeight,
		ProductWidth:       in.ProductWidth,
		ProductLength:      in.ProductLength,
		PackageWeight:      in.PackageWeight,
		PackageHeight:      in.PackageHeight,
		PackageWidth:       in.PackageWidth,
		PackageLength:      in.PackageLength,
	}
	if in.AvailableOn != nil {
		np.AvailableOn = in.AvailableOn.Time
	}
	return np
}
