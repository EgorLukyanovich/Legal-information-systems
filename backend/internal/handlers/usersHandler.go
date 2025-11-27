package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/models"
	json_resp "github.com/egor_lukyanovich/legal-information-systems/backend/pkg/json"
	"github.com/google/uuid"
)

/*
TODO: Проверку на валидность данных, '422' код
*/

type UserHandlers struct {
	q *db.Queries
}

func NewUserHandlers(q *db.Queries) *UserHandlers {
	return &UserHandlers{q: q}
}

func (u *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		json_resp.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	params := db.CreateUserParams{
		ID:        uuid.New(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		Email:     user.Email,
		Password:  user.Password, // пока так
		UserTest:  0,             // дефолт
	}

	// Проверка: мыла
	if _, err := u.q.GetUserByEmail(r.Context(), user.Email); err == nil {
		json_resp.RespondError(w, 500, "CONFLICT", "email already registered")
		return
	}

	if err := u.q.CreateUser(r.Context(), params); err != nil {
		json_resp.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"user": user,
	})
}
