package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required"`
	Writing     bool      `json:"writing"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	FCMDeviceID string    `json:"fcm_device_id"`
	LastWakeup  time.Time `json:"last_wakeup"`
	ReadOnly    bool      `json:"read_only"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SessionCreate struct {
	Name     string `json:"name" binding:"required"`
	Writing  bool   `json:"writing"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ReadOnly bool   `json:"read_only"`
}

type RegisterFCMToken struct {
	Token string `json:"token" binding:"required"`
}
