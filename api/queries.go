package main

import (
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
