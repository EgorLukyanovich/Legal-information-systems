-- name: CreateAnswer :one
INSERT INTO answers (question_id, text, is_correct) 
VALUES ($1, $2, $3) RETURNING *;

-- name: ListAnswersByQuestionIDs :many
SELECT * FROM answers WHERE question_id = ANY(@question_ids::int[]) ORDER BY id;

-- name: GetTestCorrectAnswers :many
SELECT 
    q.id AS question_id, 
    a.id AS answer_id
FROM questions q
JOIN answers a ON q.id = a.question_id
WHERE q.test_id = $1 AND a.is_correct = TRUE;