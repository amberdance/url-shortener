-- migrate:up
CREATE TABLE IF NOT EXISTS urls
(
    id           uuid,
    hash         VARCHAR(255) UNIQUE NOT NULL,
    original_url TEXT                NOT NULL,
    created_at   TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ         NULL,
    PRIMARY KEY (id),
    UNIQUE (hash)
);

-- migrate:down
DROP TABLE IF EXISTS urls;
