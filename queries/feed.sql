-- name: GetFeeds :many
SELECT *
FROM feeds
ORDER BY title;

-- name: CreateFeed :one
INSERT INTO feeds (
    title,
    feed_url,
    site_url,
    updated_at,
    checked_at
)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE feeds.id = ?;

-- name: DeleteFeed :exec
DELETE
FROM feeds
WHERE id = ?;

-- name: UpdateFeedCheckedAt :exec
UPDATE feeds
SET checked_at = ?
WHERE id = ?;
