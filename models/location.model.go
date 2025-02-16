package models

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID          uuid.UUID `json:"userId"`
	SessionID       uuid.UUID `json:"sessionId"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Accuracy        float64   `json:"accuracy"`
	Battery         float64   `json:"battery"`
	Altitude        float64   `json:"altitude"`
	Bearing         float64   `json:"bearing"`
	BearingAccuracy float64   `json:"bearingAccuracy"`
	Speed           float64   `json:"speed"`
	Provider        string    `json:"provider"`
	Timestamp       int       `json:"timestamp"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LocationCreate struct {
	UserID          uuid.UUID `json:"userId" binding:"required"`
	SessionID       uuid.UUID `json:"sessionId" binding:"required"`
	Latitude        float64   `json:"latitude" binding:"required"`
	Longitude       float64   `json:"longitude" binding:"required"`
	Accuracy        float64   `json:"accuracy"`
	Battery         float64   `json:"battery"`
	Altitude        float64   `json:"altitude"`
	Bearing         float64   `json:"bearing"`
	BearingAccuracy float64   `json:"bearingAccuracy"`
	Speed           float64   `json:"speed"`
	Provider        string    `json:"provider" binding:"required"`
	Timestamp       int       `json:"timestamp" binding:"required"`
}
