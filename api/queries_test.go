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

func TestBuildProgenitorRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM product_progenitors WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildProgenitorRetrievalQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProgenitorExistenceQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM product_progenitors WHERE id = $1 AND archived_at IS NULL)`
	actualQuery := buildProgenitorExistenceQuery(existingIDString)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProgenitorCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_progenitors (name,description,taxable,price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`
	actualQuery, actualArgs := buildProgenitorCreationQuery(exampleProgenitor)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 13, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductExistenceQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND archived_at IS NULL)`
	actualQuery := buildProductExistenceQuery(exampleSKU)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}
func TestBuildProductRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM products WHERE sku = $1 AND archived_at IS NULL`
	actualQuery := buildProductRetrievalQuery(exampleSKU)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
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

func TestBuildProductDeletionQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE products SET archived_at = NOW() WHERE sku = $1 AND archived_at IS NULL`
	actualQuery := buildProductDeletionQuery(exampleSKU)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
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

func TestBuildProductAttributeRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM product_attributes WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildProductAttributeRetrievalQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAttributeListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT count(id) over (), * FROM product_attributes WHERE product_progenitor_id = $1 AND archived_at IS NULL LIMIT 25`
	actualQuery := buildProductAttributeListQuery(existingIDString, &QueryFilter{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAttributeDeletionQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_attributes SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildProductAttributeDeletionQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAttributeUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_attributes SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductAttributeUpdateQuery(&ProductAttribute{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_attributes (name,product_progenitor_id) VALUES ($1,$2) RETURNING "id"`
	actualQuery, actualArgs := buildProductAttributeCreationQuery(&ProductAttribute{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeValueRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM product_attribute_values WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildProductAttributeValueRetrievalQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAttributeValueDeletionQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_attribute_values SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildProductAttributeValueDeletionQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildProductAttributeValueExistenceForAttributeIDQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM product_attribute_values WHERE product_attribute_id = $1 AND value = $2 AND archived_at IS NULL)`
	actualQuery, actualArgs := buildProductAttributeValueExistenceForAttributeIDQuery(1, "value")
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeValueUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_attribute_values SET updated_at = NOW(), value = $1 WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductAttributeValueUpdateQuery(&ProductAttributeValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeValueCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_attribute_values (product_attribute_id,value) VALUES ($1,$2) RETURNING "id"`
	actualQuery, actualArgs := buildProductAttributeValueCreationQuery(&ProductAttributeValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountExistenceQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT EXISTS(SELECT 1 FROM discounts WHERE id = $1 AND archived_at IS NULL)`
	actualQuery := buildDiscountExistenceQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildDiscountRetrievalQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT * FROM discounts WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildDiscountRetrievalQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}

func TestBuildDiscountDeletionQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE discounts SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`
	actualQuery := buildDiscountDeletionQuery(existingID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
}
