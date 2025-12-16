-- migrate:up
ALTER TABLE urls
    ADD user_id uuid;

CREATE INDEX user_id_idx ON urls (user_id);

-- migrate:down
ALTER TABLE urls
    DROP COLUMN IF EXISTS user_id;

DROP INDEX user_id_idx;
