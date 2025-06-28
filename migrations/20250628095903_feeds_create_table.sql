-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id         INTEGER PRIMARY KEY,
    title      TEXT NOT NULL,
    subtitle   TEXT,
    feed_url   TEXT NOT NULL UNIQUE,
    site_url   TEXT NOT NULL,
    type       TEXT NOT NULL,
    disabled   INTEGER DEFAULT 0 NOT NULL,
    checked_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
