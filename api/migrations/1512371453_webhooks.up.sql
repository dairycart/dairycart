CREATE TYPE webhook_event AS ENUM ('product_created', 'product_updated', 'product_archived');
CREATE TYPE content_type AS ENUM ('application/json', 'application/xml');

CREATE TABLE IF NOT EXISTS webhooks (
    "id" bigserial,
    "url" text NOT NULL,
    "event_type" webhook_event NOT NULL,
    "content_type" content_type NOT NULL DEFAULT 'application/json',
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS webhook_execution_logs (
    "id" bigserial,
    "webhook_id" bigint NOT NULL,
    "status_code" int NOT NULL,
    "succeeded" boolean NOT NULL DEFAULT 'false',
    "executed_on" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("webhook_id") REFERENCES "webhooks"("id")
);