package api

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

var sqlBuilder squirrel.StatementBuilderType

func init() {
	sqlBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

////////////////////////////////////////////////////////
//                                                    //
//               General Query Builders               //
//                                                    //
////////////////////////////////////////////////////////

func buildRowExistenceQuery(table string, idColumn string, id interface{}) string {
	subqueryBuilder := sqlBuilder.Select("1").From(table).Where(squirrel.Eq{idColumn: id}).Where(squirrel.Eq{"archived_at": nil})
	subquery, _, _ := subqueryBuilder.ToSql()

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS(%s)", subquery))
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildRowRetrievalQuery(table string, idColumn string, id interface{}) string {
	queryBuilder := sqlBuilder.Select("*").From(table).Where(squirrel.Eq{idColumn: id}).Where(squirrel.Eq{"archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildRowDeletionQuery(table string, idColumn string, id interface{}) string {
	queryBuilder := sqlBuilder.
		Update(table).
		Set("archived_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{idColumn: id}).
		Where(squirrel.Eq{"archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

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
		queryBuilder = queryBuilder.Where(squirrel.Gt{"p.created_at": queryFilter.CreatedAfter})
	}

	if !queryFilter.CreatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"p.created_at": queryFilter.CreatedBefore})
	}

	if !queryFilter.UpdatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"p.updated_at": queryFilter.UpdatedAfter})
	}

	if !queryFilter.UpdatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"p.updated_at": queryFilter.UpdatedBefore})
	}
	return queryBuilder
}

////////////////////////////////////////////////////////
//                                                    //
//                Product Progenitors                 //
//                                                    //
////////////////////////////////////////////////////////

func buildProgenitorRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_progenitors", "id", id)
}

func buildProgenitorExistenceQuery(id string) string {
	return buildRowExistenceQuery("product_progenitors", "id", id)
}

func buildProgenitorCreationQuery(g *ProductProgenitor) (string, []interface{}) {
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

func buildProductExistenceQuery(sku string) string {
	return buildRowExistenceQuery("products", "sku", sku)
}

func buildProductRetrievalQuery(sku string) string {
	return buildRowRetrievalQuery("products", "sku", sku)
}

func buildProductDeletionQuery(sku string) string {
	return buildRowDeletionQuery("products", "sku", sku)
}

func buildProductListQuery(queryFilter *QueryFilter) (string, []interface{}) {
	queryBuilder := sqlBuilder.
		Select("count(p.id) over (), *").
		From("products p").
		Join("product_progenitors g ON p.product_progenitor_id = g.id").
		Where(squirrel.Eq{"p.archived_at": nil}).
		Limit(uint64(queryFilter.Limit))

	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)

	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildCompleteProductRetrievalQuery(sku string) string {
	queryBuilder := sqlBuilder.
		Select("*").
		From("products p").
		Join("product_progenitors g ON p.product_progenitor_id = g.id").
		Where(squirrel.Eq{"p.sku": sku}).
		Where(squirrel.Eq{"p.archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductUpdateQuery(p *Product) (string, []interface{}) {
	productUpdateSetMap := map[string]interface{}{
		"sku":        p.SKU,
		"name":       p.Name,
		"upc":        p.UPC,
		"quantity":   p.Quantity,
		"price":      p.Price,
		"cost":       p.Cost,
		"updated_at": squirrel.Expr("NOW()"),
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
//                Product Options                  //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionExistenceQuery(id int64) string {
	return buildRowExistenceQuery("product_options", "id", id)
}

func buildProductOptionRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_options", "id", id)
}

func buildProductOptionDeletionQuery(id int64) string {
	return buildRowDeletionQuery("product_options", "id", id)
}

func buildProductOptionExistenceQueryForProductByName(name, progenitorID string) string {
	subqueryBuilder := sqlBuilder.Select("1").
		From("product_options").
		Where(squirrel.Eq{"name": name}).
		Where(squirrel.Eq{"product_progenitor_id": progenitorID}).
		Where(squirrel.Eq{"archived_at": nil})
	subquery, _, _ := subqueryBuilder.ToSql()

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS(%s)", subquery))
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductOptionListQuery(progenitorID string, queryFilter *QueryFilter) string {
	queryBuilder := sqlBuilder.
		Select("count(id) over (), *").
		From("product_options").
		Where(squirrel.Eq{"product_progenitor_id": progenitorID}).
		Where(squirrel.Eq{"archived_at": nil})
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductOptionUpdateQuery(a *ProductOption) (string, []interface{}) {
	productOptionUpdateSetMap := map[string]interface{}{
		"name":       a.Name,
		"updated_at": squirrel.Expr("NOW()"),
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
//             Product Option Values               //
//                                                    //
////////////////////////////////////////////////////////

func buildProductOptionValueExistenceQuery(id int64) string {
	return buildRowExistenceQuery("product_option_values", "id", id)
}

func buildProductOptionValueRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_option_values", "id", id)
}

func buildProductOptionValueDeletionQuery(id int64) string {
	return buildRowDeletionQuery("product_option_values", "id", id)
}

func buildProductOptionValueRetrievalForOptionIDQuery(optionID int64) string {
	queryBuilder := sqlBuilder.Select("*").
		From("product_option_values").
		Where(squirrel.Eq{"product_option_id": optionID}).
		Where(squirrel.Eq{"archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductOptionValueExistenceForOptionIDQuery(optionID int64, value string) (string, []interface{}) {
	subqueryBuilder := sqlBuilder.Select("1").
		From("product_option_values").
		Where(squirrel.Eq{"product_option_id": optionID}).
		Where(squirrel.Eq{"value": value}).
		Where(squirrel.Eq{"archived_at": nil})
	subquery, args, _ := subqueryBuilder.ToSql()

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS(%s)", subquery))
	query, _, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductOptionValueUpdateQuery(v *ProductOptionValue) (string, []interface{}) {
	productOptionUpdateSetMap := map[string]interface{}{
		"value":      v.Value,
		"updated_at": squirrel.Expr("NOW()"),
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
//                    Discounts                       //
//                                                    //
////////////////////////////////////////////////////////

func buildDiscountExistenceQuery(id string) string {
	return buildRowExistenceQuery("discounts", "id", id)
}

func buildDiscountRetrievalQuery(id string) string {
	return buildRowRetrievalQuery("discounts", "id", id)
}

func buildDiscountDeletionQuery(id string) string {
	return buildRowDeletionQuery("discounts", "id", id)
}
