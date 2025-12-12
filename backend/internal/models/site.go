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
	Theoryfull  string    `json:"theory_full" validate:"required, min=4"`
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
	QuestionID        int32   `json:"question_id"`
	SelectedAnswerIDs []int32 `json:"selected_answer_ids"`
}

type TestAnswerRequest struct {
	TestID  uuid.UUID         `json:"test_id"`
	Answers []UserAnswerInput `json:"answers"`
}

type TestAnswerResponse struct {
	TotalQuestions int  `json:"total_questions"`
	CorrectAnswers int  `json:"correct_answers"`
	IsPassed       bool `json:"is_passed"`
	ScorePercent   int  `json:"score_percent"`
}

//-- Чистые модели для создания

type CreateAnswerInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateQuestionInput struct {
	Text     string              `json:"text"`
	Multiple bool                `json:"multiple"`
	Answers  []CreateAnswerInput `json:"answers"`
}

type CreateTestInput struct {
	Name      string                `json:"test_name"`
	Questions []CreateQuestionInput `json:"questions"`
}
type TheoryInput struct {
	Name        string `json:"theory_name"`
	Description string `json:"theory_description"`
	TheoryFull  string `json:"theory_full"`
}

type ExampleInput struct {
	Name        string `json:"example_name"`
	Description string `json:"example_description"`
	FullExample string `json:"full_example"` //тут еще посмотреть как удобнее будет
}
