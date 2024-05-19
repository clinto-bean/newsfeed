-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: CreateFeed :one
INSERT INTO feeds (
    id, created_at, updated_at, name, url, user_id
) 
VALUES (
    encode(sha256(random()::text::bytea), 'hex'), $1, $2, $3, $4, $5
) RETURNING *;