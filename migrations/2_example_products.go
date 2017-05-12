package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		// create example product
		_, err := db.Exec(`
			INSERT INTO products(
				"sku",
				"name",
				"description",
				"upc",
				"price",
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
			(
				't-shirt',
				'T-Shirt',
				'This is a t-shirt. wear it, or don''t. I''m not your dad.',
				'1234567890',
				12.34,
				4.2,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				123
			),
			(
				'skateboard',
				'Skateboard',
				'This is a skateboard. Please wear a helmet.',
				'9876543210',
				99.99,
				null,
				8,
				7,
				6,
				5,
				4,
				3,
				2,
				1,
				321
			);
		`)
		if err != nil {
			return err
		}

		// create example product attributes
		_, err = db.Exec(`
			INSERT INTO product_attributes ("name", "assigned_to_product") VALUES ('Color', 1), ('Size', 1);
		`)

		if err != nil {
			return err
		}

		// create example product attribute values
		_, err = db.Exec(`
			INSERT INTO product_attribute_values ("parent_attribute", "value") VALUES(1, 'Red'), (1, 'Blue'), (1, 'Green'), (2, 'Small'), (2, 'Medium'), (2, 'Large');
		`)

		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("")
		return err
	})
}
