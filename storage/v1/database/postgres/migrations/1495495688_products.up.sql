CREATE TABLE IF NOT EXISTS product_roots (
    "id" bigserial,
    "name" text NOT NULL,
    "primary_image_id" bigint,
    "subtitle" text NOT NULL DEFAULT '',
    "description" text NOT NULL DEFAULT '',
    "sku_prefix" text NOT NULL,
    "manufacturer" text NOT NULL DEFAULT '',
    "brand" text NOT NULL DEFAULT '',
    "taxable" boolean NOT NULL DEFAULT 'false',
    "cost" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_weight" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_height" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_width" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_length" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_weight" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_height" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_width" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_length" numeric(15, 2) NOT NULL DEFAULT 0,
    "quantity_per_package" integer NOT NULL DEFAULT 1,
    "available_on" timestamp NOT NULL DEFAULT NOW(),
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("sku_prefix", "archived_on"),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS products (
    "id" bigserial,
    "product_root_id" bigint NOT NULL,
    "primary_image_id" bigint,
    "name" text NOT NULL,
    "subtitle" text NOT NULL DEFAULT '',
    "description" text NOT NULL DEFAULT '',
    "option_summary" text NOT NULL DEFAULT '',
    "sku" text NOT NULL,
    "upc" text NOT NULL DEFAULT '',
    "manufacturer" text NOT NULL DEFAULT '',
    "brand" text NOT NULL DEFAULT '',
    "quantity" integer NOT NULL DEFAULT 0,
    "taxable" boolean NOT NULL DEFAULT 'false',
    "price" numeric(15, 2) NOT NULL,
    "on_sale" boolean NOT NULL DEFAULT 'false',
    "sale_price" numeric(15, 2) NOT NULL DEFAULT 0 CONSTRAINT sale_price_must_not_be_zero CHECK(
        (sale_price != 0 AND on_sale IS TRUE)
                         OR
        (sale_price  = 0 AND on_sale IS FALSE)
    ),
    "cost" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_weight" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_height" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_width" numeric(15, 2) NOT NULL DEFAULT 0,
    "product_length" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_weight" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_height" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_width" numeric(15, 2) NOT NULL DEFAULT 0,
    "package_length" numeric(15, 2) NOT NULL DEFAULT 0,
    "quantity_per_package" integer NOT NULL DEFAULT 1,
    "available_on" timestamp NOT NULL DEFAULT NOW(),
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("sku", "archived_on"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_root_id") REFERENCES "product_roots"("id")
);
CREATE UNIQUE INDEX products_upc_empty_but_not_null_idx ON products (upc) WHERE upc != '';

CREATE TABLE IF NOT EXISTS product_images (
    "id" bigserial,
    "product_root_id" bigint NOT NULL,
    "thumbnail_url" text NOT NULL,
    "main_url" text NOT NULL,
    "original_url" text NOT NULL,
    "source_url" text NOT NULL DEFAULT '',
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_root_id") REFERENCES "product_roots"("id")
);

ALTER TABLE IF EXISTS "product_roots"
    ADD FOREIGN KEY ("primary_image_id") REFERENCES "product_images"("id");

ALTER TABLE IF EXISTS "products"
    ADD FOREIGN KEY ("primary_image_id") REFERENCES "product_images"("id");

CREATE TABLE IF NOT EXISTS product_image_bridge (
    "id" bigserial,
    "product_id" bigint NOT NULL,
    "product_image_id" bigint NOT NULL,
    UNIQUE ("product_id", "product_image_id"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products"("id"),
    FOREIGN KEY ("product_image_id") REFERENCES "product_images"("id")
);

CREATE TABLE IF NOT EXISTS product_options (
    "id" bigserial,
    "name" text NOT NULL,
    "product_root_id" bigint NOT NULL,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("product_root_id", "name"),
    UNIQUE ("name", "archived_on"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_root_id") REFERENCES "product_roots"("id")
);

CREATE TABLE IF NOT EXISTS product_option_values (
    "id" bigserial,
    "product_option_id" bigint NOT NULL,
    "value" text NOT NULL,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("product_option_id", "value"),
    UNIQUE ("value", "archived_on"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_option_id") REFERENCES "product_options"("id")
);

CREATE TABLE IF NOT EXISTS product_variant_bridge (
    "id" bigserial,
    "product_id" bigint NOT NULL,
    "product_option_value_id" bigint NOT NULL,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "archived_on" timestamp,
    FOREIGN KEY ("product_id") REFERENCES "products"("id"),
    FOREIGN KEY ("product_option_value_id") REFERENCES "product_option_values"("id")
);