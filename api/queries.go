package api

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func buildRowExistenceQuery(table string, idColumn string, id interface{}) string {
	subqueryBuilder := sqlBuilder.Select("1").From(table).Where(squirrel.Eq{idColumn: id}).Where(squirrel.Eq{"archived_at": nil})
	subquery, _, _ := subqueryBuilder.ToSql()

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS(%s)", subquery))
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProgenitorRetrievalQuery(id int64) string {
	queryBuilder := sqlBuilder.Select("*").From("product_progenitors").Where(squirrel.Eq{"id": id}).Where(squirrel.Eq{"archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProgenitorExistenceQuery(id int64) string {
	return buildRowExistenceQuery("product_progenitors", "id", id)
}

func buildProgenitorCreationQuery(g *ProductProgenitor) string {
	queryBuilder := sqlBuilder.
		Insert("product_progenitors").
		Columns(
			"name",
			"description",
			"taxable",
			"price",
			"product_weight",
			"product_height",
			"product_width",
			"product_length",
			"package_weight",
			"package_height",
			"package_width",
			"package_length",
			"created_at",
		).
		Values(
			g.Name,
			g.Description,
			g.Taxable,
			g.Price,
			g.ProductWeight,
			g.ProductHeight,
			g.ProductWidth,
			g.ProductLength,
			g.PackageWeight,
			g.PackageHeight,
			g.PackageWidth,
			g.PackageLength,
			squirrel.Expr("NOW()"),
		).
		Suffix(`RETURNING "id"`)
	query, _, _ := queryBuilder.ToSql()
	return query
}
