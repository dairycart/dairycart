CREATE TYPE discount_type AS ENUM ('percentage', 'flat_amount');
CREATE TABLE sales_and_discounts (
    "id" bigserial,
    "name" text,
    "type" discount_type,
    "amount" numeric(7, 2),
    "product_id" bigint,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products"("id")
);