package storage

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

type Storage interface {
	// Discounts
	GetDiscount(id uint64) (models.Discount, error)
	CreateDiscount(models.Discount) (uint64, time.Time, error)
	UpdateDiscount(models.Discount) (time.Time, error)
	DeleteDiscount(id uint64) (time.Time, error)
	GetDiscountByCode(string) (models.Discount, error)

	// Users
	GetUser(id uint64) (models.User, error)
	CreateUser(models.User) (uint64, time.Time, error)
	UpdateUser(models.User) (time.Time, error)
	DeleteUser(id uint64) (time.Time, error)

	// ProductRoots
	GetProductRoot(id uint64) (models.ProductRoot, error)
	CreateProductRoot(models.ProductRoot) (uint64, time.Time, error)
	UpdateProductRoot(models.ProductRoot) (time.Time, error)
	DeleteProductRoot(id uint64) (time.Time, error)

	// Products
	GetProduct(id uint64) (models.Product, error)
	CreateProduct(models.Product) (uint64, time.Time, time.Time, error)
	UpdateProduct(models.Product) (time.Time, error)
	DeleteProduct(id uint64) (time.Time, error)
	GetProductBySKU(string) (models.Product, error)

	// LoginAttempts
	GetLoginAttempt(id uint64) (models.LoginAttempt, error)
	CreateLoginAttempt(models.LoginAttempt) (uint64, time.Time, error)
	UpdateLoginAttempt(models.LoginAttempt) (time.Time, error)
	DeleteLoginAttempt(id uint64) (time.Time, error)

	// ProductOptions
	GetProductOption(id uint64) (models.ProductOption, error)
	CreateProductOption(models.ProductOption) (uint64, time.Time, error)
	UpdateProductOption(models.ProductOption) (time.Time, error)
	DeleteProductOption(id uint64) (time.Time, error)

	// ProductOptionValues
	GetProductOptionValue(id uint64) (models.ProductOptionValue, error)
	CreateProductOptionValue(models.ProductOptionValue) (uint64, time.Time, error)
	UpdateProductOptionValue(models.ProductOptionValue) (time.Time, error)
	DeleteProductOptionValue(id uint64) (time.Time, error)

	// ProductVariantBridge
	GetProductVariantBridge(id uint64) (models.ProductVariantBridge, error)
	CreateProductVariantBridge(models.ProductVariantBridge) (uint64, time.Time, error)
	UpdateProductVariantBridge(models.ProductVariantBridge) (time.Time, error)
	DeleteProductVariantBridge(id uint64) (time.Time, error)

	// PasswordResetTokens
	GetPasswordResetToken(id uint64) (models.PasswordResetToken, error)
	CreatePasswordResetToken(models.PasswordResetToken) (uint64, time.Time, error)
	UpdatePasswordResetToken(models.PasswordResetToken) (time.Time, error)
	DeletePasswordResetToken(id uint64) (time.Time, error)
}
