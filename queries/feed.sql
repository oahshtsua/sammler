-- name: GetFeeds :many
SELECT *
FROM feeds
ORDER BY title;

-- name: CreateFeed :one
INSERT INTO feeds (
    title,
    subtitle,
    feed_url,
    site_url,
    type,
    updated_at,
    checked_at
)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE feeds.id = ?;

-- name: DeleteFeed :exec
DELETE
FROM feeds
WHERE id = ?;

-- name: MarkFeedRead :exec
UPDATE entries
SET read = 1
WHERE feed_id = ?;

