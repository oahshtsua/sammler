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

-- name: GetUnreadEntries :many
SELECT feeds.title AS feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE read = 0
ORDER BY published_at DESC;

-- name: GetFeedEntries :many
SELECT feeds.title as feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE feed_id = ?
ORDER BY published_at DESC;

-- name: GetEntry :one
SELECT feeds.title as feed_title, entries.*
FROM entries
JOIN feeds
    ON entries.feed_id = feeds.id
WHERE entries.id = ?;

-- name: MarkEntriesRead :exec
UPDATE entries
SET read = 1
WHERE read = 0;

-- name: MarkEntryRead :exec
UPDATE entries
SET read = 1
WHERE id = ?;


-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = ?;
