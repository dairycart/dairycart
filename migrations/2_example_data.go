package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(`
			INSERT INTO base_products (
				"name",
				"description",
				"base_price",
				"base_product_weight",
				"base_product_height",
				"base_product_width",
				"base_product_length",
				"base_package_weight",
				"base_package_height",
				"base_package_width",
				"base_package_length"
			) VALUES
			('T-Shirt', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', 12.34, 1, 2, 3, 4, 5, 6, 7, 8),
			('Skateboard', 'This is a skateboard. Please wear a helmet.', 99.99, 8, 7, 6, 5, 4, 3, 2, 1);
		`)
		if err != nil {
			return err
		}

		fmt.Println("creating example products...")
		_, err = db.Exec(`
			INSERT INTO products(
				"base_product_id",
				"sku",
				"name",
				"description",
				"upc",
				"price",
				"on_sale",
				"sale_price",
				"product_weight",
				"product_height",
				"product_width",
				"product_length",
				"package_weight",
				"package_height",
				"package_width",
				"package_length",
				"quantity")
			VALUES
			(1, 't-shirt-small-red', 'Red T-Shirt (Small)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-small-green', 'Green T-Shirt (Small)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-small-blue', 'Blue T-Shirt (Small)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-medium-red', 'Red T-Shirt (Medium)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-medium-green', 'Green T-Shirt (Medium)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-medium-blue', 'Blue T-Shirt (Medium)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-large-red', 'Red T-Shirt (Large)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-large-green', 'Green T-Shirt (Large)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(1, 't-shirt-large-blue', 'Blue T-Shirt (Large)', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, true, 4.2, 1, 2, 3, 4, 5, 6, 7, 8, 123),
			(2, 'skateboard', 'Skateboard', 'This is a skateboard. Please wear a helmet.', '9876543210', 99.99, false, null, 8, 7, 6, 5, 4, 3, 2, 1, 321);
		`)
		if err != nil {
			return err
		}

		fmt.Println("creating example product attributes...")
		_, err = db.Exec(`
			INSERT INTO product_attributes ("name", "base_product_id") VALUES ('Color', 1), ('Size', 1);
		`)

		if err != nil {
			return err
		}

		fmt.Println("creating example product attribute values...")
		_, err = db.Exec(`
			INSERT INTO product_attribute_values ("product_attribute_id", "value") VALUES(1, 'Red'), (1, 'Blue'), (1, 'Green'), (2, 'Small'), (2, 'Medium'), (2, 'Large');
		`)

		fmt.Println("example data created!")
		return err
	}, func(db migrations.DB) error {

		fmt.Println("removing example product attribute values...")
		_, err := db.Exec(`delete from product_attribute_values where id is not null;`)
		if err != nil {
			return err
		}

		fmt.Println("removing example product attributes...")
		_, err = db.Exec(`delete from product_attributes where id is not null;`)
		if err != nil {
			return err
		}

		fmt.Println("removing example products...")
		_, err = db.Exec(`delete from products where id is not null;`)
		if err != nil {
			return err
		}

		return err
	})
}
