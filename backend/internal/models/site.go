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
	Name        string    `json:"name" validate:"required, min=4, max=50"`
	Description string    `json:"description" validate:"required, min=2, max=50"`
	FullExample string    `json:"full_example" validate:"required, min=4"`
	CreatedAt   time.Time `json:"created_at"`
}

type Question struct {
	ID        int32     `json:"id"`
	TestID    uuid.UUID `json:"test_id"`
	Text      string    `json:"text" validate:"required, min=4, max=100"`
	Multiple  bool      `json:"multiple"`
	CreatedAt time.Time `json:"created_at"`
}

type Test struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required, min=2, max=50"`
	CreatedAt time.Time `json:"created_at"`
}

type Theory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required, min=4, max=50"`
	Description string    `json:"description" validate:"required, min=2, max=50"`
	Theoryfull  string    `json:"theoryfull" validate:"required, min=4"`
	CreatedAt   time.Time `json:"created_at"`
}

type QuestionWithAnswers struct {
	Question
	Answers []Answer `json:"answers"`
}

type TestFullResponse struct {
	Test
	Questions []QuestionWithAnswers `json:"questions"`
}

type UserAnswerInput struct {
	QuestionID        int32   `json:"questionId"`
	SelectedAnswerIDs []int32 `json:"selectedAnswerIds"`
}

type TestAnswerRequest struct {
	TestID  uuid.UUID         `json:"testId"`
	Answers []UserAnswerInput `json:"answers"`
}

type TestAnswerResponse struct {
	TotalQuestions int  `json:"totalQuestions"`
	CorrectAnswers int  `json:"correctAnswers"`
	IsPassed       bool `json:"isPassed"`
	ScorePercent   int  `json:"scorePercent"`
}
