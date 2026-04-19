-- name: CreatePost :one

INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostByURL :one

SELECT *
FROM posts
WHERE url = $1;


-- name: GetPostsForUser :many

SELECT *
FROM posts
JOIN feeds ON feeds.id = posts.feed_id
WHERE user = $1
ORDER BY posts.created_at ASC
LIMIT $2;

