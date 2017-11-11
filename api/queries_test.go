package main

import (
	"testing"
	"time"

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
	exampleRoot := &ProductRoot{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
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

func TestBuildProductListQuery(t *testing.T) {
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
		 FROM products WHERE archived_on IS NULL LIMIT 25`
	actualQuery, actualArgs := buildProductListQuery(genereateDefaultQueryFilter())
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
		 FROM products WHERE archived_on IS NULL AND created_on > $1 AND created_on < $2 AND updated_on > $3 AND updated_on < $4 LIMIT 46 OFFSET 92`

	actualQuery, actualArgs := buildProductListQuery(queryFilter)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 4, len(actualArgs), argsEqualityErrorMessage)
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

func TestBuildProductUpdateQuery(t *testing.T) {
	t.Parallel()
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	expectedQuery := `UPDATE products SET cost = $1, name = $2, price = $3, quantity = $4, sku = $5, upc = $6, updated_on = NOW() WHERE id = $7 RETURNING *`
	actualQuery, actualArgs := buildProductUpdateQuery(exampleProduct)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 7, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductCreationQuery(t *testing.T) {
	t.Parallel()
	exampleProduct := &Product{
		DBRow: DBRow{
			ID:        2,
			CreatedOn: generateExampleTimeForTests(),
		},
		SKU:           "skateboard",
		Name:          "Skateboard",
		UPC:           "1234567890",
		Quantity:      123,
		Price:         99.99,
		Cost:          50.00,
		Description:   "This is a skateboard. Please wear a helmet.",
		ProductWeight: 8,
		ProductHeight: 7,
		ProductWidth:  6,
		ProductLength: 5,
		PackageWeight: 4,
		PackageHeight: 3,
		PackageWidth:  2,
		PackageLength: 1,
		AvailableOn:   generateExampleTimeForTests(),
	}

	expectedQuery := `INSERT INTO products (product_root_id,name,subtitle,description,option_summary,sku,upc,manufacturer,brand,quantity,taxable,price,on_sale,sale_price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length,quantity_per_package,available_on,updated_on) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,NOW()) RETURNING id, available_on, created_on`
	actualQuery, actualArgs := buildProductCreationQuery(exampleProduct)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 25, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildMultipleProductCreationQuery(t *testing.T) {
	t.Parallel()
	exampleProducts := []*Product{
		{
			DBRow: DBRow{
				ID:        2,
				CreatedOn: generateExampleTimeForTests(),
			},
			SKU:           "skateboard",
			Name:          "SKU ONE",
			UPC:           "1234567890",
			Quantity:      123,
			Price:         99.99,
			Cost:          50.00,
			Description:   "This is a skateboard. Please wear a helmet.",
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
			AvailableOn:   generateExampleTimeForTests(),
		},
		{
			DBRow: DBRow{
				ID:        2,
				CreatedOn: generateExampleTimeForTests(),
			},
			SKU:           "skateboard",
			Name:          "SKU TWO",
			UPC:           "1234567890",
			Quantity:      123,
			Price:         99.99,
			Cost:          50.00,
			Description:   "This is a skateboard. Please wear a helmet.",
			ProductWeight: 8,
			ProductHeight: 7,
			ProductWidth:  6,
			ProductLength: 5,
			PackageWeight: 4,
			PackageHeight: 3,
			PackageWidth:  2,
			PackageLength: 1,
			AvailableOn:   generateExampleTimeForTests(),
		},
	}

	expectedQuery := `INSERT INTO products (product_root_id,name,subtitle,description,option_summary,sku,upc,manufacturer,brand,quantity,taxable,price,on_sale,sale_price,cost,product_weight,product_height,product_width,product_length,package_weight,package_height,package_width,package_length,quantity_per_package,available_on,updated_on) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,NOW()),($25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,NOW()) RETURNING id, created_on`
	actualQuery, actualArgs := buildMultipleProductCreationQuery(exampleProducts)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 24*len(exampleProducts), len(actualArgs), argsEqualityErrorMessage)
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
	actualQuery, actualArgs := buildProductOptionListQuery(existingID, &QueryFilter{})

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_options SET name = $1, updated_on = NOW() WHERE id = $2 RETURNING updated_on`
	actualQuery, actualArgs := buildProductOptionUpdateQuery(&ProductOption{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_options (name,product_root_id) VALUES ($1,$2) RETURNING id, created_on`
	actualQuery, actualArgs := buildProductOptionCreationQuery(&ProductOption{}, exampleProductID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueUpdateQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `UPDATE product_option_values SET updated_on = NOW(), value = $1 WHERE id = $2 RETURNING updated_on`
	actualQuery, actualArgs := buildProductOptionValueUpdateQuery(&ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionValueCreationQuery(t *testing.T) {
	t.Parallel()
	expectedQuery := `INSERT INTO product_option_values (product_option_id,value) VALUES ($1,$2) RETURNING id, created_on`
	actualQuery, actualArgs := buildProductOptionValueCreationQuery(&ProductOptionValue{})
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildProductOptionCombinationExistenceQuery(t *testing.T) {
	t.Parallel()
	exampleData := []uint64{1, 4}
	expectedQuery := `SELECT EXISTS(SELECT id FROM product_variant_bridge WHERE product_option_value_id = $1 AND archived_on IS NULL) AND EXISTS(SELECT id FROM product_variant_bridge WHERE product_option_value_id = $2 AND archived_on IS NULL)`

	actualQuery, actualArgs := buildProductOptionCombinationExistenceQuery(exampleData)
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
	exampleDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	expectedQuery := `INSERT INTO discounts (name,discount_type,amount,starts_on,expires_on,requires_code,code,limited_use,number_of_uses,login_required) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, created_on`
	actualQuery, actualArgs := buildDiscountCreationQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 10, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildDiscountUpdateQuery(t *testing.T) {
	t.Parallel()
	exampleDiscount := &Discount{
		DBRow: DBRow{
			ID:        1,
			CreatedOn: generateExampleTimeForTests(),
		},
		Name:         "Example Discount",
		DiscountType: "flat_amount",
		Amount:       12.34,
		StartsOn:     generateExampleTimeForTests(),
		ExpiresOn:    NullTime{pq.NullTime{Time: generateExampleTimeForTests().Add(30 * (24 * time.Hour)), Valid: true}},
	}

	expectedQuery := `UPDATE discounts SET amount = $1, code = $2, discount_type = $3, expires_on = $4, limited_use = $5, login_required = $6, name = $7, number_of_uses = $8, requires_code = $9, starts_on = $10, updated_on = NOW() WHERE id = $11 RETURNING updated_on`
	actualQuery, actualArgs := buildDiscountUpdateQuery(exampleDiscount)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 11, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildUserSelectionQuery(t *testing.T) {
	t.Parallel()
	username := "frankzappa"
	expectedQuery := `SELECT id, first_name, last_name, username, email, password, salt, is_admin, password_last_changed_on, created_on, updated_on, archived_on FROM users WHERE username = $1 AND archived_on IS NULL`
	actualQuery, actualArgs := buildUserSelectionQuery(username)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildUserSelectionQueryByID(t *testing.T) {
	t.Parallel()
	userID := uint64(1)
	expectedQuery := `SELECT id, first_name, last_name, username, email, password, salt, is_admin, password_last_changed_on, created_on, updated_on, archived_on FROM users WHERE id = $1 AND archived_on IS NULL`
	actualQuery, actualArgs := buildUserSelectionQueryByID(userID)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 1, len(actualArgs), argsEqualityErrorMessage)
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

	expectedQuery := `INSERT INTO users (first_name,last_name,email,username,password,salt,is_admin) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_on`
	actualQuery, actualArgs := buildUserCreationQuery(user)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 7, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildPasswordResetRowCreationQuery(t *testing.T) {
	t.Parallel()
	userID := uint64(1)
	resetToken := "this_is_a_reset_token"
	expectedQuery := `INSERT INTO password_reset_tokens (user_id,token) VALUES ($1,$2)`
	actualQuery, actualArgs := buildPasswordResetRowCreationQuery(userID, resetToken)
	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildUserUpdateQueryWithoutPasswordChange(t *testing.T) {
	t.Parallel()

	user := &User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  "Password",
		Salt:      []byte("Salt"),
		IsAdmin:   true,
	}
	expectedQuery := `UPDATE users SET email = $1, first_name = $2, is_admin = $3, last_name = $4, updated_on = NOW(), username = $5 WHERE username = $6 RETURNING updated_on`
	actualQuery, actualArgs := buildUserUpdateQuery(user, false)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 6, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildUserUpdateQueryWithPasswordChange(t *testing.T) {
	t.Parallel()

	user := &User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "Email",
		Password:  "Password",
		Salt:      []byte("Salt"),
		IsAdmin:   true,
	}
	expectedQuery := `UPDATE users SET email = $1, first_name = $2, is_admin = $3, last_name = $4, password = $5, password_last_changed_on = NOW(), updated_on = NOW(), username = $6 WHERE username = $7 RETURNING updated_on`
	actualQuery, actualArgs := buildUserUpdateQuery(user, true)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 7, len(actualArgs), argsEqualityErrorMessage)
}

func TestBuildLoginAttemptCreationQuery(t *testing.T) {
	expectedQuery := `INSERT INTO login_attempts (username,successful) VALUES ($1,$2)`
	actualQuery, actualArgs := buildLoginAttemptCreationQuery("farts", true)

	assert.Equal(t, expectedQuery, actualQuery, queryEqualityErrorMessage)
	assert.Equal(t, 2, len(actualArgs), argsEqualityErrorMessage)
}
