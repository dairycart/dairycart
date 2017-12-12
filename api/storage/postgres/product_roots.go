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
        available_on,
        product_length,
        updated_on,
        sku_prefix,
        package_height,
        product_weight,
        product_width,
        quantity_per_package,
        name,
        product_height,
        package_length,
        created_on,
        cost,
        brand,
        subtitle,
        package_weight,
        archived_on,
        id,
        package_width,
        description,
        manufacturer,
        taxable
    FROM
        product_roots
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProductRoot(db storage.Querier, id uint64) (*models.ProductRoot, error) {
	p := &models.ProductRoot{}

	err := db.QueryRow(productRootSelectionQuery, id).Scan(&p.AvailableOn, &p.ProductLength, &p.UpdatedOn, &p.SKUPrefix, &p.PackageHeight, &p.ProductWeight, &p.ProductWidth, &p.QuantityPerPackage, &p.Name, &p.ProductHeight, &p.PackageLength, &p.CreatedOn, &p.Cost, &p.Brand, &p.Subtitle, &p.PackageWeight, &p.ArchivedOn, &p.ID, &p.PackageWidth, &p.Description, &p.Manufacturer, &p.Taxable)

	return p, err
}

func buildProductRootListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"available_on",
			"product_length",
			"updated_on",
			"sku_prefix",
			"package_height",
			"product_weight",
			"product_width",
			"quantity_per_package",
			"name",
			"product_height",
			"package_length",
			"created_on",
			"cost",
			"brand",
			"subtitle",
			"package_weight",
			"archived_on",
			"id",
			"package_width",
			"description",
			"manufacturer",
			"taxable",
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
			&p.AvailableOn,
			&p.ProductLength,
			&p.UpdatedOn,
			&p.SKUPrefix,
			&p.PackageHeight,
			&p.ProductWeight,
			&p.ProductWidth,
			&p.QuantityPerPackage,
			&p.Name,
			&p.ProductHeight,
			&p.PackageLength,
			&p.CreatedOn,
			&p.Cost,
			&p.Brand,
			&p.Subtitle,
			&p.PackageWeight,
			&p.ArchivedOn,
			&p.ID,
			&p.PackageWidth,
			&p.Description,
			&p.Manufacturer,
			&p.Taxable,
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
            available_on, product_length, sku_prefix, package_height, product_weight, product_width, quantity_per_package, name, product_height, package_length, cost, brand, subtitle, package_weight, package_width, description, manufacturer, taxable
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateProductRoot(db storage.Querier, nu *models.ProductRoot) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(productRootCreationQuery, &nu.AvailableOn, &nu.ProductLength, &nu.SKUPrefix, &nu.PackageHeight, &nu.ProductWeight, &nu.ProductWidth, &nu.QuantityPerPackage, &nu.Name, &nu.ProductHeight, &nu.PackageLength, &nu.Cost, &nu.Brand, &nu.Subtitle, &nu.PackageWeight, &nu.PackageWidth, &nu.Description, &nu.Manufacturer, &nu.Taxable).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const productRootUpdateQuery = `
    UPDATE product_roots
    SET
        available_on = $1, 
        product_length = $2, 
        updated_on = NOW()
        sku_prefix = $4, 
        package_height = $5, 
        product_weight = $6, 
        product_width = $7, 
        quantity_per_package = $8, 
        name = $9, 
        product_height = $10, 
        package_length = $11, 
        cost = $12, 
        brand = $13, 
        subtitle = $14, 
        package_weight = $15, 
        package_width = $16, 
        description = $17, 
        manufacturer = $18, 
        taxable = $19
    WHERE id = $19
    RETURNING updated_on;
`

func (pg *postgres) UpdateProductRoot(db storage.Querier, updated *models.ProductRoot) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productRootUpdateQuery, &updated.AvailableOn, &updated.ProductLength, &updated.SKUPrefix, &updated.PackageHeight, &updated.ProductWeight, &updated.ProductWidth, &updated.QuantityPerPackage, &updated.Name, &updated.ProductHeight, &updated.PackageLength, &updated.Cost, &updated.Brand, &updated.Subtitle, &updated.PackageWeight, &updated.PackageWidth, &updated.Description, &updated.Manufacturer, &updated.Taxable, &updated.ID).Scan(&t)
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
