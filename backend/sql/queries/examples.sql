-- name: CreateExample :one
INSERT INTO examples (name, description, full_example, created_at) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListExamples :many
SELECT id, name, description, full_example, created_at FROM examples 
ORDER BY created_at DESC;

-- name: GetExampleByName :one
SELECT * FROM examples WHERE name = $1 LIMIT 1;
