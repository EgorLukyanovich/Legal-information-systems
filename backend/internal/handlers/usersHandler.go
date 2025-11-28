package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/models"
	json_resp "github.com/egor_lukyanovich/legal-information-systems/backend/pkg/json"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

/*
TODO:
*/

type UserHandlers struct {
	q *db.Queries
}

func NewUserHandlers(q *db.Queries) *UserHandlers {
	return &UserHandlers{q: q}
}

func (u *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	valid := validator.New()
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
		Password:  user.Password,
		UserTest:  0, // дефолт
	}

	if err := valid.Struct(params); err != nil {
		json_resp.RespondError(w, 422, "INVALID_DATA", "invalid data")
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

func (u *UserHandlers) UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	var auth models.UserAuthorization
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid request body")
		return
	}

	user, err := u.q.GetUserByLogin(r.Context(), auth.Login)
	if err != nil {
		json_resp.RespondError(w, 422, "INVALID_DATA", "user not found")
		return
	}

	if auth.Password != user.Password {
		json_resp.RespondError(w, 401, "UNAUTHORIZED", "invalid password")
		return
	}

	//TODO: мб и до jwt дойду
	token := "test_token"

	json_resp.RespondJSON(w, 200, map[string]any{
		"token": token,
	})
}

func (u *UserHandlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		json_resp.RespondError(w, 401, "UNAUTHORIZED", "missing token")
		return
	}

	// TODO: мб и до jwt дойду
	// пока просто заглушка
	json_resp.RespondJSON(w, 200, map[string]any{
		"userName":  "Sava",
		"email":     "sava@mail.ru",
		"userTest":  5,
		"firstName": "Савелий",
		"lastName":  "Булатов",
	})
}
