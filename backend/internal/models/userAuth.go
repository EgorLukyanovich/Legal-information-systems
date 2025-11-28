package models

type UserAuthorization struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
