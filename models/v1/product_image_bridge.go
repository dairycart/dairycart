package models

// ProductImageBridge represents a Dairycart product image bridge
type ProductImageBridge struct {
	ID             uint64 `json:"id"`               // id
	ProductID      uint64 `json:"product_id"`       // product_id
	ProductImageID uint64 `json:"product_image_id"` // product_image_id
}

// ProductImageBridgeCreationInput is a struct to use for creating ProductImageBridges
type ProductImageBridgeCreationInput struct {
}

// ProductImageBridgeUpdateInput is a struct to use for updating ProductImageBridges
type ProductImageBridgeUpdateInput struct {
	ProductID      uint64 `json:"product_id,omitempty"`       // product_id
	ProductImageID uint64 `json:"product_image_id,omitempty"` // product_image_id
}

type ProductImageBridgeListResponse struct {
	ListResponse
	ProductImageBridges []ProductImageBridge `json:"product_image_bridge"`
}
