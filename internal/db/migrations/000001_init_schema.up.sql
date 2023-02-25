CREATE TABLE "accounts"
(
    "id"         bigserial   PRIMARY KEY,
    "owner_id"   bigint      NOT NULL,
    "balance"    bigint      NOT NULL,
    "currency"   varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries"
(
    "id"         bigserial   PRIMARY KEY,
    "account_id" bigint      NOT NULL,
    "amount"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers"
(
    "id"           bigserial   PRIMARY KEY,
    "sender_id"    bigint      NOT NULL,
    "recipient_id" bigint      NOT NULL,
    "amount"       bigint      NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("sender_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("recipient_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("owner_id");

CREATE INDEX ON "transfers" ("sender_id");

CREATE INDEX ON "transfers" ("recipient_id");

CREATE INDEX ON "transfers" ("sender_id", "recipient_id");