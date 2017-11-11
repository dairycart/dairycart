CREATE TABLE IF NOT EXISTS users (
    "id" bigserial,
    "first_name" text NOT NULL,
    "last_name" text NOT NULL,
    "username" text NOT NULL,
    "email" text NOT NULL,
    "password" text NOT NULL,
    "salt" bytea NOT NULL,
    "is_admin" bool NOT NULL DEFAULT 'false',
    "password_last_changed_on" timestamp,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    UNIQUE ("username", "archived_on"),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    "id" bigserial,
    "user_id" bigint NOT NULL,
    "token" text NOT NULL,
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "expires_on" timestamp NOT NULL DEFAULT NOW() + (15 * interval '1 minute'),
    "password_reset_on" timestamp,
    UNIQUE ("token"),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS login_attempts (
    "id" bigserial,
    "username" text NOT NULL,
    "successful" boolean NOT NULL DEFAULT 'false',
    "created_on" timestamp NOT NULL DEFAULT NOW()
)