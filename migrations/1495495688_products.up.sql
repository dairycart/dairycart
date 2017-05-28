CREATE TABLE product_progenitors (
    "id" bigserial,
    "name" text,
    "description" text,
    "taxable" boolean DEFAULT 'false',
    "price" decimal,
    "product_weight" numeric(5, 2),
    "product_height" numeric(5, 2),
    "product_width" numeric(5, 2),
    "product_length" numeric(5, 2),
    "package_weight" numeric(5, 2),
    "package_height" numeric(5, 2),
    "package_width" numeric(5, 2),
    "package_length" numeric(5, 2),
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    PRIMARY KEY ("id")
);

CREATE TABLE products (
    "id" bigserial,
    "product_progenitor_id" bigint NOT NULL,
    "sku" text,
    "name" text,
    "upc" text,
    "quantity" integer,
    "on_sale" boolean DEFAULT 'false',
    "price" numeric(7, 2),
    "sale_price" numeric(7, 2),
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    UNIQUE ("sku"),
    UNIQUE ("upc"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_progenitor_id") REFERENCES "product_progenitors"("id")
);

CREATE TABLE product_attributes (
    "id" bigserial,
    "name" text,
    "product_progenitor_id" bigint,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_progenitor_id") REFERENCES "product_progenitors"("id")
);

CREATE TABLE product_attribute_values (
    "id" bigserial,
    "product_attribute_id" bigint,
    "value" text,
    "products_created" boolean,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_attribute_id") REFERENCES "product_attributes"("id")
);