package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const productQueryBySKU = `
    SELECT
        product_height,
        package_width,
        on_sale,
        subtitle,
        product_length,
        taxable,
        cost,
        brand,
        product_root_id,
        product_weight,
        quantity,
        manufacturer,
        product_width,
        available_on,
        quantity_per_package,
        package_length,
        price,
        primary_image_id,
        option_summary,
        upc,
        name,
        package_height,
        sale_price,
        id,
        package_weight,
        updated_on,
        description,
        archived_on,
        created_on,
        sku
    FROM
        products
    WHERE
        archived_on is null
    AND
        sku = $1
`

func (pg *postgres) GetProductBySKU(db storage.Querier, sku string) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productQueryBySKU, sku).Scan(&p.ProductHeight, &p.PackageWidth, &p.OnSale, &p.Subtitle, &p.ProductLength, &p.Taxable, &p.Cost, &p.Brand, &p.ProductRootID, &p.ProductWeight, &p.Quantity, &p.Manufacturer, &p.ProductWidth, &p.AvailableOn, &p.QuantityPerPackage, &p.PackageLength, &p.Price, &p.PrimaryImageID, &p.OptionSummary, &p.UPC, &p.Name, &p.PackageHeight, &p.SalePrice, &p.ID, &p.PackageWeight, &p.UpdatedOn, &p.Description, &p.ArchivedOn, &p.CreatedOn, &p.SKU)

	return p, err
}

const productWithSKUExistenceQuery = `SELECT EXISTS(SELECT id FROM products WHERE sku = $1 and archived_on IS NULL);`

func (pg *postgres) ProductWithSKUExists(db storage.Querier, sku string) (bool, error) {
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
        product_height,
        package_width,
        on_sale,
        subtitle,
        product_length,
        taxable,
        cost,
        brand,
        product_root_id,
        product_weight,
        quantity,
        manufacturer,
        product_width,
        available_on,
        quantity_per_package,
        package_length,
        price,
        primary_image_id,
        option_summary,
        upc,
        name,
        package_height,
        sale_price,
        id,
        package_weight,
        updated_on,
        description,
        archived_on,
        created_on,
        sku
    FROM
        products
    WHERE
        product_root_id = $1
`

func (pg *postgres) GetProductsByProductRootID(db storage.Querier, productRootID uint64) ([]models.Product, error) {
	var list []models.Product

	rows, err := db.Query(productQueryByProductRootID, productRootID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ProductHeight,
			&p.PackageWidth,
			&p.OnSale,
			&p.Subtitle,
			&p.ProductLength,
			&p.Taxable,
			&p.Cost,
			&p.Brand,
			&p.ProductRootID,
			&p.ProductWeight,
			&p.Quantity,
			&p.Manufacturer,
			&p.ProductWidth,
			&p.AvailableOn,
			&p.QuantityPerPackage,
			&p.PackageLength,
			&p.Price,
			&p.PrimaryImageID,
			&p.OptionSummary,
			&p.UPC,
			&p.Name,
			&p.PackageHeight,
			&p.SalePrice,
			&p.ID,
			&p.PackageWeight,
			&p.UpdatedOn,
			&p.Description,
			&p.ArchivedOn,
			&p.CreatedOn,
			&p.SKU,
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

func (pg *postgres) ProductExists(db storage.Querier, id uint64) (bool, error) {
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
        product_height,
        package_width,
        on_sale,
        subtitle,
        product_length,
        taxable,
        cost,
        brand,
        product_root_id,
        product_weight,
        quantity,
        manufacturer,
        product_width,
        available_on,
        quantity_per_package,
        package_length,
        price,
        primary_image_id,
        option_summary,
        upc,
        name,
        package_height,
        sale_price,
        id,
        package_weight,
        updated_on,
        description,
        archived_on,
        created_on,
        sku
    FROM
        products
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProduct(db storage.Querier, id uint64) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productSelectionQuery, id).Scan(&p.ProductHeight, &p.PackageWidth, &p.OnSale, &p.Subtitle, &p.ProductLength, &p.Taxable, &p.Cost, &p.Brand, &p.ProductRootID, &p.ProductWeight, &p.Quantity, &p.Manufacturer, &p.ProductWidth, &p.AvailableOn, &p.QuantityPerPackage, &p.PackageLength, &p.Price, &p.PrimaryImageID, &p.OptionSummary, &p.UPC, &p.Name, &p.PackageHeight, &p.SalePrice, &p.ID, &p.PackageWeight, &p.UpdatedOn, &p.Description, &p.ArchivedOn, &p.CreatedOn, &p.SKU)

	return p, err
}

func buildProductListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"product_height",
			"package_width",
			"on_sale",
			"subtitle",
			"product_length",
			"taxable",
			"cost",
			"brand",
			"product_root_id",
			"product_weight",
			"quantity",
			"manufacturer",
			"product_width",
			"available_on",
			"quantity_per_package",
			"package_length",
			"price",
			"primary_image_id",
			"option_summary",
			"upc",
			"name",
			"package_height",
			"sale_price",
			"id",
			"package_weight",
			"updated_on",
			"description",
			"archived_on",
			"created_on",
			"sku",
		).
		From("products")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductList(db storage.Querier, qf *models.QueryFilter) ([]models.Product, error) {
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
			&p.ProductHeight,
			&p.PackageWidth,
			&p.OnSale,
			&p.Subtitle,
			&p.ProductLength,
			&p.Taxable,
			&p.Cost,
			&p.Brand,
			&p.ProductRootID,
			&p.ProductWeight,
			&p.Quantity,
			&p.Manufacturer,
			&p.ProductWidth,
			&p.AvailableOn,
			&p.QuantityPerPackage,
			&p.PackageLength,
			&p.Price,
			&p.PrimaryImageID,
			&p.OptionSummary,
			&p.UPC,
			&p.Name,
			&p.PackageHeight,
			&p.SalePrice,
			&p.ID,
			&p.PackageWeight,
			&p.UpdatedOn,
			&p.Description,
			&p.ArchivedOn,
			&p.CreatedOn,
			&p.SKU,
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

func (pg *postgres) GetProductCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productCreationQuery = `
    INSERT INTO products
        (
            product_height, package_width, on_sale, subtitle, product_length, taxable, cost, brand, product_root_id, product_weight, quantity, manufacturer, product_width, available_on, quantity_per_package, package_length, price, primary_image_id, option_summary, upc, name, package_height, sale_price, package_weight, description, sku
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
        )
    RETURNING
        id, created_on, available_on;
`

func (pg *postgres) CreateProduct(db storage.Querier, nu *models.Product) (createdID uint64, createdOn time.Time, availableOn time.Time, err error) {
	err = db.QueryRow(productCreationQuery, &nu.ProductHeight, &nu.PackageWidth, &nu.OnSale, &nu.Subtitle, &nu.ProductLength, &nu.Taxable, &nu.Cost, &nu.Brand, &nu.ProductRootID, &nu.ProductWeight, &nu.Quantity, &nu.Manufacturer, &nu.ProductWidth, &nu.AvailableOn, &nu.QuantityPerPackage, &nu.PackageLength, &nu.Price, &nu.PrimaryImageID, &nu.OptionSummary, &nu.UPC, &nu.Name, &nu.PackageHeight, &nu.SalePrice, &nu.PackageWeight, &nu.Description, &nu.SKU).Scan(&createdID, &createdOn, &availableOn)
	return createdID, createdOn, availableOn, err
}

const productUpdateQuery = `
    UPDATE products
    SET
        product_height = $1,
        package_width = $2,
        on_sale = $3,
        subtitle = $4,
        product_length = $5,
        taxable = $6,
        cost = $7,
        brand = $8,
        product_root_id = $9,
        product_weight = $10,
        quantity = $11,
        manufacturer = $12,
        product_width = $13,
        available_on = $14,
        quantity_per_package = $15,
        package_length = $16,
        price = $17,
        primary_image_id = $18,
        option_summary = $19,
        upc = $20,
        name = $21,
        package_height = $22,
        sale_price = $23,
        package_weight = $24,
        description = $25,
        sku = $26,
        updated_on = NOW()
    WHERE id = $27
    RETURNING updated_on;
`

func (pg *postgres) UpdateProduct(db storage.Querier, updated *models.Product) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productUpdateQuery, &updated.ProductHeight, &updated.PackageWidth, &updated.OnSale, &updated.Subtitle, &updated.ProductLength, &updated.Taxable, &updated.Cost, &updated.Brand, &updated.ProductRootID, &updated.ProductWeight, &updated.Quantity, &updated.Manufacturer, &updated.ProductWidth, &updated.AvailableOn, &updated.QuantityPerPackage, &updated.PackageLength, &updated.Price, &updated.PrimaryImageID, &updated.OptionSummary, &updated.UPC, &updated.Name, &updated.PackageHeight, &updated.SalePrice, &updated.PackageWeight, &updated.Description, &updated.SKU, &updated.ID).Scan(&t)
	return t, err
}

const productDeletionQuery = `
    UPDATE products
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProduct(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productDeletionQuery, id).Scan(&t)
	return t, err
}

const productWithProductRootIDDeletionQuery = `
    UPDATE products
    SET archived_on = NOW()
    WHERE product_root_id = $1
    RETURNING archived_on
`

func (pg *postgres) ArchiveProductsWithProductRootID(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productWithProductRootIDDeletionQuery, id).Scan(&t)
	return t, err
}
