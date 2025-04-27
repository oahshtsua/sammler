-- name: GetUnreadEntries :many
SELECT feeds.title AS feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE read = 0
ORDER BY published_at DESC;

-- name: CreateEntry :one
INSERT INTO entries (
    feed_id,
    title,
    subtitle,
    author,
    content,
    external_url,
    published_at,
    created_at
)
VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;


-- name: GetFeedEntries :many
SELECT feeds.title as feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE feed_id = ?
ORDER BY published_at DESC;

-- name: GetEntry :one
SELECT feeds.Title as feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE entries.id = ?;

-- name: UpdateEntryReadStatus :exec
UPDATE entries
SET read = ?
WHERE id = ?;
