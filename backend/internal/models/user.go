package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" validate:"required, min=2, max=50"`
	LastName  string    `json:"lastName" validate:"required, min=2, max=50"`
	UserName  string    `json:"userName" validate:"required, min=4, max=25"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required, min=6, max=25"`
	UserTest  int32     `json:"userTest"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserProfile struct {
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	UserTest  int32  `json:"userTest"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
