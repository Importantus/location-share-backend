package sessions

import (
	"location-share-backend/initializers"
	"location-share-backend/models"
)

func ListSessions(session models.Session) (sessions []models.Session, err error) {
	err = initializers.DB.Find(&sessions, &models.Session{
		UserID: session.UserID,
	}).Error

	if err != nil {
		return nil, err
	}

	return
}
