-- name: CreateTest :one
INSERT INTO tests (name) VALUES ($1) RETURNING *;

-- name: GetActiveTest :one
SELECT * FROM tests LIMIT 1;