package storage

import (
	"database/sql"
	"image"
	"time"

	"github.com/dairycart/dairymodels/v1"
)

// Querier is a generic interface that either *sql.DB or *sql.Tx can satisfy
type Querier interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Storer interface {
	// WebhookExecutionLogs
	GetWebhookExecutionLog(Querier, uint64) (*models.WebhookExecutionLog, error)
	GetWebhookExecutionLogList(Querier, *models.QueryFilter) ([]models.WebhookExecutionLog, error)
	GetWebhookExecutionLogCount(Querier, *models.QueryFilter) (uint64, error)
	WebhookExecutionLogExists(Querier, uint64) (bool, error)
	CreateWebhookExecutionLog(Querier, *models.WebhookExecutionLog) (newID uint64, createdOn time.Time, e error)
	UpdateWebhookExecutionLog(Querier, *models.WebhookExecutionLog) (time.Time, error)
	DeleteWebhookExecutionLog(Querier, uint64) (time.Time, error)

	// ProductOptions
	GetProductOption(Querier, uint64) (*models.ProductOption, error)
	GetProductOptionList(Querier, *models.QueryFilter) ([]models.ProductOption, error)
	GetProductOptionCount(Querier, *models.QueryFilter) (uint64, error)
	ProductOptionExists(Querier, uint64) (bool, error)
	CreateProductOption(Querier, *models.ProductOption) (newID uint64, createdOn time.Time, e error)
	UpdateProductOption(Querier, *models.ProductOption) (time.Time, error)
	DeleteProductOption(Querier, uint64) (time.Time, error)
	ArchiveProductOptionsWithProductRootID(Querier, uint64) (time.Time, error)
	ProductOptionWithNameExistsForProductRoot(Querier, string, uint64) (bool, error)
	GetProductOptionsByProductRootID(Querier, uint64) ([]models.ProductOption, error)

	// Users
	GetUser(Querier, uint64) (*models.User, error)
	GetUserList(Querier, *models.QueryFilter) ([]models.User, error)
	GetUserCount(Querier, *models.QueryFilter) (uint64, error)
	UserExists(Querier, uint64) (bool, error)
	CreateUser(Querier, *models.User) (newID uint64, createdOn time.Time, e error)
	UpdateUser(Querier, *models.User) (time.Time, error)
	DeleteUser(Querier, uint64) (time.Time, error)
	GetUserByUsername(Querier, string) (*models.User, error)
	UserWithUsernameExists(Querier, string) (bool, error)

	// PasswordResetTokens
	GetPasswordResetToken(Querier, uint64) (*models.PasswordResetToken, error)
	GetPasswordResetTokenList(Querier, *models.QueryFilter) ([]models.PasswordResetToken, error)
	GetPasswordResetTokenCount(Querier, *models.QueryFilter) (uint64, error)
	PasswordResetTokenExists(Querier, uint64) (bool, error)
	CreatePasswordResetToken(Querier, *models.PasswordResetToken) (newID uint64, createdOn time.Time, e error)
	UpdatePasswordResetToken(Querier, *models.PasswordResetToken) (time.Time, error)
	DeletePasswordResetToken(Querier, uint64) (time.Time, error)
	PasswordResetTokenForUserIDExists(Querier, uint64) (bool, error)
	PasswordResetTokenWithTokenExists(Querier, string) (bool, error)

	// ProductImages
	GetProductImage(Querier, uint64) (*models.ProductImage, error)
	GetProductImageList(Querier, *models.QueryFilter) ([]models.ProductImage, error)
	GetProductImageCount(Querier, *models.QueryFilter) (uint64, error)
	ProductImageExists(Querier, uint64) (bool, error)
	CreateProductImage(Querier, *models.ProductImage) (newID uint64, createdOn time.Time, e error)
	UpdateProductImage(Querier, *models.ProductImage) (time.Time, error)
	DeleteProductImage(Querier, uint64) (time.Time, error)
	GetProductImagesByProductID(Querier, uint64) ([]models.ProductImage, error)
	SetPrimaryProductImageForProduct(Querier, uint64, uint64) (time.Time, error)

	// ProductVariantBridge
	GetProductVariantBridge(Querier, uint64) (*models.ProductVariantBridge, error)
	GetProductVariantBridgeList(Querier, *models.QueryFilter) ([]models.ProductVariantBridge, error)
	GetProductVariantBridgeCount(Querier, *models.QueryFilter) (uint64, error)
	ProductVariantBridgeExists(Querier, uint64) (bool, error)
	CreateProductVariantBridge(Querier, *models.ProductVariantBridge) (newID uint64, createdOn time.Time, e error)
	UpdateProductVariantBridge(Querier, *models.ProductVariantBridge) (time.Time, error)
	DeleteProductVariantBridge(Querier, uint64) (time.Time, error)
	ArchiveProductVariantBridgesWithProductRootID(Querier, uint64) (time.Time, error)
	DeleteProductVariantBridgeByProductID(Querier, uint64) (time.Time, error)
	CreateMultipleProductVariantBridgesForProductID(Querier, uint64, []uint64) error

	// LoginAttempts
	GetLoginAttempt(Querier, uint64) (*models.LoginAttempt, error)
	GetLoginAttemptList(Querier, *models.QueryFilter) ([]models.LoginAttempt, error)
	GetLoginAttemptCount(Querier, *models.QueryFilter) (uint64, error)
	LoginAttemptExists(Querier, uint64) (bool, error)
	CreateLoginAttempt(Querier, *models.LoginAttempt) (newID uint64, createdOn time.Time, e error)
	UpdateLoginAttempt(Querier, *models.LoginAttempt) (time.Time, error)
	DeleteLoginAttempt(Querier, uint64) (time.Time, error)
	LoginAttemptsHaveBeenExhausted(Querier, string) (bool, error)

	// Webhooks
	GetWebhook(Querier, uint64) (*models.Webhook, error)
	GetWebhookList(Querier, *models.QueryFilter) ([]models.Webhook, error)
	GetWebhookCount(Querier, *models.QueryFilter) (uint64, error)
	WebhookExists(Querier, uint64) (bool, error)
	CreateWebhook(Querier, *models.Webhook) (newID uint64, createdOn time.Time, e error)
	UpdateWebhook(Querier, *models.Webhook) (time.Time, error)
	DeleteWebhook(Querier, uint64) (time.Time, error)
	GetWebhooksByEventType(db Querier, eventType string) ([]models.Webhook, error)

	// Discounts
	GetDiscount(Querier, uint64) (*models.Discount, error)
	GetDiscountList(Querier, *models.QueryFilter) ([]models.Discount, error)
	GetDiscountCount(Querier, *models.QueryFilter) (uint64, error)
	DiscountExists(Querier, uint64) (bool, error)
	CreateDiscount(Querier, *models.Discount) (newID uint64, createdOn time.Time, e error)
	UpdateDiscount(Querier, *models.Discount) (time.Time, error)
	DeleteDiscount(Querier, uint64) (time.Time, error)
	GetDiscountByCode(Querier, string) (*models.Discount, error)

	// ProductRoots
	GetProductRoot(Querier, uint64) (*models.ProductRoot, error)
	GetProductRootList(Querier, *models.QueryFilter) ([]models.ProductRoot, error)
	GetProductRootCount(Querier, *models.QueryFilter) (uint64, error)
	ProductRootExists(Querier, uint64) (bool, error)
	CreateProductRoot(Querier, *models.ProductRoot) (newID uint64, createdOn time.Time, e error)
	UpdateProductRoot(Querier, *models.ProductRoot) (time.Time, error)
	DeleteProductRoot(Querier, uint64) (time.Time, error)
	ProductRootWithSKUPrefixExists(Querier, string) (bool, error)

	// Products
	GetProduct(Querier, uint64) (*models.Product, error)
	GetProductList(Querier, *models.QueryFilter) ([]models.Product, error)
	GetProductCount(Querier, *models.QueryFilter) (uint64, error)
	ProductExists(Querier, uint64) (bool, error)
	CreateProduct(Querier, *models.Product) (newID uint64, createdOn time.Time, availableOn time.Time, e error)
	UpdateProduct(Querier, *models.Product) (time.Time, error)
	DeleteProduct(Querier, uint64) (time.Time, error)
	ArchiveProductsWithProductRootID(Querier, uint64) (time.Time, error)
	GetProductBySKU(Querier, string) (*models.Product, error)
	ProductWithSKUExists(Querier, string) (bool, error)
	GetProductsByProductRootID(Querier, uint64) ([]models.Product, error)

	// ProductOptionValues
	GetProductOptionValue(Querier, uint64) (*models.ProductOptionValue, error)
	GetProductOptionValueList(Querier, *models.QueryFilter) ([]models.ProductOptionValue, error)
	GetProductOptionValueCount(Querier, *models.QueryFilter) (uint64, error)
	ProductOptionValueExists(Querier, uint64) (bool, error)
	CreateProductOptionValue(Querier, *models.ProductOptionValue) (newID uint64, createdOn time.Time, e error)
	UpdateProductOptionValue(Querier, *models.ProductOptionValue) (time.Time, error)
	DeleteProductOptionValue(Querier, uint64) (time.Time, error)
	ArchiveProductOptionValuesWithProductRootID(Querier, uint64) (time.Time, error)
	ProductOptionValueForOptionIDExists(Querier, uint64, string) (bool, error)
	ArchiveProductOptionValuesForOption(Querier, uint64) (time.Time, error)
	GetProductOptionValuesForOption(Querier, uint64) ([]models.ProductOptionValue, error)
}

type ProductImageSet struct {
	Thumbnail image.Image
	Main      image.Image
	Original  image.Image
}

type ProductImageLocations struct {
	Thumbnail string
	Main      string
	Original  string
}

type ImageStorer interface {
	CreateThumbnails(img image.Image) ProductImageSet
	StoreImages(imgset ProductImageSet, sku string, id uint) (*ProductImageLocations, error)
}
