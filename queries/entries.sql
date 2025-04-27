-- name: CreateEntry :exec
INSERT INTO entries (
    feed_id,
    title,
    author,
    content,
    external_url,
    published_at,
    created_at
)
VALUES (
    ?, ?, ?, ?, ?, ?, ?
);

