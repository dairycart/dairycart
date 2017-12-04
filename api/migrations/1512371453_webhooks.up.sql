CREATE TYPE webhook_event AS ENUM ('product_created', 'product_updated', 'product_archived');
CREATE TABLE IF NOT EXISTS webhooks (
    "id" bigserial,
    "url" text NOT NULL,
    "event_type" webhook_event NOT NULL,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    PRIMARY KEY ("id")
);