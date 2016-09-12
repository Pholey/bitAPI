
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE "user" (
    id SERIAL,

    username text UNIQUE,
    email text,

    -- Why?
    name text,

    -- Our hash must be salty
    -- so we avoid lookup/rainbow table attacks
    password_salt bytea,
    -- We're gonna use PBKDF2
    -- that way hashing is so computatianally intensive
    -- that dictionary attacks / brute-forcing is nearly
    -- impossible
    -- Varying the number of iterations by some random
    -- amount helps somehow
    password_iterations integer,
    -- Don't store the password, just the hash
    password_hash bytea,

    -- dafaq is body?
    PRIMARY KEY(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "user";
