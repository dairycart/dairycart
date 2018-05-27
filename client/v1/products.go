package dairyclient

import (
	"github.com/dairycart/dairymodels/v1"
)

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) ProductExists(sku string) (bool, error) {
	u := dc.buildURL(nil, "product", sku)
	return dc.exists(u)
}

func (dc *V1Client) GetProduct(sku string) (*models.Product, error) {
	u := dc.buildURL(nil, "product", sku)
	p := models.Product{}

	err := dc.get(u, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) GetProducts(queryFilter map[string]string) ([]models.Product, error) {
	u := dc.buildURL(queryFilter, "products")
	pl := &models.ProductListResponse{}

	err := dc.get(u, &pl)
	if err != nil {
		return nil, err
	}

	return pl.Products, nil
}

func (dc *V1Client) CreateProduct(np models.ProductCreationInput) (*models.Product, error) {
	p := models.Product{}
	u := dc.buildURL(nil, "product")

	err := dc.post(u, np, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (dc *V1Client) UpdateProduct(sku string, up models.ProductUpdateInput) (*models.Product, error) {
	p := models.Product{}
	u := dc.buildURL(nil, "product", sku)

	err := dc.patch(u, up, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (dc *V1Client) DeleteProduct(sku string) error {
	u := dc.buildURL(nil, "product", sku)
	return dc.delete(u)
}

////////////////////////////////////////////////////////
//                                                    //
//              Product Root Functions                //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetProductRoot(rootID uint64) (*models.ProductRoot, error) {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)

	r := models.ProductRoot{}
	err := dc.get(u, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (dc *V1Client) GetProductRoots(queryFilter map[string]string) ([]models.ProductRoot, error) {
	u := dc.buildURL(queryFilter, "product_roots")

	rl := &models.ProductRootListResponse{}
	err := dc.get(u, &rl)
	if err != nil {
		return nil, err
	}

	return rl.ProductRoots, nil
}

func (dc *V1Client) DeleteProductRoot(rootID uint64) error {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)
	return dc.delete(u)
}

////////////////////////////////////////////////////////
//                                                    //
//             Product Option Functions               //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetProductOptions(productID uint64, queryFilter map[string]string) ([]models.ProductOption, error) {
	productIDString := convertIDToString(productID)
	u := dc.buildURL(queryFilter, "product", productIDString, "options")
	ol := &models.ProductOptionListResponse{}

	err := dc.get(u, &ol)
	if err != nil {
		return nil, err
	}

	return ol.ProductOptions, nil
}

func (dc *V1Client) CreateProductOption(productRootID uint64, no models.ProductOptionCreationInput) (*models.ProductOption, error) {
	productRootIDString := convertIDToString(productRootID)
	o := models.ProductOption{}
	u := dc.buildURL(nil, "product", productRootIDString, "options")

	err := dc.post(u, no, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (dc *V1Client) UpdateProductOption(optionID uint64, uo models.ProductOptionUpdateInput) (*models.ProductOption, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString)
	o := models.ProductOption{}

	err := dc.patch(u, uo, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (dc *V1Client) DeleteProductOption(optionID uint64) error {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString)
	return dc.delete(u)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) CreateProductOptionValue(optionID uint64, nv models.ProductOptionValueCreationInput) (*models.ProductOptionValue, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString, "value")
	v := models.ProductOptionValue{}

	err := dc.post(u, nv, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (dc *V1Client) UpdateProductOptionValue(valueID uint64, uv models.ProductOptionValueUpdateInput) (*models.ProductOptionValue, error) {
	valueIDString := convertIDToString(valueID)
	u := dc.buildURL(nil, "product_option_values", valueIDString)
	v := models.ProductOptionValue{}

	err := dc.patch(u, uv, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (dc *V1Client) DeleteProductOptionValue(optionID uint64) error {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_option_values", optionIDString)
	return dc.delete(u)
}
