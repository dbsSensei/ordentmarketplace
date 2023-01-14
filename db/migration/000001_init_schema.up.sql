CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "orders"
(
    "id"         serial PRIMARY KEY,
    "user_id"    int   NOT NULL,
    "status"     varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);


CREATE TABLE "order_items"
(
    "order_id"   int NOT NULL,
    "product_id" int NOT NULL,
    "quantity"   int NOT NULL
);

CREATE TABLE "products"
(
    "id"          serial PRIMARY KEY,
    "name"        varchar     NOT NULL,
    "merchant_id" int         NOT NULL,
    "price"       int         NOT NULL,
    "status"      varchar     NOT NULL,
    "category_id" int         NOT NULL,
    "created_at"  timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users"
(
    "id"                  serial PRIMARY KEY,
    "full_name"           varchar        NOT NULL,
    "email"               varchar UNIQUE NOT NULL,
    "country_code"        varchar        NOT NULL,
    "hashed_password"     varchar        NOT NULL,
    "password_changed_at" timestamptz    NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
    "created_at"          timestamptz    NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "user_id"       int              NOT NULL,
    "refresh_token" varchar          NOT NULL,
    "user_agent"    varchar          NOT NULL,
    "client_ip"     varchar          NOT NULL,
    "is_blocked"    boolean          NOT NULL DEFAULT false,
    "expires_at"    timestamptz      NOT NULL,
    "created_at"    timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "merchants"
(
    "id"           serial PRIMARY KEY,
    "admin_id"     int UNIQUE  NOT NULL,
    "name"         varchar     NOT NULL,
    "country_code" varchar     NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories"
(
    "id"        serial PRIMARY KEY,
    "name"      varchar UNIQUE NOT NULL,
    "parent_id" INT            DEFAULT 0
);

CREATE TABLE "countries"
(
    "code" varchar PRIMARY KEY,
    "name" varchar NOT NULL
);

ALTER TABLE "orders"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_items"
    ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items"
    ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "products"
    ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "products"
    ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "categories"
    ADD FOREIGN KEY ("parent_id") REFERENCES "categories" ("id");

ALTER TABLE "users"
    ADD FOREIGN KEY ("country_code") REFERENCES "countries" ("code");

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "merchants"
    ADD FOREIGN KEY ("admin_id") REFERENCES "users" ("id");

ALTER TABLE "merchants"
    ADD FOREIGN KEY ("country_code") REFERENCES "countries" ("code");
