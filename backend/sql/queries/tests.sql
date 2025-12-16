-- name: CreateTest :one
INSERT INTO tests (name) VALUES ($1) RETURNING *;

-- name: GetActiveTest :one
SELECT * FROM tests LIMIT 1;

-- name: GetTests :many
SELECT id, name, created_at FROM tests 
ORDER BY created_at DESC;