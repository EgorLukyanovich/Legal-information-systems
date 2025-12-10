package handlers

import (
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
