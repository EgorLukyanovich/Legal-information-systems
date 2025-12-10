-- name: CreateTheory :one
INSERT INTO theories (name, description, theoryFull, created_at) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListTheories :many
SELECT id, name, description, theoryFull, created_at FROM theories 
ORDER BY created_at DESC;