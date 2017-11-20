package dairymock

import (
	"strings"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m MockDB) GetProductBySKU(sku string) (*models.Product, error) {
	fn := "GetProductBySKU"
	im, ok := m.InstructionMap[fn]
	if ok && im.CallCount < im.CallLimit {
		exampleProduct := &models.Product{
			ID:            2,
			CreatedOn:     generateExampleTimeForTests(),
			SKU:           sku,
			Name:          strings.ToTitle(sku),
			UPC:           "1234567890",
			Quantity:      123,
			Price:         99.99,
			Cost:          50.00,
			Description:   "This is a product description.",
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
			AvailableOn:   generateExampleTimeForTests(),
		}
		return exampleProduct, nil
	}
	return nil, im.Error
}

func (m MockDB) GetProductByID(ID uint64) (*models.Product, error) {
	fn := "GetProductByID"
	im, ok := m.InstructionMap[fn]
	if ok && im.CallCount < im.CallLimit {
		exampleProduct := &models.Product{
			ID:            ID,
			CreatedOn:     generateExampleTimeForTests(),
			SKU:           "sku",
			Name:          "Product Name",
			UPC:           "1234567890",
			Quantity:      123,
			Price:         99.99,
			Cost:          50.00,
			Description:   "This is a product description.",
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
			AvailableOn:   generateExampleTimeForTests(),
		}
		im.CallCount++
		return exampleProduct, nil
	}
	return nil, im.Error
}
