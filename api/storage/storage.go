package storage

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

type Storage interface {
	// Basic Database Stuff
	Begin() (*sql.Tx, error)

	// LoginAttempts
	GetLoginAttempt(uint64) (*models.LoginAttempt, error)
	LoginAttemptExists(uint64) (bool, error)
	CreateLoginAttempt(*models.LoginAttempt) (uint64, time.Time, error)
	UpdateLoginAttempt(*models.LoginAttempt) (time.Time, error)
	DeleteLoginAttempt(uint64, *sql.Tx) (time.Time, error)

	// ProductRoots
	GetProductRoot(uint64) (*models.ProductRoot, error)
	ProductRootExists(uint64) (bool, error)
	CreateProductRoot(*models.ProductRoot) (uint64, time.Time, error)
	UpdateProductRoot(*models.ProductRoot) (time.Time, error)
	DeleteProductRoot(uint64, *sql.Tx) (time.Time, error)

	// ProductVariantBridge
	GetProductVariantBridge(uint64) (*models.ProductVariantBridge, error)
	ProductVariantBridgeExists(uint64) (bool, error)
	CreateProductVariantBridge(*models.ProductVariantBridge) (uint64, time.Time, error)
	UpdateProductVariantBridge(*models.ProductVariantBridge) (time.Time, error)
	DeleteProductVariantBridge(uint64, *sql.Tx) (time.Time, error)

	// Users
	GetUser(uint64) (*models.User, error)
	UserExists(uint64) (bool, error)
	CreateUser(*models.User) (uint64, time.Time, error)
	UpdateUser(*models.User) (time.Time, error)
	DeleteUser(uint64, *sql.Tx) (time.Time, error)

	// ProductOptions
	GetProductOption(uint64) (*models.ProductOption, error)
	ProductOptionExists(uint64) (bool, error)
	CreateProductOption(*models.ProductOption) (uint64, time.Time, error)
	UpdateProductOption(*models.ProductOption) (time.Time, error)
	DeleteProductOption(uint64, *sql.Tx) (time.Time, error)

	// Products
	GetProduct(uint64) (*models.Product, error)
	ProductExists(uint64) (bool, error)
	CreateProduct(*models.Product) (uint64, time.Time, time.Time, error)
	UpdateProduct(*models.Product) (time.Time, error)
	DeleteProduct(uint64, *sql.Tx) (time.Time, error)
	GetProductBySKU(string) (*models.Product, error)
	ProductWithSKUExists(string) (bool, error)

	// Discounts
	GetDiscount(uint64) (*models.Discount, error)
	DiscountExists(uint64) (bool, error)
	CreateDiscount(*models.Discount) (uint64, time.Time, error)
	UpdateDiscount(*models.Discount) (time.Time, error)
	DeleteDiscount(uint64, *sql.Tx) (time.Time, error)
	GetDiscountByCode(string) (*models.Discount, error)

	// PasswordResetTokens
	GetPasswordResetToken(uint64) (*models.PasswordResetToken, error)
	PasswordResetTokenExists(uint64) (bool, error)
	CreatePasswordResetToken(*models.PasswordResetToken) (uint64, time.Time, error)
	UpdatePasswordResetToken(*models.PasswordResetToken) (time.Time, error)
	DeletePasswordResetToken(uint64, *sql.Tx) (time.Time, error)

	// ProductOptionValues
	GetProductOptionValue(uint64) (*models.ProductOptionValue, error)
	ProductOptionValueExists(uint64) (bool, error)
	CreateProductOptionValue(*models.ProductOptionValue) (uint64, time.Time, error)
	UpdateProductOptionValue(*models.ProductOptionValue) (time.Time, error)
	DeleteProductOptionValue(uint64, *sql.Tx) (time.Time, error)
}
