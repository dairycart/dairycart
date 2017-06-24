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
		queryBuilder = queryBuilder.Where(squirrel.Gt{"p.created_on": queryFilter.CreatedAfter})
	}

	if !queryFilter.CreatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"p.created_on": queryFilter.CreatedBefore})
	}

	if !queryFilter.UpdatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"p.updated_on": queryFilter.UpdatedAfter})
	}

	if !queryFilter.UpdatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"p.updated_on": queryFilter.UpdatedBefore})
	}
	return queryBuilder
}

////////////////////////////////////////////////////////
//                                                    //
//                Product Progenitors                 //
//                                                    //
////////////////////////////////////////////////////////

func buildProgenitorCreationQuery(g *ProductProgenitor) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_progenitors").
		Columns(
			"name",
			"description",
			"taxable",
			"price",
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
			g.Name,
			g.Description,
			g.Taxable,
			g.Price,
			g.Cost,
			g.ProductWeight,
			g.ProductHeight,
			g.ProductWidth,
			g.ProductLength,
			g.PackageWeight,
			g.PackageHeight,
			g.PackageWidth,
			g.PackageLength,
		).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                     Products                       //
//                                                    //
////////////////////////////////////////////////////////

func buildProductListQuery(queryFilter *QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(`p.id as product_id,
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
				g.*`).
		From("products p").
		Join("product_progenitors g ON p.product_progenitor_id = g.id").
		Where(squirrel.Eq{"p.archived_on": nil}).
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
		Columns("product_progenitor_id", "sku", "name", "upc", "quantity", "price", "cost").
		Values(p.ProductProgenitorID, p.SKU, p.Name, p.UPC, p.Quantity, p.Price, p.Cost).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Options                    //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionListQuery(progenitorID string, queryFilter *QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select("count(id) over (), *").
		From("product_options").
		Where(squirrel.Eq{"product_progenitor_id": progenitorID}).
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

func buildProductOptionCreationQuery(a *ProductOption) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Insert("product_options").
		Columns("name", "product_progenitor_id").
		Values(a.Name, a.ProductProgenitorID).
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
		Suffix(fmt.Sprintf("RETURNING %s", discountDBColumns))
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
		Suffix(fmt.Sprintf("RETURNING %s", discountDBColumns))
	query, args, _ := queryBuilder.ToSql()
	return query, args
}
