package storage

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

// Querier is a generic interface that either *sql.DB or *sql.Tx can satisfy
type Querier interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Storer interface {
	// ProductOptions
	GetProductOption(Querier, uint64) (*models.ProductOption, error)
	GetProductOptionList(Querier, *models.QueryFilter) ([]models.ProductOption, error)
	GetProductOptionCount(Querier, *models.QueryFilter) (uint64, error)
	ProductOptionExists(Querier, uint64) (bool, error)
	CreateProductOption(Querier, *models.ProductOption) (uint64, time.Time, error)
	UpdateProductOption(Querier, *models.ProductOption) (time.Time, error)
	DeleteProductOption(Querier, uint64) (time.Time, error)
	ArchiveProductOptionsWithProductRootID(Querier, uint64) (time.Time, error)
	ProductOptionWithNameExistsForProductRoot(Querier, string, uint64) (bool, error)
	GetProductOptionsByProductRootID(Querier, uint64) ([]models.ProductOption, error)

	// LoginAttempts
	GetLoginAttempt(Querier, uint64) (*models.LoginAttempt, error)
	GetLoginAttemptList(Querier, *models.QueryFilter) ([]models.LoginAttempt, error)
	GetLoginAttemptCount(Querier, *models.QueryFilter) (uint64, error)
	LoginAttemptExists(Querier, uint64) (bool, error)
	CreateLoginAttempt(Querier, *models.LoginAttempt) (uint64, time.Time, error)
	UpdateLoginAttempt(Querier, *models.LoginAttempt) (time.Time, error)
	DeleteLoginAttempt(Querier, uint64) (time.Time, error)
	LoginAttemptsHaveBeenExhausted(Querier, string) (bool, error)

	// Discounts
	GetDiscount(Querier, uint64) (*models.Discount, error)
	GetDiscountList(Querier, *models.QueryFilter) ([]models.Discount, error)
	GetDiscountCount(Querier, *models.QueryFilter) (uint64, error)
	DiscountExists(Querier, uint64) (bool, error)
	CreateDiscount(Querier, *models.Discount) (uint64, time.Time, error)
	UpdateDiscount(Querier, *models.Discount) (time.Time, error)
	DeleteDiscount(Querier, uint64) (time.Time, error)
	GetDiscountByCode(Querier, string) (*models.Discount, error)

	// ProductVariantBridge
	GetProductVariantBridge(Querier, uint64) (*models.ProductVariantBridge, error)
	GetProductVariantBridgeList(Querier, *models.QueryFilter) ([]models.ProductVariantBridge, error)
	GetProductVariantBridgeCount(Querier, *models.QueryFilter) (uint64, error)
	ProductVariantBridgeExists(Querier, uint64) (bool, error)
	CreateProductVariantBridge(Querier, *models.ProductVariantBridge) (uint64, time.Time, error)
	UpdateProductVariantBridge(Querier, *models.ProductVariantBridge) (time.Time, error)
	DeleteProductVariantBridge(Querier, uint64) (time.Time, error)
	ArchiveProductVariantBridgesWithProductRootID(Querier, uint64) (time.Time, error)
	DeleteProductVariantBridgeByProductID(Querier, uint64) (time.Time, error)
	CreateMultipleProductVariantBridgesForProductID(Querier, uint64, []uint64) error

	// Users
	GetUser(Querier, uint64) (*models.User, error)
	GetUserList(Querier, *models.QueryFilter) ([]models.User, error)
	GetUserCount(Querier, *models.QueryFilter) (uint64, error)
	UserExists(Querier, uint64) (bool, error)
	CreateUser(Querier, *models.User) (uint64, time.Time, error)
	UpdateUser(Querier, *models.User) (time.Time, error)
	DeleteUser(Querier, uint64) (time.Time, error)
	GetUserByUsername(Querier, string) (*models.User, error)
	UserWithUsernameExists(Querier, string) (bool, error)

	// Products
	GetProduct(Querier, uint64) (*models.Product, error)
	GetProductList(Querier, *models.QueryFilter) ([]models.Product, error)
	GetProductCount(Querier, *models.QueryFilter) (uint64, error)
	ProductExists(Querier, uint64) (bool, error)
	CreateProduct(Querier, *models.Product) (uint64, time.Time, time.Time, error)
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
	CreateProductOptionValue(Querier, *models.ProductOptionValue) (uint64, time.Time, error)
	UpdateProductOptionValue(Querier, *models.ProductOptionValue) (time.Time, error)
	DeleteProductOptionValue(Querier, uint64) (time.Time, error)
	ArchiveProductOptionValuesWithProductRootID(Querier, uint64) (time.Time, error)
	ProductOptionValueForOptionIDExists(Querier, uint64, string) (bool, error)
	ArchiveProductOptionValuesForOption(Querier, uint64) (time.Time, error)
	GetProductOptionValuesForOption(Querier, uint64) ([]models.ProductOptionValue, error)

	// PasswordResetTokens
	GetPasswordResetToken(Querier, uint64) (*models.PasswordResetToken, error)
	GetPasswordResetTokenList(Querier, *models.QueryFilter) ([]models.PasswordResetToken, error)
	GetPasswordResetTokenCount(Querier, *models.QueryFilter) (uint64, error)
	PasswordResetTokenExists(Querier, uint64) (bool, error)
	CreatePasswordResetToken(Querier, *models.PasswordResetToken) (uint64, time.Time, error)
	UpdatePasswordResetToken(Querier, *models.PasswordResetToken) (time.Time, error)
	DeletePasswordResetToken(Querier, uint64) (time.Time, error)
	PasswordResetTokenForUserIDExists(Querier, uint64) (bool, error)
	PasswordResetTokenWithTokenExists(Querier, string) (bool, error)

	// ProductRoots
	GetProductRoot(Querier, uint64) (*models.ProductRoot, error)
	GetProductRootList(Querier, *models.QueryFilter) ([]models.ProductRoot, error)
	GetProductRootCount(Querier, *models.QueryFilter) (uint64, error)
	ProductRootExists(Querier, uint64) (bool, error)
	CreateProductRoot(Querier, *models.ProductRoot) (uint64, time.Time, error)
	UpdateProductRoot(Querier, *models.ProductRoot) (time.Time, error)
	DeleteProductRoot(Querier, uint64) (time.Time, error)
	ProductRootWithSKUPrefixExists(Querier, string) (bool, error)
}
