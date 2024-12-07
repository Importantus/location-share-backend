package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" binding:"required"`
	Writing    bool      `json:"writing" binding:"required"`
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	LastWakeup time.Time `json:"last_wakeup"`
	ReadOnly   bool      `json:"read_only" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SessionCreate struct {
	Name     string `json:"name" binding:"required"`
	Writing  bool   `json:"writing" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ReadOnly bool   `json:"read_only" binding:"required"`
}
