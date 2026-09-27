-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, user_name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: ClearTable :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT *
FROM users;

-- name: GetUser :one
SELECT * 
FROM users
WHERE user_name = $1;

