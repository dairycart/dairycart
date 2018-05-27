package dairyclient

import (
	"net/http"

	"github.com/dairycart/dairymodels/v1"
)

////////////////////////////////////////////////////////
//                                                    //
//                  User Functions                    //
//                                                    //
////////////////////////////////////////////////////////

// CreateUser takes a UserCreationInput and creates the user in Dairycart
func (dc *V1Client) CreateUser(nu models.UserCreationInput) (*models.User, error) {
	u := dc.buildURL(nil, "user")
	body, _ := createBodyFromStruct(nu)

	req, _ := http.NewRequest(http.MethodPost, u, body)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	ru := models.User{}
	apiErr := unmarshalBody(res, &ru)
	if apiErr != nil {
		return nil, apiErr
	}

	return &ru, nil
}

// DeleteUser deletes a user with a given ID
func (dc *V1Client) DeleteUser(userID uint64) error {
	userIDString := convertIDToString(userID)
	u := dc.buildURL(nil, "user", userIDString)
	return dc.delete(u)
}
