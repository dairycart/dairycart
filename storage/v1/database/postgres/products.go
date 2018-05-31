package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"

	"github.com/Masterminds/squirrel"
)

const productQueryBySKU = `
    SELECT
        id,
        product_root_id,
        primary_image_id,
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
    FROM
        products
    WHERE
        archived_on is null
    AND
        sku = $1
`

func (pg *postgres) GetProductBySKU(db database.Querier, sku string) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productQueryBySKU, sku).Scan(&p.ID, &p.ProductRootID, &p.PrimaryImageID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

const productWithSKUExistenceQuery = `SELECT EXISTS(SELECT id FROM products WHERE sku = $1 and archived_on IS NULL);`

func (pg *postgres) ProductWithSKUExists(db database.Querier, sku string) (bool, error) {
	var exists string

	err := db.QueryRow(productWithSKUExistenceQuery, sku).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productQueryByProductRootID = `
    SELECT
        id,
        product_root_id,
        primary_image_id,
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
    FROM
        products
    WHERE
        product_root_id = $1
`

func (pg *postgres) GetProductsByProductRootID(db database.Querier, productRootID uint64) ([]models.Product, error) {
	var list []models.Product

	rows, err := db.Query(productQueryByProductRootID, productRootID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.ProductRootID,
			&p.PrimaryImageID,
			&p.Name,
			&p.Subtitle,
			&p.Description,
			&p.OptionSummary,
			&p.SKU,
			&p.UPC,
			&p.Manufacturer,
			&p.Brand,
			&p.Quantity,
			&p.Taxable,
			&p.Price,
			&p.OnSale,
			&p.SalePrice,
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

const productExistenceQuery = `SELECT EXISTS(SELECT id FROM products WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductExists(db database.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productSelectionQuery = `
    SELECT
        id,
        product_root_id,
        primary_image_id,
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
    FROM
        products
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProduct(db database.Querier, id uint64) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productSelectionQuery, id).Scan(&p.ID, &p.ProductRootID, &p.PrimaryImageID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

func buildProductListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"product_root_id",
			"primary_image_id",
			"name",
			"subtitle",
			"description",
			"option_summary",
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
			"created_on",
			"updated_on",
			"archived_on",
		).
		From("products")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductList(db database.Querier, qf *models.QueryFilter) ([]models.Product, error) {
	var list []models.Product
	query, args := buildProductListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.ProductRootID,
			&p.PrimaryImageID,
			&p.Name,
			&p.Subtitle,
			&p.Description,
			&p.OptionSummary,
			&p.SKU,
			&p.UPC,
			&p.Manufacturer,
			&p.Brand,
			&p.Quantity,
			&p.Taxable,
			&p.Price,
			&p.OnSale,
			&p.SalePrice,
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

func buildProductCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("products")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetProductCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productCreationQuery = `
    INSERT INTO products
        (
            product_root_id, primary_image_id, name, subtitle, description, option_summary, sku, upc, manufacturer, brand, quantity, taxable, price, on_sale, sale_price, cost, product_weight, product_height, product_width, product_length, package_weight, package_height, package_width, package_length, quantity_per_package, available_on
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
        )
    RETURNING
        id, created_on, available_on;
`

func (pg *postgres) CreateProduct(db database.Querier, nu *models.Product) (createdID uint64, createdOn time.Time, availableOn time.Time, err error) {
	err = db.QueryRow(productCreationQuery, &nu.ProductRootID, &nu.PrimaryImageID, &nu.Name, &nu.Subtitle, &nu.Description, &nu.OptionSummary, &nu.SKU, &nu.UPC, &nu.Manufacturer, &nu.Brand, &nu.Quantity, &nu.Taxable, &nu.Price, &nu.OnSale, &nu.SalePrice, &nu.Cost, &nu.ProductWeight, &nu.ProductHeight, &nu.ProductWidth, &nu.ProductLength, &nu.PackageWeight, &nu.PackageHeight, &nu.PackageWidth, &nu.PackageLength, &nu.QuantityPerPackage, &nu.AvailableOn).Scan(&createdID, &createdOn, &availableOn)
	return createdID, createdOn, availableOn, err
}

const productUpdateQuery = `
    UPDATE products
    SET
        product_root_id = $1,
        primary_image_id = $2,
        name = $3,
        subtitle = $4,
        description = $5,
        option_summary = $6,
        sku = $7,
        upc = $8,
        manufacturer = $9,
        brand = $10,
        quantity = $11,
        taxable = $12,
        price = $13,
        on_sale = $14,
        sale_price = $15,
        cost = $16,
        product_weight = $17,
        product_height = $18,
        product_width = $19,
        product_length = $20,
        package_weight = $21,
        package_height = $22,
        package_width = $23,
        package_length = $24,
        quantity_per_package = $25,
        available_on = $26,
        updated_on = NOW()
    WHERE id = $27
    RETURNING updated_on;
`

func (pg *postgres) UpdateProduct(db database.Querier, updated *models.Product) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productUpdateQuery, &updated.ProductRootID, &updated.PrimaryImageID, &updated.Name, &updated.Subtitle, &updated.Description, &updated.OptionSummary, &updated.SKU, &updated.UPC, &updated.Manufacturer, &updated.Brand, &updated.Quantity, &updated.Taxable, &updated.Price, &updated.OnSale, &updated.SalePrice, &updated.Cost, &updated.ProductWeight, &updated.ProductHeight, &updated.ProductWidth, &updated.ProductLength, &updated.PackageWeight, &updated.PackageHeight, &updated.PackageWidth, &updated.PackageLength, &updated.QuantityPerPackage, &updated.AvailableOn, &updated.ID).Scan(&t)
	return t, err
}

const productDeletionQuery = `
    UPDATE products
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProduct(db database.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productDeletionQuery, id).Scan(&t)
	return t, err
}

const productWithProductRootIDDeletionQuery = `
    UPDATE products
    SET archived_on = NOW()
    WHERE product_root_id = $1
    RETURNING archived_on
`

func (pg *postgres) ArchiveProductsWithProductRootID(db database.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productWithProductRootIDDeletionQuery, id).Scan(&t)
	return t, err
}
