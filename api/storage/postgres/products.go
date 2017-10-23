package postgres

import (
	"github.com/dairycart/dairycart/api/storage/models"
)

func (pg Postgres) GetProductBySKU(sku string) (models.Product, error) {
	var p models.Product

	query := `
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

        FROM products
        WHERE sku = $1
    `

	err := pg.DB.QueryRow(query, sku).Scan(&p.ID, &p.ProductRootID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)
	return p, err
}

func (pg Postgres) GetProductByID(ID uint64) (models.Product, error) {
	var p models.Product
	query := `
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

        FROM products
        WHERE id = $1
    `

	err := pg.DB.QueryRow(query, ID).Scan(&p.ID, &p.ProductRootID, &p.Name, &p.Subtitle, &p.Description, &p.OptionSummary, &p.SKU, &p.UPC, &p.Manufacturer, &p.Brand, &p.Quantity, &p.Taxable, &p.Price, &p.OnSale, &p.SalePrice, &p.Cost, &p.ProductWeight, &p.ProductHeight, &p.ProductWidth, &p.ProductLength, &p.PackageWeight, &p.PackageHeight, &p.PackageWidth, &p.PackageLength, &p.QuantityPerPackage, &p.AvailableOn, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)
	return p, err
}
