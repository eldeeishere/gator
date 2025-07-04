-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: GetFeeds :many
SELECT f.name, f.url, u.name from feeds f 
INNER JOIN users u ON f.user_id = u.id;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name from users;

-- name: GetUser :one
SELECT id, created_at, updated_at, name
FROM users
WHERE name = $1
LIMIT 1;