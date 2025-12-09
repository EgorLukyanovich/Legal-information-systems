package models

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         int32     `json:"id"`
	QuestionID int32     `json:"question_id"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
}

type Example struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FullExample string    `json:"full_example"`
	CreatedAt   time.Time `json:"created_at"`
}

type Question struct {
	ID        int32     `json:"id"`
	TestID    uuid.UUID `json:"test_id"`
	Text      string    `json:"text"`
	Multiple  bool      `json:"multiple"`
	CreatedAt time.Time `json:"created_at"`
}

type Test struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Theory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Theoryfull  string    `json:"theoryfull"`
	CreatedAt   time.Time `json:"created_at"`
}
