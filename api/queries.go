package api

import "github.com/Masterminds/squirrel"

func applyQueryFilterToQueryBuilder(queryBuilder squirrel.SelectBuilder, queryFilter *QueryFilter) squirrel.SelectBuilder {
	if queryFilter.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(queryFilter.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(25)
	}

	if queryFilter.Page > 1 {
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
		Select("count(p.id) over (), *").
		From("products p").
		Join("product_progenitors g ON p.product_progenitor_id = g.id").
		Where(squirrel.Eq{"p.archived_on": nil}).
		Limit(uint64(queryFilter.Limit))

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)

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
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)
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
		Select("count(id) over (), *").
		From("discounts").
		Where(squirrel.Or{squirrel.Eq{"expires_on": nil}, squirrel.Gt{"expires_on": "NOW()"}}).
		Where(squirrel.Eq{"archived_on": nil}).
		Limit(uint64(queryFilter.Limit))

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}
