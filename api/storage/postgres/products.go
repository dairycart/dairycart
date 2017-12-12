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
        product_width,
        package_length,
        sale_price,
        description,
        package_weight,
        price,
        product_weight,
        quantity,
        product_root_id,
        product_height,
        taxable,
        brand,
        product_length,
        created_on,
        available_on,
        quantity_per_package,
        on_sale,
        name,
        sku,
        manufacturer,
        subtitle,
        package_width,
        cost,
        id,
        package_height,
        archived_on,
        option_summary,
        updated_on,
        upc
    FROM
        products
    WHERE
        archived_on is null
    AND
        sku = $1
`

func (pg *postgres) GetProductBySKU(db storage.Querier, sku string) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productQueryBySKU, sku).Scan(&p.ProductWidth, &p.PackageLength, &p.SalePrice, &p.Description, &p.PackageWeight, &p.Price, &p.ProductWeight, &p.Quantity, &p.ProductRootID, &p.ProductHeight, &p.Taxable, &p.Brand, &p.ProductLength, &p.CreatedOn, &p.AvailableOn, &p.QuantityPerPackage, &p.OnSale, &p.Name, &p.SKU, &p.Manufacturer, &p.Subtitle, &p.PackageWidth, &p.Cost, &p.ID, &p.PackageHeight, &p.ArchivedOn, &p.OptionSummary, &p.UpdatedOn, &p.UPC)

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
        product_width,
        package_length,
        sale_price,
        description,
        package_weight,
        price,
        product_weight,
        quantity,
        product_root_id,
        product_height,
        taxable,
        brand,
        product_length,
        created_on,
        available_on,
        quantity_per_package,
        on_sale,
        name,
        sku,
        manufacturer,
        subtitle,
        package_width,
        cost,
        id,
        package_height,
        archived_on,
        option_summary,
        updated_on,
        upc
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
			&p.ProductWidth,
			&p.PackageLength,
			&p.SalePrice,
			&p.Description,
			&p.PackageWeight,
			&p.Price,
			&p.ProductWeight,
			&p.Quantity,
			&p.ProductRootID,
			&p.ProductHeight,
			&p.Taxable,
			&p.Brand,
			&p.ProductLength,
			&p.CreatedOn,
			&p.AvailableOn,
			&p.QuantityPerPackage,
			&p.OnSale,
			&p.Name,
			&p.SKU,
			&p.Manufacturer,
			&p.Subtitle,
			&p.PackageWidth,
			&p.Cost,
			&p.ID,
			&p.PackageHeight,
			&p.ArchivedOn,
			&p.OptionSummary,
			&p.UpdatedOn,
			&p.UPC,
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
        product_width,
        package_length,
        sale_price,
        description,
        package_weight,
        price,
        product_weight,
        quantity,
        product_root_id,
        product_height,
        taxable,
        brand,
        product_length,
        created_on,
        available_on,
        quantity_per_package,
        on_sale,
        name,
        sku,
        manufacturer,
        subtitle,
        package_width,
        cost,
        id,
        package_height,
        archived_on,
        option_summary,
        updated_on,
        upc
    FROM
        products
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProduct(db storage.Querier, id uint64) (*models.Product, error) {
	p := &models.Product{}

	err := db.QueryRow(productSelectionQuery, id).Scan(&p.ProductWidth, &p.PackageLength, &p.SalePrice, &p.Description, &p.PackageWeight, &p.Price, &p.ProductWeight, &p.Quantity, &p.ProductRootID, &p.ProductHeight, &p.Taxable, &p.Brand, &p.ProductLength, &p.CreatedOn, &p.AvailableOn, &p.QuantityPerPackage, &p.OnSale, &p.Name, &p.SKU, &p.Manufacturer, &p.Subtitle, &p.PackageWidth, &p.Cost, &p.ID, &p.PackageHeight, &p.ArchivedOn, &p.OptionSummary, &p.UpdatedOn, &p.UPC)

	return p, err
}

func buildProductListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"product_width",
			"package_length",
			"sale_price",
			"description",
			"package_weight",
			"price",
			"product_weight",
			"quantity",
			"product_root_id",
			"product_height",
			"taxable",
			"brand",
			"product_length",
			"created_on",
			"available_on",
			"quantity_per_package",
			"on_sale",
			"name",
			"sku",
			"manufacturer",
			"subtitle",
			"package_width",
			"cost",
			"id",
			"package_height",
			"archived_on",
			"option_summary",
			"updated_on",
			"upc",
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
			&p.ProductWidth,
			&p.PackageLength,
			&p.SalePrice,
			&p.Description,
			&p.PackageWeight,
			&p.Price,
			&p.ProductWeight,
			&p.Quantity,
			&p.ProductRootID,
			&p.ProductHeight,
			&p.Taxable,
			&p.Brand,
			&p.ProductLength,
			&p.CreatedOn,
			&p.AvailableOn,
			&p.QuantityPerPackage,
			&p.OnSale,
			&p.Name,
			&p.SKU,
			&p.Manufacturer,
			&p.Subtitle,
			&p.PackageWidth,
			&p.Cost,
			&p.ID,
			&p.PackageHeight,
			&p.ArchivedOn,
			&p.OptionSummary,
			&p.UpdatedOn,
			&p.UPC,
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
            product_width, package_length, sale_price, description, package_weight, price, product_weight, quantity, product_root_id, product_height, taxable, brand, product_length, available_on, quantity_per_package, on_sale, name, sku, manufacturer, subtitle, package_width, cost, package_height, option_summary, upc
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
        )
    RETURNING
        id, created_on, available_on;
`

func (pg *postgres) CreateProduct(db storage.Querier, nu *models.Product) (uint64, time.Time, time.Time, error) {
	var (
		createdID   uint64
		createdAt   time.Time
		availableOn time.Time
	)

	err := db.QueryRow(productCreationQuery, &nu.ProductWidth, &nu.PackageLength, &nu.SalePrice, &nu.Description, &nu.PackageWeight, &nu.Price, &nu.ProductWeight, &nu.Quantity, &nu.ProductRootID, &nu.ProductHeight, &nu.Taxable, &nu.Brand, &nu.ProductLength, &nu.AvailableOn, &nu.QuantityPerPackage, &nu.OnSale, &nu.Name, &nu.SKU, &nu.Manufacturer, &nu.Subtitle, &nu.PackageWidth, &nu.Cost, &nu.PackageHeight, &nu.OptionSummary, &nu.UPC).Scan(&createdID, &createdAt, &availableOn)
	return createdID, createdAt, availableOn, err
}

const productUpdateQuery = `
    UPDATE products
    SET
        product_width = $1, 
        package_length = $2, 
        sale_price = $3, 
        description = $4, 
        package_weight = $5, 
        price = $6, 
        product_weight = $7, 
        quantity = $8, 
        product_root_id = $9, 
        product_height = $10, 
        taxable = $11, 
        brand = $12, 
        product_length = $13, 
        available_on = $14, 
        quantity_per_package = $15, 
        on_sale = $16, 
        name = $17, 
        sku = $18, 
        manufacturer = $19, 
        subtitle = $20, 
        package_width = $21, 
        cost = $22, 
        package_height = $23, 
        option_summary = $24, 
        updated_on = NOW()
        upc = $26
    WHERE id = $26
    RETURNING updated_on;
`

func (pg *postgres) UpdateProduct(db storage.Querier, updated *models.Product) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productUpdateQuery, &updated.ProductWidth, &updated.PackageLength, &updated.SalePrice, &updated.Description, &updated.PackageWeight, &updated.Price, &updated.ProductWeight, &updated.Quantity, &updated.ProductRootID, &updated.ProductHeight, &updated.Taxable, &updated.Brand, &updated.ProductLength, &updated.AvailableOn, &updated.QuantityPerPackage, &updated.OnSale, &updated.Name, &updated.SKU, &updated.Manufacturer, &updated.Subtitle, &updated.PackageWidth, &updated.Cost, &updated.PackageHeight, &updated.OptionSummary, &updated.UPC, &updated.ID).Scan(&t)
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
