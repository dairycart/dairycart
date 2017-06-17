package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	aTimestamp                = 232747200
	anOlderTimestamp          = aTimestamp + 10000
	existingID                = 1
	existingIDString          = "1"
	queryEqualityErrorMessage = "Generated SQL query should match expected SQL query"
	argsEqualityErrorMessage  = "Generated SQL arguments should match expected arguments"
)

// Note: comparing interface equality with assert is impossible as far as I can tell,
// so generally these tests ensure that the correct number of args are returned.

func TestBuildProgenitorCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_progenitors (name,description,taxable,price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`
	actualQuery, actualArgs := buildProgenitorCreationQuery(exampleProgenitor)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 13, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT p.id as product_id,
				p.product_progenitor_id,
				p.sku,
				p.name as product_name,
				p.upc,
				p.quantity,
				p.price as product_price,
				p.cost as product_cost,
				p.created_on as product_created_on,
				p.updated_on as product_updated_on,
				p.archived_on as product_archived_on,
				g.* FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductListQuery(defaultQueryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 0, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQueryAndPartiallyCustomQueryFilter(t *testing.T) {
	t.Parallel()
	queryFilter := &QueryFilter{
		Page:          3,
		Limit:         25,
		UpdatedBefore: time.Unix(int64(aTimestamp), 0),
		UpdatedAfter:  time.Unix(int64(anOlderTimestamp), 0),
	}

	expectedQuery := `SELECT p.id as product_id,
				p.product_progenitor_id,
				p.sku,
				p.name as product_name,
				p.upc,
				p.quantity,
				p.price as product_price,
				p.cost as product_cost,
				p.created_on as product_created_on,
				p.updated_on as product_updated_on,
				p.archived_on as product_archived_on,
				g.* FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_on IS NULL AND p.updated_on > $1 AND p.updated_on < $2 LIMIT 25 OFFSET 50`

	actualQuery, actualArgs := buildProductListQuery(queryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQueryAndCompletelyCustomQueryFilter(t *testing.T) {
	t.Parallel()
	queryFilter := &QueryFilter{
		Page:          3,
		Limit:         46,
		UpdatedBefore: time.Unix(int64(aTimestamp), 0),
		UpdatedAfter:  time.Unix(int64(anOlderTimestamp), 0),
		CreatedBefore: time.Unix(int64(aTimestamp), 0),
		CreatedAfter:  time.Unix(int64(anOlderTimestamp), 0),
	}

	expectedQuery := `SELECT p.id as product_id,
				p.product_progenitor_id,
				p.sku,
				p.name as product_name,
				p.upc,
				p.quantity,
				p.price as product_price,
				p.cost as product_cost,
				p.created_on as product_created_on,
				p.updated_on as product_updated_on,
				p.archived_on as product_archived_on,
				g.* FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_on IS NULL AND p.created_on > $1 AND p.created_on < $2 AND p.updated_on > $3 AND p.updated_on < $4 LIMIT 46 OFFSET 92`

	actualQuery, actualArgs := buildProductListQuery(queryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 4, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE products SET cost = $1, name = $2, price = $3, quantity = $4, sku = $5, upc = $6, updated_on = NOW() WHERE id = $7 RETURNING *`
	actualQuery, actualArgs := buildProductUpdateQuery(exampleProduct)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 7, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO products (product_progenitor_id,sku,name,upc,quantity,price,cost) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	actualQuery, actualArgs := buildProductCreationQuery(exampleProduct)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 7, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT count(id) over (), * FROM product_options WHERE product_progenitor_id = $1 AND archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductOptionListQuery(existingIDString, &QueryFilter{})

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_options SET name = $1, updated_on = NOW() WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductOptionUpdateQuery(&ProductOption{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_options (name,product_progenitor_id) VALUES ($1,$2) RETURNING "id"`
	actualQuery, actualArgs := buildProductOptionCreationQuery(&ProductOption{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_option_values SET updated_on = NOW(), value = $1 WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductOptionValueUpdateQuery(&ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_option_values (product_option_id,value) VALUES ($1,$2) RETURNING "id"`
	actualQuery, actualArgs := buildProductOptionValueCreationQuery(&ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT id, name, type, amount, starts_on, expires_on, requires_code, code, limited_use, number_of_uses, login_required, created_on, updated_on, archived_on FROM discounts WHERE (expires_on IS NULL OR expires_on > $1) AND archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildDiscountListQuery(defaultQueryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO discounts (name,type,amount,starts_on,expires_on,requires_code,code,limited_use,number_of_uses,login_required) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, name, type, amount, starts_on, expires_on, requires_code, code, limited_use, number_of_uses, login_required, created_on, updated_on, archived_on`
	actualQuery, actualArgs := buildDiscountCreationQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 10, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE discounts SET amount = $1, code = $2, expires_on = $3, limited_use = $4, login_required = $5, name = $6, number_of_uses = $7, requires_code = $8, starts_on = $9, type = $10, updated_on = NOW() WHERE id = $11 RETURNING id, name, type, amount, starts_on, expires_on, requires_code, code, limited_use, number_of_uses, login_required, created_on, updated_on, archived_on`
	actualQuery, actualArgs := buildDiscountUpdateQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 11, len(actualArgs), argsEqualityErrorMessage)
}
