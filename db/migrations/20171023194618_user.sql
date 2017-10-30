
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TABLE "user";

CREATE TABLE "user" (
    id SERIAL,

    username text UNIQUE NOT NULL,
    email text,

    password_salt bytea NOT NULL,
    -- PBKDF2
    password_iterations integer NOT NULL,
    password_hash bytea NOT NULL,

    PRIMARY KEY(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "user";
