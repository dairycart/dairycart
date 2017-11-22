package main

import (
	"fmt"
	"strings"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/Masterminds/squirrel"
)

func applyQueryFilterToQueryBuilder(queryBuilder squirrel.SelectBuilder, queryFilter *models.QueryFilter, includeOffset bool) squirrel.SelectBuilder {
	if queryFilter == nil {
		return queryBuilder
	}

	if queryFilter.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(queryFilter.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(25)
	}

	if queryFilter.Page > 1 && includeOffset {
		offset := (queryFilter.Page - 1) * uint64(queryFilter.Limit)
		queryBuilder = queryBuilder.Offset(offset)
	}

	if !queryFilter.CreatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"created_on": queryFilter.CreatedAfter})
	}

	if !queryFilter.CreatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"created_on": queryFilter.CreatedBefore})
	}

	if !queryFilter.UpdatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"updated_on": queryFilter.UpdatedAfter})
	}

	if !queryFilter.UpdatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"updated_on": queryFilter.UpdatedBefore})
	}
	return queryBuilder
}

func buildCountQuery(table string, queryFilter *models.QueryFilter) string {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select("count(id)").
		From(table).
		Where(squirrel.Eq{"archived_on": nil})

	// setting this to false so we always get a count
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, false)

	query, _, _ := queryBuilder.ToSql()
	return query
}

////////////////////////////////////////////////////////
//                                                    //
//                   Product Roots                    //
//                                                    //
////////////////////////////////////////////////////////

func buildProductRootListQuery(queryFilter *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		// note, this has to look ugly and disjointed because otherwise my editor
		// will delete the trailing space and the tests will fail. Womp womp.
		Select(`id,
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
		`).
		From("product_roots").
		Where(squirrel.Eq{"archived_on": nil}).
		Limit(uint64(queryFilter.Limit))

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductRootCreationQuery(r *models.ProductRoot) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_roots").
		Columns(
			"name",
			"subtitle",
			"description",
			"sku_prefix",
			"manufacturer",
			"brand",
			"available_on",
			"quantity_per_package",
			"taxable",
			"cost",
			"product_weight",
			"product_height",
			"product_width",
			"product_length",
			"package_weight",
			"package_height",
			"package_width",
			"package_length",
		).
		Values(
			r.Name,
			r.Subtitle,
			r.Description,
			r.SKUPrefix,
			r.Manufacturer,
			r.Brand,
			r.AvailableOn,
			r.QuantityPerPackage,
			r.Taxable,
			r.Cost,
			r.ProductWeight,
			r.ProductHeight,
			r.ProductWidth,
			r.ProductLength,
			r.PackageWeight,
			r.PackageHeight,
			r.PackageWidth,
			r.PackageLength,
		).
		Suffix(`RETURNING id, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                     Products                       //
//                                                    //
////////////////////////////////////////////////////////

func getProductCreationColumns() []string {
	c := []string{
		"product_root_id",
		"name",
		"subtitle",
		"description",
		"option_summary",
		"sku",
		"upc",
		"manufacturer",
		"brand",
		"quantity",
		"taxable",
		"price",
		"on_sale",
		"sale_price",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"quantity_per_package",
		"available_on",
		"updated_on",
	}
	return c
}

func buildProductListQuery(queryFilter *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		// note, this has to look ugly and disjointed because otherwise my editor
		// will delete the trailing space and the tests will fail. Womp womp.
		Select(`id,
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
		`).
		From("products").
		Where(squirrel.Eq{"archived_on": nil})

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductAssociatedWithRootListQuery(rootID uint64) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		// note, this has to look ugly and disjointed because otherwise my editor
		// will delete the trailing space and the tests will fail. Womp womp.
		Select(`id,
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
		`).
		From("products").
		Where(squirrel.Eq{"archived_on": nil}).
		Where(squirrel.Eq{"product_root_id": rootID})

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductUpdateQuery(p *models.Product) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	productUpdateSetMap := map[string]interface{}{
		"sku":        p.SKU,
		"name":       p.Name,
		"upc":        p.UPC,
		"quantity":   p.Quantity,
		"price":      p.Price,
		"cost":       p.Cost,
		"updated_on": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("products").
		SetMap(productUpdateSetMap).
		Where(squirrel.Eq{"id": p.ID}).
		Suffix(`RETURNING *`)
		// Suffix(`RETURNING updated_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductCreationQuery(p *models.Product) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	columns := getProductCreationColumns()

	values := []interface{}{
		p.ProductRootID,
		p.Name,
		p.Subtitle,
		p.Description,
		p.OptionSummary,
		p.SKU,
		p.UPC,
		p.Manufacturer,
		p.Brand,
		p.Quantity,
		p.Taxable,
		p.Price,
		p.OnSale,
		p.SalePrice,
		p.Cost,
		p.ProductWeight,
		p.ProductHeight,
		p.ProductWidth,
		p.ProductLength,
		p.PackageWeight,
		p.PackageHeight,
		p.PackageWidth,
		p.PackageLength,
		p.QuantityPerPackage,
		p.AvailableOn,
		squirrel.Expr("NOW()"),
	}

	queryBuilder := sqlBuilder.
		Insert("products").
		Columns(columns...).
		Values(values...).
		Suffix(`RETURNING id, available_on, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildMultipleProductCreationQuery(ps []*models.Product) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("products").
		Columns(getProductCreationColumns()...)

	for _, p := range ps {
		values := []interface{}{
			p.ProductRootID,
			p.Name,
			p.Subtitle,
			p.Description,
			p.OptionSummary,
			p.SKU,
			p.Manufacturer,
			p.Brand,
			p.Quantity,
			p.Taxable,
			p.Price,
			p.OnSale,
			p.SalePrice,
			p.Cost,
			p.ProductWeight,
			p.ProductHeight,
			p.ProductWidth,
			p.ProductLength,
			p.PackageWeight,
			p.PackageHeight,
			p.PackageWidth,
			p.PackageLength,
			p.QuantityPerPackage,
			p.AvailableOn,
			squirrel.Expr("NOW()"),
		}

		queryBuilder = queryBuilder.Values(values...)
	}
	queryBuilder = queryBuilder.Suffix(`RETURNING id, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Options                    //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionListQuery(productRootID uint64, queryFilter *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(productOptionsHeaders).
		From("product_options").
		Where(squirrel.Eq{"product_root_id": productRootID}).
		Where(squirrel.Eq{"archived_on": nil})
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionUpdateQuery(a *models.ProductOption) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	productOptionUpdateSetMap := map[string]interface{}{
		"name":       a.Name,
		"updated_on": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_options").
		SetMap(productOptionUpdateSetMap).
		Where(squirrel.Eq{"id": a.ID}).
		Suffix(`RETURNING updated_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionCreationQuery(a *models.ProductOption, productRootID uint64) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_options").
		Columns("name", "product_root_id").
		Values(a.Name, productRootID).
		Suffix(`RETURNING id, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//               Product Option Values                //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionValueUpdateQuery(v *models.ProductOptionValue) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	productOptionUpdateSetMap := map[string]interface{}{
		"value":      v.Value,
		"updated_on": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_option_values").
		SetMap(productOptionUpdateSetMap).
		Where(squirrel.Eq{"id": v.ID}).
		Suffix(`RETURNING updated_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionValueCreationQuery(v *models.ProductOptionValue) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_option_values").
		Columns("product_option_id", "value").
		Values(v.ProductOptionID, v.Value).
		Suffix(`RETURNING id, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionCombinationExistenceQuery(optionValueIDs []uint64) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder
	var subqueries []string
	var subargs []interface{}
	for _, id := range optionValueIDs {
		q, a, _ := sqlBuilder.
			Select("id").
			From("product_variant_bridge").
			Where(squirrel.Eq{"product_option_value_id": id}).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		subqueries = append(subqueries, fmt.Sprintf("(%s)", q))
		subargs = append(subargs, a)
	}
	prequery := strings.Join(subqueries, " AND EXISTS")

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS%s", prequery))
	query, _, _ := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	return query, subargs
}

////////////////////////////////////////////////////////
//                                                    //
//                     Discounts                      //
//                                                    //
////////////////////////////////////////////////////////

func buildDiscountListQuery(queryFilter *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(`id,
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
			archived_on`).
		From("discounts").
		Where(squirrel.Or{squirrel.Eq{"expires_on": nil}, squirrel.Gt{"expires_on": "NOW()"}}).
		Where(squirrel.Eq{"archived_on": nil})

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildDiscountCreationQuery(d *models.Discount) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("discounts").
		Columns("name", "discount_type", "amount", "starts_on", "expires_on", "requires_code", "code", "limited_use", "number_of_uses", "login_required").
		Values(d.Name, d.DiscountType, d.Amount, d.StartsOn, d.ExpiresOn, d.RequiresCode, d.Code, d.LimitedUse, d.NumberOfUses, d.LoginRequired).
		Suffix("RETURNING id, created_on")
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildDiscountUpdateQuery(d *models.Discount) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	updateSetMap := map[string]interface{}{
		"name":           d.Name,
		"discount_type":  d.DiscountType,
		"amount":         d.Amount,
		"starts_on":      d.StartsOn,
		"expires_on":     d.ExpiresOn.Time,
		"requires_code":  d.RequiresCode,
		"code":           d.Code,
		"limited_use":    d.LimitedUse,
		"number_of_uses": d.NumberOfUses,
		"login_required": d.LoginRequired,
		"updated_on":     squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("discounts").
		SetMap(updateSetMap).
		Where(squirrel.Eq{"id": d.ID}).
		Suffix("RETURNING updated_on")
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                       Auth                         //
//                                                    //
////////////////////////////////////////////////////////

func buildUserSelectionQuery(username string) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(usersTableHeaders).
		From("users").
		Where(squirrel.Eq{"username": username}).
		Where(squirrel.Eq{"archived_on": nil})

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildUserSelectionQueryByID(userID uint64) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(usersTableHeaders).
		From("users").
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"archived_on": nil})

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildUserCreationQuery(u *models.User) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("users").
		Columns(
			"first_name",
			"last_name",
			"email",
			"username",
			"password",
			"salt",
			"is_admin",
		).
		Values(
			u.FirstName,
			u.LastName,
			u.Email,
			u.Username,
			u.Password,
			u.Salt,
			u.IsAdmin,
		).
		Suffix(`RETURNING id, created_on`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildPasswordResetRowCreationQuery(userID uint64, resetToken string) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("password_reset_tokens").
		Columns(
			"user_id",
			"token",
		).
		Values(
			userID,
			resetToken,
		)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildUserUpdateQuery(u *models.User, passwordChanged bool) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	updateSetMap := map[string]interface{}{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"username":   u.Username,
		"email":      u.Email,
		"is_admin":   u.IsAdmin,
		"updated_on": squirrel.Expr("NOW()"),
	}
	if passwordChanged {
		updateSetMap["password"] = u.Password
		updateSetMap["password_last_changed_on"] = squirrel.Expr("NOW()")
	}

	queryBuilder := sqlBuilder.
		Update("users").
		SetMap(updateSetMap).
		Where(squirrel.Eq{"username": u.Username}).
		Suffix("RETURNING updated_on")
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildLoginAttemptCreationQuery(username string, successful bool) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("login_attempts").
		Columns(
			"username",
			"successful",
		).
		Values(
			username,
			successful,
		)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}
