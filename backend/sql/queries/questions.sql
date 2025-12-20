-- name: CreateQuestion :one
INSERT INTO questions (test_id, text, multiple) 
VALUES ($1, $2, $3) RETURNING *;

-- name: ListQuestionsByTestID :many
SELECT * FROM questions WHERE test_id = $1 ORDER BY id;

-- name: GetQuestionByTestAndText :one
SELECT * FROM questions
WHERE test_id = $1 AND text = $2
LIMIT 1;
