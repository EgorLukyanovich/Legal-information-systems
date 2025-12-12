package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/db"
	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/models"
	json_resp "github.com/egor_lukyanovich/legal-information-systems/backend/pkg/json"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		json_resp.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "could not hash password")
		return
	}

	params := db.CreateUserParams{
		ID:        uuid.New(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		Email:     user.Email,
		Password:  string(hashedPassword),
		UserTest:  0, // дефолт
	}

	if err := valid.Struct(params); err != nil {
		json_resp.RespondError(w, 422, "INVALID_DATA", "invalid data")
		return
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
	//а чтобы не палить
	user.Password = ""

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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password)); err != nil {
		json_resp.RespondError(w, 401, "UNAUTHORIZED", "invalid password")
		return
	}

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		log.Fatal("JWT secret is not found in .env")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // протухнет через 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "token generation failed")
		return
	}

	json_resp.RespondJSON(w, 200, map[string]any{
		"token": tokenString,
	})

}

func (u *UserHandlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to get user context")
		return
	}

	user, err := u.q.GetUserByID(r.Context(), userID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "database error: "+err.Error())
		return
	}

	json_resp.RespondJSON(w, 200, map[string]any{
		"userName":  user.UserName,
		"email":     user.Email,
		"userTest":  user.UserTest,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})
}
