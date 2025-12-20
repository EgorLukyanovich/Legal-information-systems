package seeds

import (
	"context"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
)

func SeedAdminLawTest(ctx context.Context, q *db.Queries) error {

	_, err := q.GetTestByName(ctx, "Административное право")
	if err == nil {
		return nil
	}

	test, err := q.CreateTest(ctx, "Административное право")
	if err != nil {
		return err
	}

	questions := []struct {
		Text    string
		Answers []struct {
			Text      string
			IsCorrect bool
		}
	}{
		{
			Text: "17-летний Денис в автобусе громко включал музыку и нецензурно выражался. Какое правонарушение он совершил?",
			Answers: []struct {
				Text      string
				IsCorrect bool
			}{
				{"Уголовное преступление по ст. 213 УК РФ", false},
				{"Нарушение правил пользования транспортом", false},
				{"Административное правонарушение по ст. 20.1 КоАП РФ", true},
				{"Оскорбление по ст. 5.61 КоАП РФ", false},
			},
		},
		{
			Text: "16-летний Артём управлял автомобилем без водительского удостоверения. Какая ответственность наступает?",
			Answers: []struct {
				Text      string
				IsCorrect bool
			}{
				{"Никто не будет привлечён", false},
				{"Только отец — штраф", false},
				{"Артём и отец — административная ответственность", true},
				{"Только Артём — штраф", false},
			},
		},
	}

	for _, qData := range questions {

		_, err := q.GetQuestionByTestAndText(ctx, db.GetQuestionByTestAndTextParams{TestID: test.ID, Text: qData.Text})
		if err == nil {
			continue
		}

		qRow, err := q.CreateQuestion(ctx, db.CreateQuestionParams{
			TestID:   test.ID,
			Text:     qData.Text,
			Multiple: false,
		})
		if err != nil {
			return err
		}

		for _, a := range qData.Answers {

			if _, err := q.GetAnswerByQuestionAndText(ctx, db.GetAnswerByQuestionAndTextParams{QuestionID: qRow.ID, Text: a.Text}); err == nil {
				continue
			}

			if _, err := q.CreateAnswer(ctx, db.CreateAnswerParams{
				QuestionID: qRow.ID,
				Text:       a.Text,
				IsCorrect:  a.IsCorrect,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
