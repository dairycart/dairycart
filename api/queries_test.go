// +build !migrated

package main

import (
	"testing"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	aTimestamp       = 232747200
	anOlderTimestamp = aTimestamp + 10000
	existingID       = 1

	queryEqualityErrorMessage = "Generated SQL query should match expected SQL query"
	argsEqualityErrorMessage  = "Generated SQL arguments should match expected arguments"
)

// Note: comparing interface equality with assert is impossible as far as I can tell,
// so generally these tests ensure that the correct number of args are returned.

func TestBuildProductRootCreationQuery(t *testing.T) {
	t.Parallel()
	exampleRoot := &models.ProductRoot{
		ID:            2,
		CreatedOn:     generateExampleTimeForTests(),
		Name:          "Skateboard",
		Description:   "This is a skateboard. Please wear a helmet.",
		Cost:          50.00,
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
	}
	expectedQuery := `INSERT INTO product_roots (name,subtitle,description,sku_prefix,manufacturer,brand,available_on,quantity_per_package,taxable,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING id, created_on`
	actualQuery, actualArgs := buildProductRootCreationQuery(exampleRoot)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 18, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductRootListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT id,
			name,
			subtitle,
			description,
			sku_prefix,
			manufacturer,
			brand,
			taxable,
			cost,
			product_weight,
			product_height,
			product_width,
			product_length,
			package_weight,
			package_height,
			package_width,
			package_length,
			quantity_per_package,
			available_on,
			created_on,
			updated_on,
			archived_on
		 FROM product_roots WHERE archived_on IS NULL LIMIT 25`
	actualQuery, _ := buildProductRootListQuery(genereateDefaultQueryFilter())
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAssociatedWithRootListQuery(t *testing.T) {
	t.Parallel()

	expectedQuery := `SELECT id,
			product_root_id,
			name,
			subtitle,
			description,
			option_summary,
			sku,
			upc,
			manufacturer,
			brand,
			quantity,
			taxable,
			price,
			on_sale,
			sale_price,
			cost,
			product_weight,
			product_height,
			product_width,
			product_length,
			package_weight,
			package_height,
			package_width,
			package_length,
			quantity_per_package,
			available_on,
			created_on,
			updated_on,
			archived_on
		 FROM products WHERE archived_on IS NULL AND product_root_id = $1`
	actualQuery, actualArgs := buildProductAssociatedWithRootListQuery(123)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT id,
		name,
		product_root_id,
		created_on,
		updated_on,
		archived_on
	 FROM product_options WHERE product_root_id = $1 AND archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductOptionListQuery(existingID, &models.QueryFilter{})

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_options SET name = $1, updated_on = NOW() WHERE id = $2 RETURNING updated_on`
	actualQuery, actualArgs := buildProductOptionUpdateQuery(&models.ProductOption{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_options (name,product_root_id) VALUES ($1,$2) RETURNING id, created_on`
	actualQuery, actualArgs := buildProductOptionCreationQuery(&models.ProductOption{}, exampleProductID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_option_values SET updated_on = NOW(), value = $1 WHERE id = $2 RETURNING updated_on`
	actualQuery, actualArgs := buildProductOptionValueUpdateQuery(&models.ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_option_values (product_option_id,value) VALUES ($1,$2) RETURNING id, created_on`
	actualQuery, actualArgs := buildProductOptionValueCreationQuery(&models.ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT id,
			name,
			discount_type,
			amount,
			starts_on,
			expires_on,
			requires_code,
			code,
			limited_use,
			number_of_uses,
			login_required,
			created_on,
			updated_on,
			archived_on FROM discounts WHERE (expires_on IS NULL OR expires_on > $1) AND archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildDiscountListQuery(genereateDefaultQueryFilter())
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountCreationQuery(t *testing.T) {
	t.Parallel()
	exampleDiscount := &models.Discount{
		ID:           1,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    models.NullTime{NullTime: pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	expectedQuery := `INSERT INTO discounts (name,discount_type,amount,starts_on,expires_on,requires_code,code,limited_use,number_of_uses,login_required) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, created_on`
	actualQuery, actualArgs := buildDiscountCreationQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 10, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountUpdateQuery(t *testing.T) {
	t.Parallel()
	exampleDiscount := &models.Discount{
		ID:           1,
		CreatedOn:    generateExampleTimeForTests(),
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    models.NullTime{NullTime: pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	expectedQuery := `UPDATE discounts SET amount = $1, code = $2, discount_type = $3, expires_on = $4, limited_use = $5, login_required = $6, name = $7, number_of_uses = $8, requires_code = $9, starts_on = $10, updated_on = NOW() WHERE id = $11 RETURNING updated_on`
	actualQuery, actualArgs := buildDiscountUpdateQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 11, len(actualArgs), argsEqualityErrorMessage)
}
