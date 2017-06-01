CREATE TABLE product_progenitors (
    "id" bigserial,
    "name" text NOT NULL,
    "description" text NOT NULL,
    "taxable" boolean DEFAULT 'false',
    "price" numeric(7, 2) NOT NULL,
    "cost" numeric(7, 2) NOT NULL,
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
    "product_progenitor_id" bigint,
    "sku" text NOT NULL,
    "name" text NOT NULL,
    "upc" text,
    "quantity" integer,
    "price" numeric(7, 2) NOT NULL,
    "cost" numeric(7, 2) NOT NULL,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    UNIQUE ("sku"),
    UNIQUE ("product_progenitor_id", "name"),
    UNIQUE ("upc"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_progenitor_id") REFERENCES "product_progenitors"("id")
);

CREATE TABLE product_attributes (
    "id" bigserial,
    "name" text NOT NULL,
    "product_progenitor_id" bigint,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    UNIQUE ("product_progenitor_id", "name"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_progenitor_id") REFERENCES "product_progenitors"("id")
);

CREATE TABLE product_attribute_values (
    "id" bigserial,
    "product_attribute_id" bigint,
    "value" text NOT NULL,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    UNIQUE ("product_attribute_id", "value"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_attribute_id") REFERENCES "product_attributes"("id")
);