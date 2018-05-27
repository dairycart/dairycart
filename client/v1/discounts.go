package dairyclient

import (
	"github.com/dairycart/dairymodels/v1"
)

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetDiscountByID(discountID uint64) (*models.Discount, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	d := models.Discount{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (dc *V1Client) GetDiscounts(queryFilter map[string]string) ([]models.Discount, error) {
	u := dc.buildURL(nil, "discounts")
	d := &models.DiscountListResponse{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return d.Discounts, nil
}

func (dc *V1Client) CreateDiscount(nd models.DiscountCreationInput) (*models.Discount, error) {
	d := models.Discount{}
	u := dc.buildURL(nil, "discount")

	err := dc.post(u, nd, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dc *V1Client) UpdateDiscount(discountID uint64, ud models.DiscountUpdateInput) (*models.Discount, error) {
	d := models.Discount{}
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)

	err := dc.patch(u, ud, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dc *V1Client) DeleteDiscount(discountID uint64) error {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	return dc.delete(u)
}
