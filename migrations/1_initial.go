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
			"taxable" boolean DEFAULT 'false',
			"customer_can_set_pricing" boolean default 'false',
			PRIMARY KEY ("id")
		);

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
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping products table...")
		_, err := db.Exec(`DROP TABLE products`)
		return err
	})
}
