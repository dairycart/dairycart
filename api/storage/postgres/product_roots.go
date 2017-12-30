package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const productRootWithSKUPrefixExistenceQuery = `SELECT EXISTS(SELECT id FROM product_roots WHERE sku_prefix = $1 and archived_on IS NULL);`

func (pg *postgres) ProductRootWithSKUPrefixExists(db storage.Querier, skuPrefix string) (bool, error) {
	var exists string

	err := db.QueryRow(productRootWithSKUPrefixExistenceQuery, skuPrefix).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productRootExistenceQuery = `SELECT EXISTS(SELECT id FROM product_roots WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductRootExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productRootExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productRootSelectionQuery = `
    SELECT
        id,
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
    FROM
        product_roots
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProductRoot(db storage.Querier, id uint64) (*models.ProductRoot, error) {
	p := &models.ProductRoot{}

	err := db.QueryRow(productRootSelectionQuery, id).Scan(&p.ID, &p.Name, &p.Subtitle, &p.Description, &p.SKUPrefix, &p.Manufacturer, &p.Brand, &p.Taxable, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

func buildProductRootListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"name",
			"subtitle",
			"description",
			"sku_prefix",
			"manufacturer",
			"brand",
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
			"quantity_per_package",
			"available_on",
			"created_on",
			"updated_on",
			"archived_on",
		).
		From("product_roots")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductRootList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductRoot, error) {
	var list []models.ProductRoot
	query, args := buildProductRootListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductRoot
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Subtitle,
			&p.Description,
			&p.SKUPrefix,
			&p.Manufacturer,
			&p.Brand,
			&p.Taxable,
			&p.Cost,
			&p.ProductWeight,
			&p.ProductHeight,
			&p.ProductWidth,
			&p.ProductLength,
			&p.PackageWeight,
			&p.PackageHeight,
			&p.PackageWidth,
			&p.PackageLength,
			&p.QuantityPerPackage,
			&p.AvailableOn,
			&p.CreatedOn,
			&p.UpdatedOn,
			&p.ArchivedOn,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildProductRootCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("product_roots")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetProductRootCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductRootCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productRootCreationQuery = `
    INSERT INTO product_roots
        (
            name, subtitle, description, sku_prefix, manufacturer, brand, taxable, cost, product_weight, product_height, product_width, product_length, package_weight, package_height, package_width, package_length, quantity_per_package, available_on
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateProductRoot(db storage.Querier, nu *models.ProductRoot) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(productRootCreationQuery, &nu.Name, &nu.Subtitle, &nu.Description, &nu.SKUPrefix, &nu.Manufacturer, &nu.Brand, &nu.Taxable, &nu.Cost, &nu.ProductWeight, &nu.ProductHeight, &nu.ProductWidth, &nu.ProductLength, &nu.PackageWeight, &nu.PackageHeight, &nu.PackageWidth, &nu.PackageLength, &nu.QuantityPerPackage, &nu.AvailableOn).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const productRootUpdateQuery = `
    UPDATE product_roots
    SET
        name = $1,
        subtitle = $2,
        description = $3,
        sku_prefix = $4,
        manufacturer = $5,
        brand = $6,
        taxable = $7,
        cost = $8,
        product_weight = $9,
        product_height = $10,
        product_width = $11,
        product_length = $12,
        package_weight = $13,
        package_height = $14,
        package_width = $15,
        package_length = $16,
        quantity_per_package = $17,
        available_on = $18,
        updated_on = NOW()
    WHERE id = $19
    RETURNING updated_on;
`

func (pg *postgres) UpdateProductRoot(db storage.Querier, updated *models.ProductRoot) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productRootUpdateQuery, &updated.Name, &updated.Subtitle, &updated.Description, &updated.SKUPrefix, &updated.Manufacturer, &updated.Brand, &updated.Taxable, &updated.Cost, &updated.ProductWeight, &updated.ProductHeight, &updated.ProductWidth, &updated.ProductLength, &updated.PackageWeight, &updated.PackageHeight, &updated.PackageWidth, &updated.PackageLength, &updated.QuantityPerPackage, &updated.AvailableOn, &updated.ID).Scan(&t)
	return t, err
}

const productRootDeletionQuery = `
    UPDATE product_roots
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProductRoot(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productRootDeletionQuery, id).Scan(&t)
	return t, err
}
