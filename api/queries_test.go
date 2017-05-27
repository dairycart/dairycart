package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	queryEqualityErrorMessage = "Generated SQL query should match expected SQL query"
	argsEqualityErrorMessage  = "Generated SQL arguments should match expected arguments"
)

func TestBuildRowExistenceQuery(t *testing.T) {
	expected := `SELECT EXISTS(SELECT 1 FROM things WHERE stuff = $1 AND archived_at IS NULL)`
	actual := buildRowExistenceQuery("things", "stuff", "abritrary")
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildRowRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM things WHERE stuff = $1 AND archived_at IS NULL`
	actual := buildRowRetrievalQuery("things", "stuff", "abritrary")
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildRowDeletionQuery(t *testing.T) {
	expected := `UPDATE things SET archived_at = NOW() WHERE stuff = $1 AND archived_at IS NULL`
	actual := buildRowDeletionQuery("things", "stuff", "abritrary")
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProgenitorRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM product_progenitors WHERE id = $1 AND archived_at IS NULL`
	actual := buildProgenitorRetrievalQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProgenitorExistenceQuery(t *testing.T) {
	expected := `SELECT EXISTS(SELECT 1 FROM product_progenitors WHERE id = $1 AND archived_at IS NULL)`
	actual := buildProgenitorExistenceQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProgenitorCreationQuery(t *testing.T) {
	expectedQuery := `INSERT INTO product_progenitors (name,description,taxable,price,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`
	actualQuery, actualArgs := buildProgenitorCreationQuery(exampleProgenitor)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	// comparing interface equality with assert is impossible as far as I can tell
	assert.Equal(t, 12, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductExistenceQuery(t *testing.T) {
	expected := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND archived_at IS NULL)`
	actual := buildProductExistenceQuery(exampleSKU)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}
func TestBuildProductRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM products WHERE sku = $1 AND archived_at IS NULL`
	actual := buildProductRetrievalQuery(exampleSKU)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildAllProductsRetrievalQuery(t *testing.T) {
	expectedQuery := `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.archived_at IS NULL LIMIT 25`
	actualQuery, actualArgs := buildAllProductsRetrievalQuery(defaultQueryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 0, len(actualArgs), queryEqualityErrorMessage)
}
func TestBuildProductDeletionQuery(t *testing.T) {
	expected := `UPDATE products SET archived_at = NOW() WHERE sku = $1 AND archived_at IS NULL`
	actual := buildProductDeletionQuery(exampleSKU)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}
func TestBuildCompleteProductRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM products p JOIN product_progenitors g ON p.product_progenitor_id = g.id WHERE p.sku = $1 AND p.archived_at IS NULL`
	actual := buildCompleteProductRetrievalQuery(exampleSKU)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}
func TestBuildProductUpdateQuery(t *testing.T) {
	expectedQuery := `UPDATE products SET name = $1, on_sale = $2, price = $3, quantity = $4, sale_price = $5, sku = $6, upc = $7, updated_at = NOW() WHERE id = $8 RETURNING *`
	actualQuery, actualArgs := buildProductUpdateQuery(exampleProduct)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	// comparing interface equality with assert is impossible as far as I can tell
	assert.Equal(t, 8, len(actualArgs), argsEqualityErrorMessage)
}
func TestBuildProductCreationQuery(t *testing.T) {
	expected := `INSERT INTO products (product_progenitor_id,sku,name,upc,quantity,on_sale,price,sale_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *`
	actual, _ := buildProductCreationQuery(exampleProduct)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM product_attributes WHERE id = $1 AND archived_at IS NULL`
	actual := buildProductAttributeRetrievalQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeListQuery(t *testing.T) {
	expected := `SELECT * FROM product_attributes WHERE product_progenitor_id = $1 AND archived_at IS NULL`
	actual := buildProductAttributeListQuery("1")
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeDeletionQuery(t *testing.T) {
	expected := `UPDATE product_attributes SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`
	actual := buildProductAttributeDeletionQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeUpdateQuery(t *testing.T) {
	expectedQuery := `UPDATE product_attributes SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductAttributeUpdateQuery(&ProductAttribute{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeCreationQuery(t *testing.T) {
	expectedQuery := `INSERT INTO product_attributes (name,product_progenitor_id) VALUES ($1,$2) RETURNING *`
	actualQuery, actualArgs := buildProductAttributeCreationQuery(&ProductAttribute{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	// comparing interface equality with assert is impossible as far as I can tell
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeValueRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM product_attribute_values WHERE id = $1 AND archived_at IS NULL`
	actual := buildProductAttributeValueRetrievalQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeValueDeletionQuery(t *testing.T) {
	expected := `UPDATE product_attribute_values SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`
	actual := buildProductAttributeValueDeletionQuery(1)
	assert.Equal(t, expected, actual, queryEqualityErrorMessage)
}

func TestBuildProductAttributeValueUpdateQuery(t *testing.T) {
	expectedQuery := `UPDATE product_attribute_values SET updated_at = NOW(), value = $1 WHERE id = $2 RETURNING *`
	actualQuery, actualArgs := buildProductAttributeValueUpdateQuery(&ProductAttributeValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductAttributeValueCreationQuery(t *testing.T) {
	expectedQuery := `INSERT INTO product_attribute_values (product_attribute_id,value) VALUES ($1,$2) RETURNING *`
	actualQuery, actualArgs := buildProductAttributeValueCreationQuery(&ProductAttributeValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	// comparing interface equality with assert is impossible as far as I can tell
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}
