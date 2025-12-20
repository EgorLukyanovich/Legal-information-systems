package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/models"
	json_resp "github.com/egor_lukyanovich/legal-information-systems/backend/pkg/json"
)

type SiteHandlers struct {
	q *db.Queries
}

func NewSiteHandlers(queries *db.Queries) *SiteHandlers {
	return &SiteHandlers{q: queries}
}

func (s *SiteHandlers) CreateTheory(w http.ResponseWriter, r *http.Request) {
	var input models.TheoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid body")
		return
	}

	newTheory, err := s.q.CreateTheory(r.Context(), db.CreateTheoryParams{
		Name:        input.Name,
		Description: input.Description,
		Theoryfull:  input.TheoryFull,
	})
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to create theory")
		return
	}

	json_resp.RespondJSON(w, 200, newTheory)
}

func (s *SiteHandlers) CreateExample(w http.ResponseWriter, r *http.Request) {
	var input models.ExampleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid body")
		return
	}

	newExample, err := s.q.CreateExample(r.Context(), db.CreateExampleParams{
		Name:        input.Name,
		Description: input.Description,
		FullExample: input.FullExample,
	})
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to create example")
		return
	}

	json_resp.RespondJSON(w, 200, newExample)
}

func (s *SiteHandlers) CreateTest(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTestInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid body")
		return
	}

	ctx := r.Context()

	newTest, err := s.q.CreateTest(ctx, input.Name)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to create test")
		return
	}

	for _, qInput := range input.Questions {
		newQuestion, err := s.q.CreateQuestion(ctx, db.CreateQuestionParams{
			TestID:   newTest.ID,
			Text:     qInput.Text,
			Multiple: qInput.Multiple,
		})
		if err != nil {
			json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to create question")
			return
		}

		for _, aInput := range qInput.Answers {
			_, err := s.q.CreateAnswer(ctx, db.CreateAnswerParams{
				QuestionID: newQuestion.ID,
				Text:       aInput.Text,
				IsCorrect:  aInput.IsCorrect,
			})
			if err != nil {
				json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to create answer")
				return
			}
		}
	}
	//какого то хуя все варианты ответа по базе false после создания теста, и get questions высерает собственно правильный ответ
	//решилось или нет пока не тестил
	json_resp.RespondJSON(w, 200, map[string]string{"status": "Test created successfully"})
}

func (s *SiteHandlers) GetTheories(w http.ResponseWriter, r *http.Request) {
	theoriesDB, err := s.q.ListTheories(r.Context())
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch theories")
		return
	}

	response := make([]models.Theory, 0, len(theoriesDB))
	for _, t := range theoriesDB {
		response = append(response, models.Theory{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Theoryfull:  t.Theoryfull,
			CreatedAt:   t.CreatedAt,
		})
	}

	json_resp.RespondJSON(w, 200, response)
}

func (s *SiteHandlers) GetExamples(w http.ResponseWriter, r *http.Request) {
	examplesDB, err := s.q.ListExamples(r.Context())
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch examples")
		return
	}

	response := make([]models.Example, 0, len(examplesDB))
	for _, e := range examplesDB {
		response = append(response, models.Example{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			FullExample: e.FullExample,
			CreatedAt:   e.CreatedAt,
		})
	}

	json_resp.RespondJSON(w, 200, response)
}

func (s *SiteHandlers) GetQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	testDB, err := s.q.GetActiveTest(ctx)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch test")
		return
	}

	questionsDB, err := s.q.ListQuestionsByTestID(ctx, testDB.ID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch questions")
		return
	}

	questionIDs := make([]int32, len(questionsDB))
	for i, q := range questionsDB {
		questionIDs[i] = q.ID
	}

	answersDB, err := s.q.ListAnswersByQuestionIDs(ctx, questionIDs)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch answers")
		return
	}

	// Ключ - ID вопроса
	answersMap := make(map[int32][]models.Answer)

	for _, ans := range answersDB {
		qID := int32(ans.QuestionID)

		answersMap[qID] = append(answersMap[qID], models.Answer{
			ID:         int32(ans.ID),
			QuestionID: qID,
			Text:       ans.Text,
			IsCorrect:  false,
			CreatedAt:  ans.CreatedAt,
		})
	}

	questionsResp := make([]models.QuestionWithAnswers, 0, len(questionsDB))

	for _, q := range questionsDB {
		qID := int32(q.ID)

		questionsResp = append(questionsResp, models.QuestionWithAnswers{
			Question: models.Question{
				ID:        qID,
				TestID:    q.TestID,
				Text:      q.Text,
				Multiple:  q.Multiple,
				CreatedAt: q.CreatedAt,
			},
			Answers: answersMap[qID],
		})
	}

	finalResponse := models.TestFullResponse{
		Test: models.Test{
			ID:        testDB.ID,
			Name:      testDB.Name,
			CreatedAt: testDB.CreatedAt,
		},
		Questions: questionsResp,
	}

	json_resp.RespondJSON(w, 200, finalResponse)
}

func (s *SiteHandlers) GetTests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	testDB, err := s.q.GetTests(ctx)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch test")
		return
	}

	response := make([]models.Example, 0, len(testDB))
	for _, e := range testDB {
		response = append(response, models.Example{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
		})
	}

	json_resp.RespondJSON(w, 200, response)

}

func (s *SiteHandlers) SubmitTestAnswers(w http.ResponseWriter, r *http.Request) {
	var req models.TestAnswerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid request body")
		return
	}

	correctAnswersDB, err := s.q.GetTestCorrectAnswers(r.Context(), req.TestID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to check answers")
		return
	}

	correctMap := make(map[int32][]int32)
	for _, row := range correctAnswersDB {
		qID := int32(row.QuestionID)
		aID := int32(row.AnswerID)
		correctMap[qID] = append(correctMap[qID], aID)
	}

	score := 0
	for _, userAns := range req.Answers {
		correctIDs, exists := correctMap[userAns.QuestionID]
		if !exists {
			continue // вопроса нет в базе или у него нет правильных ответов
		}

		if len(userAns.SelectedAnswerIDs) != len(correctIDs) {
			continue
		}

		isCorrect := true
		for _, selectedID := range userAns.SelectedAnswerIDs {
			found := false
			for _, correctID := range correctIDs {
				if selectedID == correctID {
					found = true
					break
				}
			}
			if !found {
				isCorrect = false
				break
			}
		}

		if isCorrect {
			score++
		}
	}

	resp := models.TestAnswerResponse{
		TestID: req.TestID,
		Score:  score,
	}

	json_resp.RespondJSON(w, 200, resp)
}
