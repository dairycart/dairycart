package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

const productQueryBySKU = `
    SELECT
        id,
        product_root_id,
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

func (pg *Postgres) GetProductBySKU(sku string) (*models.Product, error) {
	p := &models.Product{}

	err := pg.DB.QueryRow(productQueryBySKU, sku).Scan(&p.ID, &p.ProductRootID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

const productWithSKUExistenceQuery = `SELECT EXISTS(SELECT id FROM products WHERE sku = $1 and archived_on IS NULL);`

func (pg *Postgres) ProductWithSKUExists(sku string) (bool, error) {
	var exists string

	err := pg.DB.QueryRow(productWithSKUExistenceQuery, sku).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productExistenceQuery = `SELECT EXISTS(SELECT id FROM products WHERE id = $1 and archived_on IS NULL);`

func (pg *Postgres) ProductExists(id uint64) (bool, error) {
	var exists string

	err := pg.DB.QueryRow(productExistenceQuery, id).Scan(&exists)
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

func (pg *Postgres) GetProduct(id uint64) (*models.Product, error) {
	p := &models.Product{}

	err := pg.DB.QueryRow(productSelectionQuery, id).Scan(&p.ID, &p.ProductRootID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

const productCreationQuery = `
    INSERT INTO products
        (
            product_root_id, name, subtitle, description, option_summary, sku, upc, manufacturer, brand, quantity, taxable, price, on_sale, sale_price, cost, product_weight, product_height, product_width, product_length, package_weight, package_height, package_width, package_length, quantity_per_package, available_on
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
        )
    RETURNING
        id, created_on, available_on;
`

func (pg *Postgres) CreateProduct(nu *models.Product) (uint64, time.Time, time.Time, error) {
	var (
		createdID   uint64
		createdAt   time.Time
		availableOn time.Time
	)
	err := pg.DB.QueryRow(productCreationQuery, &nu.ProductRootID, &nu.Name, &nu.Subtitle, &nu.Description, &nu.OptionSummary, &nu.SKU, &nu.UPC, &nu.Manufacturer, &nu.Brand, &nu.Quantity, &nu.Taxable, &nu.Price, &nu.OnSale, &nu.SalePrice, &nu.Cost, &nu.ProductWeight, &nu.ProductHeight, &nu.ProductWidth, &nu.ProductLength, &nu.PackageWeight, &nu.PackageHeight, &nu.PackageWidth, &nu.PackageLength, &nu.QuantityPerPackage, &nu.AvailableOn).Scan(&createdID, &createdAt, &availableOn)

	return createdID, createdAt, availableOn, err
}

const productUpdateQuery = `
    UPDATE products
    SET
        product_root_id = $1, 
        name = $2, 
        subtitle = $3, 
        description = $4, 
        option_summary = $5, 
        sku = $6, 
        upc = $7, 
        manufacturer = $8, 
        brand = $9, 
        quantity = $10, 
        taxable = $11, 
        price = $12, 
        on_sale = $13, 
        sale_price = $14, 
        cost = $15, 
        product_weight = $16, 
        product_height = $17, 
        product_width = $18, 
        product_length = $19, 
        package_weight = $20, 
        package_height = $21, 
        package_width = $22, 
        package_length = $23, 
        quantity_per_package = $24, 
        available_on = $25, 
        updated_on = NOW()
    WHERE id = $26
    RETURNING updated_on;
`

func (pg *Postgres) UpdateProduct(updated *models.Product) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productUpdateQuery, &updated.ProductRootID, &updated.Name, &updated.Subtitle, &updated.Description, &updated.OptionSummary, &updated.SKU, &updated.UPC, &updated.Manufacturer, &updated.Brand, &updated.Quantity, &updated.Taxable, &updated.Price, &updated.OnSale, &updated.SalePrice, &updated.Cost, &updated.ProductWeight, &updated.ProductHeight, &updated.ProductWidth, &updated.ProductLength, &updated.PackageWeight, &updated.PackageHeight, &updated.PackageWidth, &updated.PackageLength, &updated.QuantityPerPackage, &updated.AvailableOn, &updated.ID).Scan(&t)
	return t, err
}

const productDeletionQuery = `
    UPDATE products
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteProduct(id uint64) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productDeletionQuery, id).Scan(&t)
	return t, err
}