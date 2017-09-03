CREATE TABLE IF NOT EXISTS product_roots (
    "id" bigserial,
    "name" text NOT NULL,
    "subtitle" text NOT NULL DEFAULT '',
    "description" text NOT NULL DEFAULT '',
    "sku_prefix" text NOT NULL,
    "manufacturer" text NOT NULL DEFAULT '',
    "brand" text NOT NULL DEFAULT '',
    "taxable" boolean DEFAULT 'false',
    "cost" numeric(15, 2) NOT NULL,
    "product_weight" numeric(15, 2),
    "product_height" numeric(15, 2),
    "product_width" numeric(15, 2),
    "product_length" numeric(15, 2),
    "package_weight" numeric(15, 2),
    "package_height" numeric(15, 2),
    "package_width" numeric(15, 2),
    "package_length" numeric(15, 2),
    "quantity_per_package" integer DEFAULT 1,
    "available_on" timestamp NOT NULL DEFAULT NOW(),
    "created_on" timestamp DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("sku_prefix", "archived_on"),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS products (
    "id" bigserial,
    "product_root_id" bigint NOT NULL,
    "name" text NOT NULL,
    "subtitle" text NOT NULL DEFAULT '',
    "description" text NOT NULL DEFAULT '',
    "option_summary" text NOT NULL DEFAULT '',
    "sku" text NOT NULL,
    "upc" text NOT NULL DEFAULT '',
    "manufacturer" text NOT NULL DEFAULT '',
    "brand" text NOT NULL DEFAULT '',
    "quantity" integer NOT NULL DEFAULT 0,
    "taxable" boolean DEFAULT 'false',
    "price" numeric(15, 2) NOT NULL,
    "on_sale" boolean DEFAULT 'false',
    "sale_price" numeric(15, 2) NOT NULL DEFAULT 0 CONSTRAINT sale_price_must_not_be_zero CHECK(
        (sale_price != 0 AND on_sale IS TRUE)
                         OR
        (sale_price  = 0 AND on_sale IS FALSE)
    ),
    "cost" numeric(15, 2) NOT NULL,
    "product_weight" numeric(15, 2),
    "product_height" numeric(15, 2),
    "product_width" numeric(15, 2),
    "product_length" numeric(15, 2),
    "package_weight" numeric(15, 2),
    "package_height" numeric(15, 2),
    "package_width" numeric(15, 2),
    "package_length" numeric(15, 2),
    "quantity_per_package" integer DEFAULT 1,
    "available_on" timestamp DEFAULT NOW(),
    "created_on" timestamp DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("sku", "archived_on"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_root_id") REFERENCES "product_roots"("id")
);
CREATE UNIQUE INDEX products_upc_empty_but_not_null_idx ON products (upc) WHERE upc != '';

CREATE TABLE IF NOT EXISTS product_options (
    "id" bigserial,
    "name" text NOT NULL,
    "product_root_id" bigint NOT NULL,
    "created_on" timestamp DEFAULT NOW(),
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
    "created_on" timestamp DEFAULT NOW(),
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
    "created_on" timestamp DEFAULT NOW(),
    "archived_on" timestamp,
    FOREIGN KEY ("product_id") REFERENCES "products"("id"),
    FOREIGN KEY ("product_option_value_id") REFERENCES "product_option_values"("id")
);