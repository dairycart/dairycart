CREATE TABLE IF NOT EXISTS users (
    "id" bigserial,
    "first_name" text NOT NULL,
    "last_name" text NOT NULL,
    "username" text NOT NULL,
    "email" text NOT NULL,
    "password" text NOT NULL,
    "salt" bytea NOT NULL,
    "is_admin" bool DEFAULT 'false',
    "password_last_changed_on" timestamp,
    "created_on" timestamp DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("username"),
    PRIMARY KEY ("id")
);