-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, feed_name, feed_url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: UserWhoCreatedFeed :many
SELECT feed_name, feed_url, user_name
FROM feeds
JOIN users
ON feeds.user_id = users.id;