CREATE TABLE "users"
(
    "id"          uuid PRIMARY KEY,
    "name"        varchar     NOT NULL,
    "email"       varchar     NOT NULL UNIQUE,
    "password"    varchar     NOT NULL,
    "avatar"      text                 DEFAULT NULL,
    "verified_at" timestamptz          DEFAULT NULL,
    "created_at"  timestamptz NOT NULL DEFAULT (now()),
    "deleted_at"  timestamptz          DEFAULT NULL
);

CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "user_id"       uuid        NOT NULL,
    "access_token"  varchar     NOT NULL,
    "refresh_token" varchar     NOT NULL,
    "user_agent"    varchar     NOT NULL,
    "client_ip"     varchar     NOT NULL,
    "is_blocked"    boolean     NOT NULL DEFAULT false,
    "expires_at"    timestamptz NOT NULL,
    "created_at"    timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id")
        ON UPDATE CASCADE
        ON DELETE CASCADE;

CREATE TABLE "roles"
(
    "id"         bigserial PRIMARY KEY,
    "name"       varchar     NOT NULL UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_role"
(
    "user_id"    uuid        NOT NULL,
    "role_id"    bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    UNIQUE (user_id, role_id)
);

ALTER TABLE "user_role"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id")
        ON UPDATE CASCADE
        ON DELETE CASCADE;

ALTER TABLE "user_role"
    ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id")
        ON UPDATE CASCADE
        ON DELETE CASCADE;

INSERT INTO "roles" (name)
VALUES ('root'),
       ('admin'),
       ('normal');

