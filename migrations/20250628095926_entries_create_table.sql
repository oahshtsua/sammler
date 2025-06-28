-- +goose Up
-- +goose StatementBegin
CREATE TABLE entries (
    id            INTEGER PRIMARY KEY,
    feed_id       INTEGER NOT NULL,
    title         TEXT NOT NULL,
    author        TEXT,
    content       TEXT NOT NULL,
    external_url  TEXT NOT NULL,
    published_at  TEXT NOT NULL,
    read          INTEGER DEFAULT 0 NOT NULL,
    starred       INTEGER DEFAULT 0 NOT NULL,
    created_at    TEXT NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE entries;
-- +goose StatementEnd

