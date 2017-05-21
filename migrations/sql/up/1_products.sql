CREATE TABLE product_progenitors (
    "id" bigserial,
    "name" text,
    "description" text,
    "taxable" boolean DEFAULT 'false',
    "price" decimal,
    "product_weight" decimal,
    "product_height" decimal,
    "product_width" decimal,
    "product_length" decimal,
    "package_weight" decimal,
    "package_height" decimal,
    "package_width" decimal,
    "package_length" decimal,
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
    "price" decimal,
    "sale_price" decimal,
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