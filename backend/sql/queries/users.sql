-- name: CreateUser :exec
INSERT INTO users (
    id,
    first_name,
    last_name,
    user_name,
    email,
    password,
    user_test
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id,
    first_name,
    last_name,
    user_name,
    email,
    user_test,
    created_at,
    updated_at;

-- name: GetUserByID :one
SELECT
    id,
    first_name,
    last_name,
    user_name,
    email,
    user_test,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- name: GetUserByLogin :one
SELECT
    id,
    first_name,
    last_name,
    user_name,
    email,
    password,
    user_test,
    created_at,
    updated_at
FROM users
WHERE user_name = $1 OR email = $1;

-- name: GetUserByEmail :one
SELECT
    id,
    first_name,
    last_name,
    user_name,
    email,
    password,
    user_test,
    created_at,
    updated_at
FROM users
WHERE email = $1;


-- name: UpdateUserUserTest :one
UPDATE users
SET
    user_test = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    first_name,
    last_name,
    user_name,
    email,
    user_test,
    created_at,
    updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
