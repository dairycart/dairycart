package main

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

func TestBuildProductListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT id,
		name,
		subtitle,
		description,
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
	 FROM products WHERE archived_on IS NULL LIMIT 25`
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

	expectedQuery := `SELECT id,
		name,
		subtitle,
		description,
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
	 FROM products WHERE archived_on IS NULL AND updated_on > $1 AND updated_on < $2 LIMIT 25 OFFSET 50`

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

	expectedQuery := `SELECT id,
		name,
		subtitle,
		description,
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
	 FROM products WHERE archived_on IS NULL AND created_on > $1 AND created_on < $2 AND updated_on > $3 AND updated_on < $4 LIMIT 46 OFFSET 92`

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
	expectedQuery := `INSERT INTO products (name,subtitle,description,sku,upc,manufacturer,brand,quantity,taxable,price,on_sale,sale_price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length,quantity_per_package,available_on,updated_on) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,NOW()) RETURNING "id"`
	actualQuery, actualArgs := buildProductCreationQuery(exampleProduct)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 23, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionListQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `SELECT count(id) over (), id,
		name,
		product_id,
		created_on,
		updated_on,
		archived_on
	 FROM product_options WHERE product_id = $1 AND archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductOptionListQuery(existingID, &QueryFilter{})

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
	expectedQuery := `INSERT INTO product_options (name,product_id) VALUES ($1,$2) RETURNING "id"`
	actualQuery, actualArgs := buildProductOptionCreationQuery(&ProductOption{}, exampleProduct.ID)
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
	expectedQuery := "SELECT \n\t\tid,\n\t\tname,\n\t\ttype,\n\t\tamount,\n\t\tstarts_on,\n\t\texpires_on,\n\t\trequires_code,\n\t\tcode,\n\t\tlimited_use,\n\t\tnumber_of_uses,\n\t\tlogin_required,\n\t\tcreated_on,\n\t\tupdated_on,\n\t\tarchived_on\n\t FROM discounts WHERE (expires_on IS NULL OR expires_on > $1) AND archived_on IS NULL LIMIT 25"
	actualQuery, actualArgs := buildDiscountListQuery(defaultQueryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO discounts (name,type,amount,starts_on,expires_on,requires_code,code,limited_use,number_of_uses,login_required) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING
		id,
		name,
		type,
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
		archived_on
	`
	actualQuery, actualArgs := buildDiscountCreationQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 10, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE discounts SET amount = $1, code = $2, expires_on = $3, limited_use = $4, login_required = $5, name = $6, number_of_uses = $7, requires_code = $8, starts_on = $9, type = $10, updated_on = NOW() WHERE id = $11 RETURNING
		id,
		name,
		type,
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
		archived_on
	`
	actualQuery, actualArgs := buildDiscountUpdateQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 11, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildUserCreationQuery(t *testing.T) {
	t.Parallel()
	user := &User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  "Password",
		Salt:      []byte("Salt"),
		IsAdmin:   true,
	}

	expectedQuery := `INSERT INTO users (first_name,last_name,email,password,salt,is_admin) VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	actualQuery, actualArgs := buildUserCreationQuery(user)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 6, len(actualArgs), argsEqualityErrorMessage)
}
