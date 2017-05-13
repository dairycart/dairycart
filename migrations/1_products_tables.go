package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating products tables...")
		_, err := db.Exec(`CREATE TABLE base_products (
			"id" bigserial,
			"name" text,
			"description" text,
			"on_sale" boolean DEFAULT 'false',
			"customer_can_set_pricing" boolean default 'false',
			"base_price" decimal,
			"base_sale_price" decimal,
			"base_taxable" boolean DEFAULT 'false',
			"base_product_weight" decimal,
			"base_product_height" decimal,
			"base_product_width" decimal,
			"base_product_length" decimal,
			"base_package_weight" decimal,
			"base_package_height" decimal,
			"base_package_width" decimal,
			"base_package_length" decimal,
			"active" boolean DEFAULT 'true',
			"created" timestamp DEFAULT NOW(),
			"archived" timestamp,
			PRIMARY KEY ("id")
		);`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`CREATE TABLE products (
			"id" bigserial,
			"base_product_id" bigint NOT NULL,
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
    		UNIQUE ("sku"),
    		UNIQUE ("upc"),
			PRIMARY KEY ("id"),
			FOREIGN KEY ("base_product_id") REFERENCES "base_products"("id")
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
		fmt.Println("products tables created!")

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
		if err != nil {
			return err
		}

		_, err = db.Exec(`DROP TABLE base_products`)
		return err
	})
}
