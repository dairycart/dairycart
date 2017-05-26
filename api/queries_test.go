package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProgenitorRetrievalQuery(t *testing.T) {
	expected := `SELECT * FROM product_progenitors WHERE id = $1 AND archived_at IS NULL`
	actual := buildProgenitorRetrievalQuery(1)
	assert.Equal(t, expected, actual, "Generated SQL query should match expected SQL query")
}

func TestBuildRowExistenceQuery(t *testing.T) {
	expected := `SELECT EXISTS(SELECT 1 FROM things WHERE stuff = $1 AND archived_at IS NULL)`
	actual := buildRowExistenceQuery("things", "stuff", "abritrary")
	assert.Equal(t, expected, actual, "Generated SQL query should match expected SQL query")
}

func TestBuildProgenitorExistenceQuery(t *testing.T) {
	expected := `SELECT EXISTS(SELECT 1 FROM product_progenitors WHERE id = $1 AND archived_at IS NULL)`
	actual := buildProgenitorExistenceQuery(1)
	assert.Equal(t, expected, actual, "Generated SQL query should match expected SQL query")
}

func TestBuildProgenitorCreationQuery(t *testing.T) {
	expected := `INSERT INTO product_progenitors (name,description,taxable,price,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW()) RETURNING "id"`
	actual := buildProgenitorCreationQuery(exampleProgenitor)
	assert.Equal(t, expected, actual, "Generated SQL query should match expected SQL query")
}
