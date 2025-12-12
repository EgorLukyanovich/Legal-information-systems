package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name" validate:"required, min=2, max=50"`
	LastName  string    `json:"last_name" validate:"required, min=2, max=50"`
	UserName  string    `json:"user_name" validate:"required, min=4, max=25"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required, min=6, max=25"`
	UserTest  int32     `json:"user_test"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	UserTest  int32  `json:"user_test"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
