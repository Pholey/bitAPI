
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE user_channel (
  id          SERIAL REFERENCES public.user ON DELETE CASCADE,
  channel_name text UNIQUE NOT NULL,

  PRIMARY KEY(channel_name)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user_channel;
