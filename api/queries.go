package api

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

/////////////////////////////////////////////////////////////////////////////
//                        ,---.           ,---.                            //
//                       / /"`.\.--"""--./,'"\ \                           //
//                       \ \    _       _    / /                           //
//                        `./  / __   __ \  \,'                            //
//                         /    /_O)_(_O\    \                             //
//                         |  .-'  ___  `-.  |                             //
//                      .--|       \_/       |--.                          //
//                    ,'    \   \   |   /   /    `.                        //
//                   /       `.  `--^--'  ,'       \                       //
//                .-"""""-.    `--.___.--'     .-"""""-.                   //
//   .-----------/         \------------------/         \--------------.   //
//   | .---------\         /----------------- \         /------------. |   //
//   | |          `-`--`--'                    `--'--'-'             | |   //
//   | |                                                             | |   //
//   | |                Generalized Query Builders                   | |   //
//   | |                                                             | |   //
//   | |_____________________________________________________________| |   //
//   |_________________________________________________________________|   //
//                      )__________|__|__________(                         //
//                     |            ||            |                        //
//                     |____________||____________|                        //
//                       ),-----.(      ),-----.(                          //
//                     ,'   ==.   \    /  .==    `.                        //
//                    /            )  (            \                       //
//                    `==========='    `==========='                       //
/////////////////////////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////////////////////////////
//                                                                              //
//   ........................................................................   //
//   :      ,~~.          ,~~.          ,~~.          ,~~.          ,~~.    :   //
//   :     (  6 )-_,     (  6 )-_,     (  6 )-_,     (  6 )-_,     (  6 )-_,:   //
//   :(\___ )=='-'  (\___ )=='-'  (\___ )=='-'  (\___ )=='-'  (\___ )=='-'  :   //
//   : \ .   ) )     \ .   ) )     \ .   ) )     \ .   ) )     \ .   ) )    :   //
//   :  \ `-' /       \ `-' /       \ `-' /       \ `-' /       \ `-' /     :   //
//   : ~'`~'`~'`~`~'`~'`~'`~'`~`~'`~'`~'`~'`~`~'`~'`~'`~'`~'`~`~'`~'`~'`~'` :   //
//   :      ,~~.    ..........................................      ,~~.    :   //
//   :     (  9 )-_,:                                        :     (  9 )-_,:   //
//   :(\___ )=='-'  :                                        :(\___ )=='-'  :   //
//   : \ .   ) )    :       Product Progenitor Queries       : \ .   ) )    :   //
//   :  \ `-' /     :                                        :  \ `-' /     :   //
//   :   `~j-'      :                                        :   `~j-'      :   //
//   :     '=:      :........................................:     '=:      :   //
//   :      ,~~.          ,~~.          ,~~.          ,~~.          ,~~.    :   //
//   :     (  6 )-_,     (  6 )-_,     (  6 )-_,     (  6 )-_,     (  6 )-_,:   //
//   :(\___ )=='-'  (\___ )=='-'  (\___ )=='-'  (\___ )=='-'  (\___ )=='-'  :   //
//   : \ .   ) )     \ .   ) )     \ .   ) )     \ .   ) )     \ .   ) )    :   //
//   :  \ `-' /       \ `-' /       \ `-' /       \ `-' /       \ `-' /     :   //
//   : ~'`~'`~'`~`~'`~'`~'`~'`~`~'`~'`~'`~'`~`~'`~'`~'`~'`~`~'`~'`~'`~'`~'` :   //
//   :......................................................................:   //
//                                                                              //
//////////////////////////////////////////////////////////////////////////////////

func buildProgenitorRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_progenitors", "id", id)
}

func buildProgenitorExistenceQuery(id int64) string {
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

/////////////////////////////////////////////////////
//                                ,--._,--.        //
//                              ,'  ,'   ,-`.      //
//                   (`-.__    /  ,'   /           //
//                    `.   `--'        \__,--'-.   //
//    ______________    `--/       ,-.  ______/    //
//   |              |     (o-.     ,o- /           //
//   |   Product    |     `. ;         \           //
//   |   Queries    |      |:           \          //
//   |______________|     ,'`       ,    \         //
//                   \    (o o ,  --'     :        //
//                    \-   \--','.        ;        //
//                          `;;  :       /         //
//                           ;'  ;  ,' ,'          //
//                           ,','  :  '            //
//                           \ \   :               //
//                            `                    //
/////////////////////////////////////////////////////

func buildProductExistenceQuery(sku string) string {
	return buildRowExistenceQuery("products", "sku", sku)
}

func buildProductRetrievalQuery(sku string) string {
	return buildRowRetrievalQuery("products", "sku", sku)
}

func buildProductDeletionQuery(sku string) string {
	return buildRowDeletionQuery("products", "sku", sku)
}

func buildAllProductsRetrievalQuery() string {
	queryBuilder := sqlBuilder.
		Select("*").
		From("products p").
		Join("product_progenitors g ON p.product_progenitor_id = g.id").
		Where(squirrel.Eq{"p.archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
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
		"on_sale":    p.OnSale,
		"price":      p.Price,
		"sale_price": p.SalePrice,
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
		Columns("product_progenitor_id", "sku", "name", "upc", "quantity", "on_sale", "price", "sale_price").
		Values(p.ProductProgenitorID, p.SKU, p.Name, p.UPC, p.Quantity, p.OnSale, p.Price, p.SalePrice).
		Suffix(`RETURNING *`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}
