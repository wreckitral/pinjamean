-- 000001_create_users.up.sql
CREATE TABLE users (
    uuid            UUID PRIMARY KEY,
    email           TEXT NOT NULL,
    passwordHashed  TEXT NOT NULL,
    role            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL
);
