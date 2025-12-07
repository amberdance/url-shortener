-- migrate:up
CREATE TABLE IF NOT EXISTS urls
(
    id             uuid,
    hash           VARCHAR(255) UNIQUE NOT NULL,
    original_url   TEXT                NOT NULL,
    created_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ         NULL,
    correlation_id varchar(255),
    PRIMARY KEY (id),
    UNIQUE (hash),
    UNIQUE (correlation_id),
    UNIQUE (original_url)
);

-- migrate:down
DROP TABLE IF EXISTS urls;
