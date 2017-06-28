package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func buildCountQuery(table string, queryFilter *QueryFilter) string {
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

func applyQueryFilterToQueryBuilder(queryBuilder squirrel.SelectBuilder, queryFilter *QueryFilter, includeOffset bool) squirrel.SelectBuilder {
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

////////////////////////////////////////////////////////
//                                                    //
//                     Products                       //
//                                                    //
////////////////////////////////////////////////////////

func buildProductListQuery(queryFilter *QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(productTableHeaders).
		From("products").
		Where(squirrel.Eq{"archived_on": nil}).
		Limit(uint64(queryFilter.Limit))

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductUpdateQuery(p *Product) (string, []interface{}) {
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
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductCreationQuery(p *Product) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("products").
		Columns(
			"name",
			"subtitle",
			"description",
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
		).
		Values(
			p.Name,
			p.Subtitle.String,
			p.Description,
			p.SKU,
			p.UPC.String,
			p.Manufacturer.String,
			p.Brand.String,
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
		).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Options                    //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionListQuery(productID uint64, queryFilter *QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(fmt.Sprintf("count(id) over (), %s", productOptionsHeaders)).
		From("product_options").
		Where(squirrel.Eq{"product_id": productID}).
		Where(squirrel.Eq{"archived_on": nil})
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionUpdateQuery(a *ProductOption) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	productOptionUpdateSetMap := map[string]interface{}{
		"name":       a.Name,
		"updated_on": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_options").
		SetMap(productOptionUpdateSetMap).
		Where(squirrel.Eq{"id": a.ID}).
		Suffix(`RETURNING *`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionCreationQuery(a *ProductOption, productID uint64) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_options").
		Columns("name", "product_id").
		Values(a.Name, productID).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//               Product Option Values                //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionValueUpdateQuery(v *ProductOptionValue) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	productOptionUpdateSetMap := map[string]interface{}{
		"value":      v.Value,
		"updated_on": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_option_values").
		SetMap(productOptionUpdateSetMap).
		Where(squirrel.Eq{"id": v.ID}).
		Suffix(`RETURNING *`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionValueCreationQuery(v *ProductOptionValue) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_option_values").
		Columns("product_option_id", "value").
		Values(v.ProductOptionID, v.Value).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                     Discounts                      //
//                                                    //
////////////////////////////////////////////////////////

func buildDiscountListQuery(queryFilter *QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(discountDBColumns).
		From("discounts").
		Where(squirrel.Or{squirrel.Eq{"expires_on": nil}, squirrel.Gt{"expires_on": "NOW()"}}).
		Where(squirrel.Eq{"archived_on": nil})

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter, true)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildDiscountCreationQuery(d *Discount) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("discounts").
		Columns("name", "type", "amount", "starts_on", "expires_on", "requires_code", "code", "limited_use", "number_of_uses", "login_required").
		Values(d.Name, d.Type, d.Amount, d.StartsOn, d.ExpiresOn, d.RequiresCode, d.Code, d.LimitedUse, d.NumberOfUses, d.LoginRequired).
		Suffix(fmt.Sprintf("RETURNING%s", discountDBColumns))
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildDiscountUpdateQuery(d *Discount) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	updateSetMap := map[string]interface{}{
		"name":           d.Name,
		"type":           d.Type,
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
		Suffix(fmt.Sprintf("RETURNING%s", discountDBColumns))
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                       Users                        //
//                                                    //
////////////////////////////////////////////////////////

func buildUserCreationQuery(u *User) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("users").
		Columns(
			"first_name",
			"last_name",
			"email",
			"password",
			"salt",
			"is_admin",
		).
		Values(
			u.FirstName,
			u.LastName,
			u.Email,
			u.Password,
			u.Salt,
			u.IsAdmin,
		).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}
