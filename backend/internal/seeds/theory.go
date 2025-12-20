package seeds

import (
	"context"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
)

func SeedAdminTheory(ctx context.Context, q *db.Queries) error {

	if _, err := q.GetTheoryByName(ctx, "Административное право"); err == nil {
		return nil
	}

	_, err := q.CreateTheory(ctx, db.CreateTheoryParams{
		Name:        "Административное право",
		Description: "Основы административной ответственности несовершеннолетних",
		Theoryfull: `
Административная ответственность наступает с 16 лет (ст. 2.3 КоАП РФ).

К административным правонарушениям относятся:
— мелкое хулиганство (ст. 20.1 КоАП РФ)
— распитие алкоголя в общественных местах (ст. 20.20 КоАП РФ)
— нарушение ПДД (гл. 12 КоАП РФ)

Дела несовершеннолетних, как правило, рассматриваются Комиссией по делам несовершеннолетних.
`,
	})
	return err
}
