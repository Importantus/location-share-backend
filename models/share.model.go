package models

import (
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey"`
	ValidUntil *time.Time `json:"valid_until"`
	SharedWith uuid.UUID  `json:"shared_with"`
	SharedBy   uuid.UUID  `json:"shared_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ShareCreate struct {
	ValidUntil *time.Time `json:"valid_until"`
	SharedWith uuid.UUID  `json:"shared_with" binding:"required"`
}
