CREATE TABLE "users"
(
    "id"                  bigserial PRIMARY KEY,
    "username"            varchar NOT NULL UNIQUE,
    "email"               varchar NOT NULL UNIQUE,
    "password"            varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner_id", "currency");