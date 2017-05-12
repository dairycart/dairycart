package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating products table...")
		_, err := db.Exec(`CREATE TABLE products (
			"id" bigserial,
			"sku" text,
			"name" text,
			"description" text,
			"upc" text,
			"quantity" integer,
			"on_sale" boolean DEFAULT 'false',
			"price" decimal,
			"sale_price" decimal,
			"taxable" boolean DEFAULT 'false',
			"customer_can_set_pricing" boolean default 'false',
			"product_weight" decimal,
			"product_height" decimal,
			"product_width" decimal,
			"product_length" decimal,
			"package_weight" decimal,
			"package_height" decimal,
			"package_width" decimal,
			"package_length" decimal,
			"active" boolean DEFAULT 'true',
			"created" timestamp DEFAULT NOW(),
			"archived" timestamp,
			PRIMARY KEY ("id")
		);`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`CREATE TABLE product_attributes (
			"id" bigserial,
			"name" text,
			"assigned_to_product" bigint,
			"active" boolean DEFAULT 'true',
			"created" timestamp DEFAULT NOW(),
			"archived" timestamp,
			PRIMARY KEY ("id"),
			FOREIGN KEY ("assigned_to_product") REFERENCES "products"("id")
		);`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`CREATE TABLE product_attribute_values (
			"id" serial,
			"parent_attribute" bigint,
			"value" text,
			"active" boolean DEFAULT 'true',
			"created" timestamp DEFAULT NOW(),
			"archived" timestamp,
			PRIMARY KEY ("id"),
			FOREIGN KEY ("parent_attribute") REFERENCES "product_attributes"("id")
		);`)
		if err != nil {
			return err
		}

		// create example product
		_, err = db.Exec(`
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
		fmt.Println("dropping products tables...")
		_, err := db.Exec(`DROP TABLE product_attribute_values`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`DROP TABLE product_attributes`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`DROP TABLE products`)
		return err
	})
}
