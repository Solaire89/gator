-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, feed_name, feed_url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UserWhoCreatedFeed :many
SELECT feed_name, feed_url, user_name
FROM feeds
JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE feed_url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;    

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;