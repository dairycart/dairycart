CREATE TYPE discount_type AS ENUM ('percentage', 'flat_amount');
CREATE TABLE discounts (
    "id" bigserial,
    "name" text NOT NULL,
    "type" discount_type,
    "amount" numeric(7, 2),
    "product_id" bigint,
    "starts_on" timestamp NOT NULL,
    "expires_on" timestamp NOT NULL,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp,
    "archived_at" timestamp,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products"("id")
);