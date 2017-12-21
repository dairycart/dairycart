CREATE TYPE discount_type AS ENUM ('percentage', 'flat_amount');
CREATE TABLE IF NOT EXISTS discounts (
    "id" bigserial,
    "name" text NOT NULL,
    "discount_type" discount_type NOT NULL DEFAULT 'percentage',
    "amount" numeric(15, 2) NOT NULL DEFAULT 0,
    "expires_on" timestamp,
    "requires_code" boolean NOT NULL DEFAULT FALSE,
    "code" text NOT NULL DEFAULT '' CONSTRAINT code_must_be_provided CHECK(
        (code != '' AND requires_code IS TRUE)
                        OR
        (code = '' AND requires_code IS FALSE)
    ),
    "limited_use" boolean NOT NULL DEFAULT FALSE,
    "number_of_uses" bigint NOT NULL DEFAULT 0 CONSTRAINT use_number_must_be_provided CHECK(
        (number_of_uses != 0 AND limited_use IS TRUE)
                            OR
        (number_of_uses = 0 AND limited_use IS FALSE)
    ),
    "login_required" boolean NOT NULL DEFAULT FALSE,
    "starts_on" timestamp NOT NULL DEFAULT NOW(),
    "created_on" timestamp NOT NULL DEFAULT NOW(),
    "updated_on" timestamp,
    "archived_on" timestamp,
    PRIMARY KEY ("id")
);