package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreate struct {
	Username           string  `json:"username" binding:"required"`
	Name               string  `json:"name" binding:"required"`
	Email              string  `json:"email" binding:"required"`
	Password           string  `json:"password" binding:"required"`
	RegistrationSecret *string `json:"registration_secret"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUpdate struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Username string    `json:"username"`
}
