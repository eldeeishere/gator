-- name: CreateUser :one
INSERT INTO users (id, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: AddFeed :one
INSERT INTO feeds (id, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: CreatePost :one
INSERT INTO posts (id, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPostsForUsers :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at
LIMIT $2;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name from feeds f 
INNER JOIN users u ON f.user_id = u.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4)
    RETURNING *
)
SELECT 
    inserted_feed_follow.*, 
    feeds.name AS feed_name, 
    users.name AS user_name 
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;


-- name: GetFeedsUrl :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;


-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name 
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds set last_fetech_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :many
SELECT * FROM feeds ORDER BY last_fetech_at NULLS FIRST;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name from users;

-- name: GetUser :one
SELECT id, created_at, updated_at, name
FROM users
WHERE name = $1
LIMIT 1;