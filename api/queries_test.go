package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	existingID                = 1
	existingIDString          = "1"
	queryEqualityErrorMessage = "Generated SQL query should match expected SQL query"
	argsEqualityErrorMessage  = "Generated SQL arguments should match expected arguments"
)

// Note: comparing interface equality with assert is impossible as far as I can tell,
// so generally these tests ensure that the correct number of args are returned.

func TestBuildRowExistenceQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM things WHERE stuff = $1 AND archived_at IS NULL)`
	actualQuery := buildRowExistenceQuery("things", "stuff", "abritrary")
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildRowRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM things WHERE stuff = $1 AND archived_at IS NULL`
	actualQuery := buildRowRetrievalQuery("things", "stuff", "abritrary")
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildRowDeletionQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE things SET archived_at = NOW() WHERE stuff = $1 AND archived_at IS NULL`
	actualQuery := buildRowDeletionQuery("things", "stuff", "abritrary")
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProgenitorCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_progenitors (name,description,taxable,price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`
	actualQuery, actualArgs := buildProgenitorCreationQuery(exampleProgenitor)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 13, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT count(p.id) over (), * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_at IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductListQuery(defaultQueryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 0, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQueryAndPartiallyCustomQueryFilter(t *testing.T) {
	t.Parallel()
	queryFilter := &QueryFilter{
		Page:          3,
		Limit:         25,
		UpdatedBefore: time.Unix(int64(232747200), 0),
		UpdatedAfter:  time.Unix(int64(232747200+10000), 0),
	}
	expectedQuery := `SELECT count(p.id) over (), * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_at IS NULL AND p.updated_at > $1 AND p.updated_at < $2 LIMIT 25 OFFSET 50`
	actualQuery, actualArgs := buildProductListQuery(queryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductListQueryAndCompletelyCustomQueryFilter(t *testing.T) {
	t.Parallel()
	queryFilter := &QueryFilter{
		Page:          3,
		Limit:         46,
		UpdatedBefore: time.Unix(int64(232747200), 0),
		UpdatedAfter:  time.Unix(int64(232747200+10000), 0),
		CreatedBefore: time.Unix(int64(232747200), 0),
		CreatedAfter:  time.Unix(int64(232747200+10000), 0),
	}
	expectedQuery := `SELECT count(p.id) over (), * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_at IS NULL AND p.created_at > $1 AND p.created_at < $2 AND p.updated_at > $3 AND p.updated_at < $4 LIMIT 46 OFFSET 92`
	actualQuery, actualArgs := buildProductListQuery(queryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 4, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildCompleteProductRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.sku = $1 AND p.archived_at IS NULL`
	actualQuery := buildCompleteProductRetrievalQuery(exampleSKU)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE products SET cost = $1, name = $2, price = $3, quantity = $4, sku = $5, upc = $6, updated_at = NOW() WHERE id = $7 RETURNING *`
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
	expectedQuery := `SELECT count(id) over (), * FROM product_options WHERE product_progenitor_id = $1 AND archived_at IS NULL LIMIT 25`
	actualQuery := buildProductOptionListQuery(existingIDString, &QueryFilter{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductOptionUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_options SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING *`
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

func TestBuildProductOptionValueExistenceForOptionIDQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM product_option_values WHERE product_option_id = $1 AND value = $2 AND archived_at IS NULL)`
	actualQuery, actualArgs := buildProductOptionValueExistenceForOptionIDQuery(1, "value")
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_option_values SET updated_at = NOW(), value = $1 WHERE id = $2 RETURNING *`
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
