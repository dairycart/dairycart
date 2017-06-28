CREATE TABLE IF NOT EXISTS users (
    "id" bigserial,
    "first_name" text NOT NULL,
    "last_name" text NOT NULL,
    "email" text NOT NULL,
    "password" text NOT NULL,
    "salt" bytea NOT NULL,
    "is_admin" bool DEFAULT 'false',
    "created_on" timestamp DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    PRIMARY KEY ("id")
);