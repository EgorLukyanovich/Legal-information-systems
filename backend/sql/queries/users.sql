-- name: CreateUser :one
INSERT INTO users (
id,
created_at,
updated_at
)
VALUES (
$1, $2, $3
)
RETURNING
id,
created_at,
updated_at;

-- name: GetAllUsers :many
SELECT
id,
created_at,
updated_at
FROM users
ORDER BY created_at;

-- name: GetUserByID :one
SELECT
id,
created_at,
updated_at
FROM users
WHERE id = $1;

-- name: UpdateUserUpdatedAt :one
UPDATE users
SET
updated_at = $2
WHERE id = $1
RETURNING
id,
created_at,
updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;