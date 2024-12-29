package sessions

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"

	"github.com/google/uuid"
)

func GetSession(sessionId uuid.UUID) (session models.Session, error customerrors.APIError) {

	// Check if the session exists
	err := initializers.DB.Where(&models.Session{
		ID: sessionId,
	}).First(&session).Error

	if err != nil {
		return models.Session{}, customerrors.ErrSessionNotFound
	}

	return session, customerrors.Success
}
