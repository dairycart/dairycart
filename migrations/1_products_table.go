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
			"on_sale" boolean DEFAULT 'false',
			"price" decimal,
			"weight" decimal,
			"height" decimal,
			"width" decimal,
			"length" decimal,
			"quantity" integer,
			"child_of" bigint DEFAULT null,
			"taxable" boolean DEFAULT 'false',
			"customer_can_set_pricing" boolean default 'false',
			PRIMARY KEY ("id")
		);`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE products ADD FOREIGN KEY ("child_of") REFERENCES "products"("id");`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO products (
				"sku",
				"name",
				"description",
				"upc",
				"price",
				"weight",
				"height",
				"width",
				"length",
				"quantity") VALUES (
				'example-sku',
				'example product',
				'this is a test product',
				'1234567890',
				12.34,
				1,
				2,
				3,
				4,
				123);
		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO products (
				"sku",
				"name",
				"description",
				"upc",
				"price",
				"weight",
				"height",
				"width",
				"length",
				"quantity",
				"child_of") VALUES (
				'example-child',
				'example child product',
				'this is a test product child',
				'0987654321',
				12.34,
				9,
				7,
				8,
				6,
				234,
				1);
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping products table...")
		_, err := db.Exec(`DROP TABLE products`)
		return err
	})
}
