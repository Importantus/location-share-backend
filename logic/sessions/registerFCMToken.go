package sessions

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
)

func RegisterFCMToken(session models.Session, token string) (result models.Session, error customerrors.APIError) {
	// Save the CM token to the session
	session.FCMDeviceID = token
	err := initializers.DB.Save(&session).Error
	if err != nil {
		return result, customerrors.ErrFCMTokenRegistrationFailed
	}
	return session, customerrors.Success
}
