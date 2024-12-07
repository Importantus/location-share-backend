package models

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID          uuid.UUID `json:"user_id"`
	SessionID       uuid.UUID `json:"session_id"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Accuracy        float64   `json:"accuracy"`
	Battery         float64   `json:"battery"`
	Altitude        float64   `json:"altitude"`
	Bearing         float64   `json:"bearing"`
	BearingAccuracy float64   `json:"bearing_accuracy"`
	Speed           float64   `json:"speed"`
	Provider        string    `json:"provider"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
