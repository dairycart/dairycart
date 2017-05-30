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
//   | |               Generalized Query Builders                    | |   //
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
//                   \_   (o o ,  --'     :        //
//                     `   \--','.        ;        //
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

func buildAllProductsRetrievalQuery(queryFilter *QueryFilter) (string, []interface{}) {
	queryBuilder := sqlBuilder.
		Select("*").
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

/////////////////////////////////////////////////////////////
//                                                         //
//          _..--""-.                  .-""--.._           //
//      _.-'         \ __...----...__ /         '-._       //
//    .'      .:::...,'              ',...:::.      '.     //
//   (     .'``'''::;                  ;::'''``'.     )    //
//    \             '-)              (-'             /     //
//     \             /                \             /      //
//      \          .'.-.            .-.'.          /       //
//       \         | \0|            |0/ |         /        //
//       |          \  |   .-==-.   |  /          |        //
//        \          `/`;          ;`\`          /         //
//         '.._      (_ |  .-==-.  | _)      _..'          //
//             `"`"-`/ `/'        '\` \`-"`"`              //
//                  / /`;   .==.   ;`\ \                   //
//            .---./_/   \  .==.  /   \ \                  //
//           / '.    `-.__)       |    `"                  //
//          | =(`-.        '==.   ;                        //
//           \  '. `-.           /         Product         //
//            \_:_)   `"--.....-'         Attributes       //
//                                                         //
/////////////////////////////////////////////////////////////

func buildProductAttributeExistenceQuery(id int64) string {
	return buildRowExistenceQuery("product_attributes", "id", id)
}

func buildProductAttributeRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_attributes", "id", id)
}

func buildProductAttributeDeletionQuery(id int64) string {
	return buildRowDeletionQuery("product_attributes", "id", id)
}

func buildProductAttributeExistenceQueryForProductByName(name, progenitorID string) string {
	subqueryBuilder := sqlBuilder.Select("1").
		From("product_attributes").
		Where(squirrel.Eq{"name": name}).
		Where(squirrel.Eq{"product_progenitor_id": progenitorID}).
		Where(squirrel.Eq{"archived_at": nil})
	subquery, _, _ := subqueryBuilder.ToSql()

	queryBuilder := sqlBuilder.Select(fmt.Sprintf("EXISTS(%s)", subquery))
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductAttributeListQuery(progenitorID string, queryFilter *QueryFilter) string {
	queryBuilder := sqlBuilder.
		Select("*").
		From("product_attributes").
		Where(squirrel.Eq{"product_progenitor_id": progenitorID}).
		Where(squirrel.Eq{"archived_at": nil})
	queryBuilder = applyQueryFilterToQueryBuilder(queryBuilder, queryFilter)
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductAttributeUpdateQuery(a *ProductAttribute) (string, []interface{}) {
	productAttributeUpdateSetMap := map[string]interface{}{
		"name":       a.Name,
		"updated_at": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_attributes").
		SetMap(productAttributeUpdateSetMap).
		Where(squirrel.Eq{"id": a.ID}).
		Suffix(`RETURNING *`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductAttributeCreationQuery(a *ProductAttribute) (string, []interface{}) {
	queryBuilder := sqlBuilder.
		Insert("product_attributes").
		Columns("name", "product_progenitor_id").
		Values(a.Name, a.ProductProgenitorID).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

/////////////////////////////////////////////////////
//                                    ___          //
//                                ,-""   `.        //
//      Product                 ,'  _   e )`-._    //
//         Attribute           /  ,' `-._<.===-'   //
//            Values          /  /                 //
//                           /  ;                  //
//               _.--.__    /   ;                  //
//  (`._    _.-""       "--'    |                  //
//  <_  `-""                     \                 //
//   <`-                          :                //
//    (__   <__.                  ;                //
//      `-.   '-.__.      _.'    /                 //
//         \      `-.__,-'    _,'                  //
//          `._    ,    /__,-'                     //
//             ""._\__,'< <____                    //
//                  | |  `----.`.                  //
//                  | |        \ `.                //
//                  ; |___      \-``               //
//                  \   --<                        //
//                   `.`.<                         //
//                     `-'                         //
//                                                 //
/////////////////////////////////////////////////////

func buildProductAttributeValueExistenceQuery(id int64) string {
	return buildRowExistenceQuery("product_attribute_values", "id", id)
}

func buildProductAttributeValueRetrievalQuery(id int64) string {
	return buildRowRetrievalQuery("product_attribute_values", "id", id)
}

func buildProductAttributeValueDeletionQuery(id int64) string {
	return buildRowDeletionQuery("product_attribute_values", "id", id)
}

func buildProductAttributeValueRetrievalForAttributeIDQuery(attributeID int64) string {
	queryBuilder := sqlBuilder.Select("*").
		From("product_attribute_values").
		Where(squirrel.Eq{"product_attribute_id": attributeID}).
		Where(squirrel.Eq{"archived_at": nil})
	query, _, _ := queryBuilder.ToSql()
	return query
}

func buildProductAttributeValueUpdateQuery(v *ProductAttributeValue) (string, []interface{}) {
	productAttributeUpdateSetMap := map[string]interface{}{
		"value":      v.Value,
		"updated_at": squirrel.Expr("NOW()"),
	}
	queryBuilder := sqlBuilder.
		Update("product_attribute_values").
		SetMap(productAttributeUpdateSetMap).
		Where(squirrel.Eq{"id": v.ID}).
		Suffix(`RETURNING *`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}

func buildProductAttributeValueCreationQuery(v *ProductAttributeValue) (string, []interface{}) {
	queryBuilder := sqlBuilder.
		Insert("product_attribute_values").
		Columns("product_attribute_id", "value").
		Values(v.ProductAttributeID, v.Value).
		Suffix(`RETURNING "id"`)
	query, args, _ := queryBuilder.ToSql()
	return query, args
}
